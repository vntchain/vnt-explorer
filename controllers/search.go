package controllers

import (
	"strings"
	"github.com/vntchain/vnt-explorer/models"
	"errors"
)

type SearchController struct {
	BaseController
}

func (this *SearchController) Search() {
	sBody := &SearchBody{}
	var err error
	keyword := this.Ctx.Input.Param(":keyword")
	if keyword == "" {
		this.ReturnErrorMsg("Error happend: %s", "Wrong format of keyword")
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
			err = errors.New("wrong format of keyword")
		}
	}

	if err != nil {
		this.ReturnErrorMsg("Error happend: %s", err.Error())
	} else {
		this.ReturnData(sBody)
	}
}
