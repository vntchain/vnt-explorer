package controllers

import (
	"github.com/astaxie/beego"
	"github.com/vntchain/Masami/backEnd/protocols"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) ReturnErrorMsg(msg string) {
	beego.Error(msg)
	c.Ctx.Output.SetStatus(500)
	c.Data["json"] = &protocols.ErrorMessage{
		Message: msg,
	}
	c.ServeJSON()
}