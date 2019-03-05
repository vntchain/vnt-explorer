package utils

import (
	"github.com/vntchain/vnt-explorer/common"
	"encoding/json"
	"bytes"
	"fmt"
	"github.com/astaxie/beego"
	"net/http"
	"io/ioutil"
)

var rpcHost = beego.AppConfig.String("node::rpc_host")
var rpcPort = beego.AppConfig.String("node::rpc_port")
var rpcApi = fmt.Sprintf("http://%s:%s/", rpcHost, rpcPort)

func CallRpc(rpc *common.Rpc) *common.Response {
	rpcJson, err := json.Marshal(rpc)

	buf := bytes.NewBuffer(rpcJson)

	if err != nil {
		msg := fmt.Sprint("Failed to parse rpc %s", err.Error())
		beego.Error(msg)
		panic(msg)
	}

	beego.Info("Will call rpc with request: ", buf.String())

	resp, err := http.Post(rpcApi, common.H_ContentType, buf)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		msg := fmt.Sprintf("Failed to read response body: %s", err.Error())
		beego.Error(msg)
		panic(msg)
	}

	obj := new(common.Response)
	err = json.Unmarshal(body, obj)
	if err != nil {
		msg := fmt.Sprintf("Failed to unmarshal json: %s", err.Error())
		beego.Error(msg)
		panic(msg)
	}

	if obj.Error != nil {
		msg := fmt.Sprintf("Rpc returned with error: code: %d, error: %s", obj.Error.Code, obj.Error.Message)
		beego.Error(msg)
		panic(msg)
	}

	return obj
}