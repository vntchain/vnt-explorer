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
		this.ReturnData(tokenBalance, nil)
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
	order := this.GetString("order")
	fields := this.getFields()

	tokenBalance := &models.TokenBalance{}
	dbItemList, err := tokenBalance.List("", tokenAddress, order, offset, limit, fields)
	if err != nil {
		this.ReturnErrorMsg("Failed to list TokenBalance: %s", err.Error())
	} else {
		count := make(map[string]int64)
		count["count"], err = tokenBalance.Count("", tokenAddress)
		if err != nil {
			this.ReturnErrorMsg("Failed to list TokenBalance: %s", err.Error())
			return
		}
		this.ReturnData(dbItemList, count)
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
	order := this.GetString("order")
	fields := this.getFields()

	tokenBalance := &models.TokenBalance{}
	dbItemList, err := tokenBalance.List(account, "", order, offset, limit, fields)
	if err != nil {
		this.ReturnErrorMsg("Failed to list TokenBalance: %s", err.Error())
	} else {
		count := make(map[string]int64)
		count["count"], err = tokenBalance.Count(account, "")
		if err != nil {
			this.ReturnErrorMsg("Failed to list TokenBalance: %s", err.Error())
			return
		}
		this.ReturnData(dbItemList, count)
	}

}

func (this *TokenBalanceController) TokenCount() {
	account := this.Ctx.Input.Param(":address")

	tokenBalance := &models.TokenBalance{}
	count, err := tokenBalance.Count(account, "")
	if err != nil {
		this.ReturnErrorMsg("Failed to get account token count: %s", err.Error())
	} else {
		this.ReturnData(count, nil)
	}
}

func (this *TokenBalanceController) HolderCount() {
	tokenAddress := this.Ctx.Input.Param(":address")

	tokenBalance := &models.TokenBalance{}
	count, err := tokenBalance.Count("", tokenAddress)
	if err != nil {
		this.ReturnErrorMsg("Failed to get account token count: %s", err.Error())
	} else {
		this.ReturnData(count, nil)
	}
}