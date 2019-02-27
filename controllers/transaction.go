package controllers

import (
	"github.com/vntchain/vnt-explorer/models"
	"encoding/json"
	"github.com/vntchain/vnt-explorer/common"
	"github.com/astaxie/beego"
)

type TransactionController struct {
	BaseController
}

func (this *TransactionController) Post() {
	tx := &models.Transaction{}
	body := this.Ctx.Input.RequestBody
	err := json.Unmarshal(body, tx)

	if err != nil {
		this.ReturnErrorMsg("Wrong format of transaction: %s", err.Error())
		return
	}

	err = tx.Insert()

	if err != nil {
		this.ReturnErrorMsg("Failed to create transaction: %s", err.Error())
	} else {
		this.ReturnData(tx)
	}
}

func (this *TransactionController) List() {
	offset, err := this.GetInt64("offset");
	if err != nil {
		beego.Warn("Failed to read offset: ", err.Error())
		offset = common.DefaultOffset
	}

	limit, err := this.GetInt64("limit")
	if err != nil {
		beego.Warn("Failed to read limit: ", err.Error())
		limit = common.DefaultPageSize
	}

	fields := this.getFields()

	block := this.GetString("block")
	account := this.GetString("account")

	isToken, err := this.GetInt("istoken")
	if err != nil {
		beego.Warn("Failed to read istoken: ", err.Error())
		isToken = -1
	}

	tx := &models.Transaction{}
	txs, err := tx.List(offset, limit, block, account, isToken, fields...)

	if err != nil {
		this.ReturnErrorMsg("Failed to list transactions: ", err.Error())
	} else {
		this.ReturnData(txs)
	}
}

func (this *TransactionController) Get() {
	//beego.Info("params", this.Ctx.Input.Params())
	txHash := this.Ctx.Input.Param(":tx_hash")
	if len(txHash) == 0 {
		this.ReturnErrorMsg("Failed to get block number or hash", "")
		return
	}

	fields := this.getFields()
	beego.Info("Will read colums: ", fields, "txhash", txHash)

	tx := &models.Transaction{}
	dbTx, err := tx.Get(txHash, fields...)
	if err != nil {
		this.ReturnErrorMsg("Failed to read transaction: %s", err.Error())
	} else {
		this.ReturnData(dbTx)
	}
}

func (this *TransactionController) Count() {
	block := this.GetString("block")
	account := this.GetString("account")

	isToken, err := this.GetInt("istoken")
	if err != nil {
		beego.Warn("Failed to read istoken: ", err.Error())
		isToken = -1
	}

	tx := &models.Transaction{}
	count, err := tx.Count(block, account, isToken)

	if err != nil {
		this.ReturnErrorMsg("Failed to count transactions: ", err.Error())
	} else {
		this.ReturnData(count)
	}
}