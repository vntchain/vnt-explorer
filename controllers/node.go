package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/common"
	"github.com/vntchain/vnt-explorer/models"
)

type NodeController struct {
	BaseController
}

func (this *NodeController) Post() {
	node := &models.Node{}
	body := this.Ctx.Input.RequestBody
	err := json.Unmarshal(body, node)
	if err != nil {
		this.ReturnErrorMsg("Wrong format of Node: %s", err.Error(), "")
		return
	}

	err = node.Insert()
	if err != nil {
		this.ReturnErrorMsg("Failed to create Node: %s", err.Error(), "")
	} else {
		this.ReturnData(node, nil)
	}
}

func (this *NodeController) List() {
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

	node := &models.Node{}
	nodes, err := node.List(order, offset, limit, fields)
	if err != nil {
		this.ReturnErrorMsg("Failed to list nodes: %s", err.Error(), "")
	} else {
		count := make(map[string]int64)
		count["count"], err = node.Count(-1)
		if err != nil {
			this.ReturnErrorMsg("Failed to list nodes: %s", err.Error(), "")
			return
		}
		this.ReturnData(nodes, count)
	}

}

func (this *NodeController) Count() {
	//status, err := this.GetInt("status")
	//if err != nil {
	//	beego.Warn("Failed to read status: ", err.Error())
	//	status = common.DefaultNodeStatus
	//}

	node := &models.Node{}

	superCount, err := node.Count(1)
	if err != nil {
		this.ReturnErrorMsg("Failed to get node count: %s", err.Error(), "")
		return
	}

	candiCount, err := node.Count(0)
	if err != nil {
		this.ReturnErrorMsg("Failed to get node count: %s", err.Error(), "")
		return
	}

	type Result struct {
		Super	int64
		Candi	int64
		Total	int64
	}

	result := &Result{
		superCount,
		candiCount,
		superCount + candiCount,
	}

	this.ReturnData(result, nil)
}

func (this *NodeController) Get() {
	//beego.Info("params", this.Ctx.Input.Params())
	address := this.Ctx.Input.Param(":address")
	if len(address) == 0 {
		this.ReturnErrorMsg("Failed to get address", "", "")
		return
	}

	node := &models.Node{}
	dbItem, err := node.Get(address)
	if err != nil {
		this.ReturnErrorMsg("Failed to read node: %s", err.Error(), "")
	} else {
		this.ReturnData(dbItem, nil)
	}
}
