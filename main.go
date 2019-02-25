package main

import (
	_ "github.com/vntchain/vnt-explorer/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.BConfig.Log.AccessLogs = true
	beego.Run()
}

