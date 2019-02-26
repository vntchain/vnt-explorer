package controllers

import (
	"encoding/json"
	"github.com/vntchain/vnt-explorer/models"
	"github.com/astaxie/beego/orm"
	"strconv"
	"github.com/astaxie/beego"
)

type TestController struct {
	BaseController
}

func (this *TestController) Get() {
	id := this.Ctx.Input.Param(":id")
	beego.Info("Params: ", this.Ctx.Input.Params())
	beego.Info("Will get test of id: ", id)
	test := new(models.Test)
	var err error
	test.Id, err = strconv.Atoi(id)
	if err != nil {
		this.ReturnErrorMsg("Invalid id, err : %s", err.Error())
	}

	o := orm.NewOrm()

	err = o.Read(test)
	if err != nil {
		this.ReturnErrorMsg("read error, err : %s", err.Error())
	}

	this.Ctx.Output.SetStatus(200)
	this.Data["json"] = test
	this.ServeJSON()
}

func (this *TestController) Post() {
	testBody := &models.Test{}
	body := this.Ctx.Input.RequestBody
	beego.Info("Will create a Test", "body", body)
	err := json.Unmarshal(body, testBody)

	if err != nil {
		this.ReturnErrorMsg("fail err : %s", err.Error())
	}
	test := new(models.Test)
	test.Name = testBody.Name
	o := orm.NewOrm()
	id, err := o.Insert(test)
	test.Id = int(id)
	if err != nil {
		this.ReturnErrorMsg("fail, err : %s", err.Error())
	} else {
		this.ReturnData(200, test)
	}
}