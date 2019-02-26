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

	err, dbblock := block.Insert()

	if err != nil {
		this.ReturnErrorMsg("Failed to create block: %s", err.Error())
	} else {
		this.ReturnData(200, dbblock)
	}
}

func (this *BlockController) List() {
	offset, err := this.GetInt("offset");
	if err != nil {
		beego.Warn("Failed to read offset: ", err.Error())
		offset = common.DefaultOffset
	}

	limit, err := this.GetInt("limit")
	if err != nil {
		beego.Warn("Failed to read limit: ", err.Error())
		limit = common.DefaultPageSize
	}

	block := &models.Block{}
	err, blocks := block.List(offset, limit)
	if err != nil {
		this.ReturnErrorMsg("Failed to list blocks: %s", err.Error())
	} else {
		this.ReturnData(200, blocks)
	}

}

func (this *BlockController) Get() {
	//beego.Info("params", this.Ctx.Input.Params())
	nOrh := this.Ctx.Input.Param(":n_or_h")
	if len(nOrh) == 0 {
		this.ReturnErrorMsg("Failed to get block number or hash", "")
		return
	}

	block := &models.Block{}
	err, dbblock := block.Get(nOrh)
	if err != nil {
		this.ReturnErrorMsg("Failed to read block: %s", err.Error())
	} else {
		this.ReturnData(200, dbblock)
	}
}