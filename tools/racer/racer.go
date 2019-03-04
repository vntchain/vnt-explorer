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

		//localHgt = 89
		//rmtHgt = 20000
		beego.Info(fmt.Sprintf("Local height: %d, rmtHeight: %d", localHgt, rmtHgt))
		if localHgt == rmtHgt {
			time.Sleep(1 * time.Second)
			continue
		}

		for localHgt < rmtHgt {
			block, txs, witnesses := data.GetBlock(localHgt+1)

			beego.Info("Block:", block)
			beego.Info("txs:", txs)
			beego.Info("witness:", witnesses)

			err := block.Insert()
			if err != nil {
				msg := fmt.Sprintf("Failed to insert transaction: %s", err.Error())
				panic(msg)
			}

			for _, txHash := range txs {
				tx := data.GetTx(fmt.Sprintf("%v", txHash))
				beego.Info("Got transaction: ", tx)
				err := tx.Insert()
				if err != nil {
					msg := fmt.Sprintf("Failed to insert transaction: %s", err.Error())
					panic(msg)
				}
			}

			localHgt = localHgt + 1
		}
		//break
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



