package controllers

import (
	"strings"
	"github.com/vntchain/vnt-explorer/models"
	"github.com/vntchain/vnt-explorer/common"
	"github.com/astaxie/beego/orm"
	"fmt"
)

type ErrorWrongKey struct {
	format	string
	keyword	string
}

func (e ErrorWrongKey) Error() string {
	return fmt.Sprintf(e.format, e.keyword)
}


type SearchController struct {
	BaseController
}

func (this *SearchController) Search() {
	sBody := &SearchBody{}
	var err error
	keyword := this.Ctx.Input.Param(":keyword")
	if keyword == "" {
		this.ReturnErrorMsg("Error happend: %s", "Wrong format of keyword", common.ERROR_WRONG_KEYWORD)
		return
	}

	if !strings.HasPrefix(keyword, "0x") {
		b := &models.Block{}
		b, err = b.Get(keyword)

		sBody.Block = b
	} else {
		switch len(keyword) {
		case 42:
			acct := &models.Account{}
			acct, err = acct.Get(keyword)
			if err == nil {
				sBody.Account = acct
			}
			break
		case 66:
			b := &models.Block{}
			b, err = b.Get(keyword)

			if err == nil {
				sBody.Block = b
				break
			}

			t := &models.Transaction{}
			t, err = t.Get(keyword)
			if err == nil {
				sBody.Tx = t
			}
			break
		default:
			err = ErrorWrongKey{"wrong format of keyword: %s", keyword}
		}
	}

	if err != nil {
		if err == orm.ErrNoRows {
			this.ReturnErrorMsg("Not data found %s", err.Error(), common.ERROR_NOT_FOUND)
		} else {

		}
		switch err.(type) {
		case models.ErrorBlockNumber:
			this.ReturnErrorMsg("Block Number incorrect, %s", err.Error(), common.ERROR_WRONG_KEYWORD)
			break
		case ErrorWrongKey:
			this.ReturnErrorMsg("%s", err.Error(), common.ERROR_WRONG_KEYWORD)
			break
		default:
			this.ReturnErrorMsg("Error happend: %s", err.Error(), common.ERROR_SEARCH_ERROR)
		}
	} else {
		this.ReturnData(sBody, nil)
	}
}
