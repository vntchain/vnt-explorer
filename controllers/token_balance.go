package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/common"
	"github.com/vntchain/vnt-explorer/models"
)

type TokenBalanceController struct {
	BaseController
}

func (this *TokenBalanceController) Post() {
	tokenBalance := &models.TokenBalance{}
	body := this.Ctx.Input.RequestBody
	err := json.Unmarshal(body, tokenBalance)
	if err != nil {
		this.ReturnErrorMsg("Wrong format of TokenBalance: %s", err.Error())
		return
	}

	err = tokenBalance.Insert()
	if err != nil {
		this.ReturnErrorMsg("Failed to create TokenBalance: %s", err.Error())
	} else {
		this.ReturnData(tokenBalance)
	}
}

func (this *TokenBalanceController) ListByToken() {
	tokenAddress := this.Ctx.Input.Param(":address")

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

	tokenBalance := &models.TokenBalance{}
	dbItemList, err := tokenBalance.List("", tokenAddress, offset, limit)
	if err != nil {
		this.ReturnErrorMsg("Failed to list TokenBalance: %s", err.Error())
	} else {
		this.ReturnData(dbItemList)
	}

}

func (this *TokenBalanceController) ListByAccount() {
	account := this.Ctx.Input.Param(":address")

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

	tokenBalance := &models.TokenBalance{}
	dbItemList, err := tokenBalance.List(account, "", offset, limit)
	if err != nil {
		this.ReturnErrorMsg("Failed to list TokenBalance: %s", err.Error())
	} else {
		this.ReturnData(dbItemList)
	}

}

func (this *TokenBalanceController) CountByToken() {
	token := &models.TokenBalance{}
	count, err := token.CountByToken()
	if err != nil {
		this.ReturnErrorMsg("Failed to get token count: %s", err.Error())
	} else {
		this.ReturnData(count)
	}
}
