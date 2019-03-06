package controllers

import (
	"github.com/vntchain/vnt-explorer/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/common"
)

type BlockController struct {
	BaseController
}

func (this *BlockController) Post() {
	block := &models.Block{}
	body := this.Ctx.Input.RequestBody
	err := json.Unmarshal(body, block)
	if err != nil {
		this.ReturnErrorMsg("Wrong format of block: %s", err.Error())
		return
	}

	err = block.Insert()

	if err != nil {
		this.ReturnErrorMsg("Failed to create block: %s", err.Error())
	} else {
		this.ReturnData(block)
	}
}

func (this *BlockController) List() {
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

	order := this.GetString("order")

	fields := this.getFields()

	block := &models.Block{}
	blocks, err := block.List(offset, limit, order, fields...)
	if err != nil {
		this.ReturnErrorMsg("Failed to list blocks: %s", err.Error())
	} else {
		this.ReturnData(blocks)
	}

}

func (this *BlockController) Get() {
	//beego.Info("params", this.Ctx.Input.Params())
	nOrh := this.Ctx.Input.Param(":n_or_h")
	if len(nOrh) == 0 {
		this.ReturnErrorMsg("Failed to get block number or hash", "")
		return
	}

	fields := this.getFields()
	beego.Info("Will read colums: ", fields, "number", nOrh)

	block := &models.Block{}
	dbblock, err := block.Get(nOrh, fields...)
	if err != nil {
		this.ReturnErrorMsg("Failed to read block: %s", err.Error())
	} else {
		this.ReturnData(dbblock)
	}
}

func (this *BlockController) Count() {
	block := &models.Block{}
	count, err := block.Count()
	if err != nil {
		this.ReturnErrorMsg("Failed to get block count: %s", err.Error())
	} else {
		this.ReturnData(count)
	}
}