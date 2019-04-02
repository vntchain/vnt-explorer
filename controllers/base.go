package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"github.com/vntchain/vnt-explorer/common"
)

type BaseController struct {
	beego.Controller
}

type Response struct {
	Ok   int         	`json:"ok"`
	Err  *common.ErrorMessage  `json:"err"`
	Data interface{} 	`json:"data"`
	Extra interface{} 	`json:"extra"`
}

func makeResp(err string, errorCode string, data interface{}, extra interface{}) *Response {
	isOk := len(err) == 0
	var ok int
	if isOk {
		ok = 1
	} else {
		ok = 0
	}

	var Error *common.ErrorMessage = nil
	if err != "" {
		Error = &common.ErrorMessage{
			Code: errorCode,
			Message: err,
		}
	}

	resp := &Response{
		ok,
		Error,
		data,
		extra,
	}

	beego.Info("Response: ", resp)

	return resp
}

func (c *BaseController) ReturnErrorMsg(format, err string, errCode string) {
	msg := fmt.Sprintf(format, err)
	beego.Error(msg)
	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = makeResp(msg, errCode, nil, nil)
	c.ServeJSON()
}

func (c *BaseController) ReturnData(data interface{}, extra interface{}) {
	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = makeResp("", "", data, extra)
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
