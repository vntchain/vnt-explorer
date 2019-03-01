package main

import (
	"github.com/astaxie/beego"
	"fmt"
	"time"
	"github.com/vntchain/vnt-explorer/tools/racer/data"
)

func main() {
	rpcHost := beego.AppConfig.String("node::rpc_host")
	rpcPort := beego.AppConfig.String("node::rpc_port")

	beego.Info("rpc host: ", rpcHost)
	beego.Info("rpc port: ", rpcPort)

	for {
		rmtHgt, localHgt := checkHeight()
		beego.Info(fmt.Sprintf("Local height: %d, rmtHeight: %d", localHgt, rmtHgt))

		if localHgt == rmtHgt {
			time.Sleep(1)
			continue
		}

		for localHgt < rmtHgt {
			data.GetBlock(localHgt+1)
			break;
		}

		break
	}
}

func checkHeight() (int64, int64) {
	rmtHgt := data.GetRemoteHeight()
	localHgt := data.GetLocalHeight()

	if localHgt > rmtHgt {
		msg := fmt.Sprintf("Local height %d is bigger than remote height: %d, please check your remote node", localHgt, rmtHgt)
		beego.Error(msg)
		panic(msg)
	}

	return rmtHgt, localHgt
}



