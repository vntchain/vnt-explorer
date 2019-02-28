package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/common"
	"github.com/vntchain/vnt-explorer/models"
)

type AccountController struct {
	BaseController
}

func (this *AccountController) Post() {
	account := &models.Account{}
	body := this.Ctx.Input.RequestBody
	err := json.Unmarshal(body, account)
	if err != nil {
		this.ReturnErrorMsg("Wrong format of Account: %s", err.Error())
		return
	}

	err = account.Insert()
	if err != nil {
		this.ReturnErrorMsg("Failed to create Account: %s", err.Error())
	} else {
		this.ReturnData(account)
	}
}

func (this *AccountController) List() {
	offset, err := this.GetInt("offset")
	if err != nil {
		beego.Warn("Failed to read offset: ", err.Error())
		offset = common.DefaultOffset
	}

	limit, err := this.GetInt("limit")
	if err != nil {
		beego.Warn("Failed to read limit: ", err.Error())
		limit = common.DefaultPageSize
	}

	account := &models.Account{}
	accounts, err := account.List(offset, limit)
	if err != nil {
		this.ReturnErrorMsg("Failed to list accounts: %s", err.Error())
	} else {
		this.ReturnData(accounts)
	}

}

func (this *AccountController) Get() {
	//beego.Info("params", this.Ctx.Input.Params())
	address := this.Ctx.Input.Param(":address")
	if len(address) == 0 {
		this.ReturnErrorMsg("Failed to get address", "")
		return
	}

	account := &models.Account{}
	dbaccount, err := account.Get(address)
	if err != nil {
		this.ReturnErrorMsg("Failed to read account: %s", err.Error())
	} else {
		this.ReturnData(dbaccount)
	}
}

func (this *AccountController) Count() {
	account := &models.Account{}
	count, err := account.Count()
	if err != nil {
		this.ReturnErrorMsg("Failed to get accounts count: %s", err.Error())
	} else {
		this.ReturnData(count)
	}
}
