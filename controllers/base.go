package controllers

import (
	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/common"
	"fmt"
	"strings"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) ReturnErrorMsg(format, err string) {
	msg := fmt.Sprintf(format, err)
	beego.Error(msg)
	c.Ctx.Output.SetStatus(500)
	c.Data["json"] = &common.ErrorMessage{
		Message: msg,
	}
	c.ServeJSON()
}

func (c *BaseController) ReturnData(status int, data interface{}) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = data

	c.ServeJSON()
}

func (c *BaseController) getFields() []string {
	fieldStr := c.GetString("fields")
	if len(fieldStr) > 0 {
		return strings.Split(fieldStr, ",")
	} else {
		return make([]string, 0)
	}
}