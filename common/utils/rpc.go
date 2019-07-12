package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/common"
	"io/ioutil"
	"net/http"
)

var rpcHost = beego.AppConfig.String("node::rpc_host")
var rpcPort = beego.AppConfig.String("node::rpc_port")
var rpcApi = fmt.Sprintf("http://%s:%s/", rpcHost, rpcPort)

func CallRpc(rpc *common.Rpc) (error, *common.Response, *common.Error) {
	rpcJson, err := json.Marshal(rpc)

	buf := bytes.NewBuffer(rpcJson)

	if err != nil {
		msg := fmt.Sprint("Failed to parse rpc %s", err.Error())
		beego.Error(msg)
		return errors.New(msg), nil, nil
	}

	beego.Debug("Will call rpc with request: ", buf.String())

	resp, err := http.Post(rpcApi, common.H_ContentType, buf)

	if resp == nil || resp.Body == nil {
		msg := fmt.Sprintf("Failed to get resp body")
		beego.Error(msg)
		return errors.New(msg), nil, nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		msg := fmt.Sprintf("Failed to read response body: %s", err.Error())
		beego.Error(msg)
		return errors.New(msg), nil, nil
	}

	obj := new(common.Response)
	err = json.Unmarshal(body, obj)
	if err != nil {
		msg := fmt.Sprintf("Failed to unmarshal json: %s", err.Error())
		beego.Error(msg)
		return errors.New(msg), nil, nil
	}

	if obj.Error != nil {
		msg := fmt.Sprintf("Rpc returned with error: code: %d, error: %s", obj.Error.Code, obj.Error.Message)
		beego.Warn(msg)
		return errors.New(msg), nil, obj.Error
	}

	return nil, obj, obj.Error
}
