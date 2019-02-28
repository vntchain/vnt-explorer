package main

import "github.com/astaxie/beego"

func main() {
	rpcHost := beego.AppConfig.String("node::rpc_host")
	rpcPort := beego.AppConfig.String("node::rpc_port")

	beego.Info("rpc host: ", rpcHost)
	beego.Info("rpc port: ", rpcPort)


}

func checkHeight() {
	
}