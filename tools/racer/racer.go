package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/tools/racer/data"
	"time"
)

func main() {

	//amount := token.GetMount("0x2b437e35b08ce2d995922f0f07dd67e94ab85b88", "0x122369f04f32269598789998de33e3d56e2c507a", 4201)
	//beego.Info("Amount is: ", amount)
	//supply := token.GetTotalSupply("0x2b437e35b08ce2d995922f0f07dd67e94ab85b88", 4201)
	//beego.Info("supply is: ", supply)
	//dec := token.GetDecimals("0x2b437e35b08ce2d995922f0f07dd67e94ab85b88", 4201)
	//beego.Info("dec is: ", dec)
	//symbol := token.GetSymbol("0x2b437e35b08ce2d995922f0f07dd67e94ab85b88", 4201)
	//beego.Info("symbol is: ", symbol)
	//tokenName := token.GetTokenName("0x2b437e35b08ce2d995922f0f07dd67e94ab85b88", 4201)
	//beego.Info("tokenName is: ", tokenName)
	//
	//return

	rpcHost := beego.AppConfig.String("node::rpc_host")
	rpcPort := beego.AppConfig.String("node::rpc_port")

	beego.Info("rpc host: ", rpcHost)
	beego.Info("rpc port: ", rpcPort)

	for {
		rmtHgt, localHgt := checkHeight()

		//localHgt = 89
		beego.Info(fmt.Sprintf("Local height: %d, rmtHeight: %d", localHgt, rmtHgt))
		if localHgt >= rmtHgt {
			time.Sleep(1 * time.Second)
			continue
		} else {
			nodes := data.GetNodes()
			for _, node := range nodes {
				if err := node.Insert(); err != nil {
					msg := fmt.Sprintf("Failed to insert node: %s", err.Error())
					panic(msg)
				}
			}
		}

		for localHgt < rmtHgt {
			block, txs, witnesses := data.GetBlock(localHgt + 1)

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
				tx.TimeStamp = block.TimeStamp
				beego.Info("Got transaction: ", tx)
				err := tx.Insert()
				if err != nil {
					msg := fmt.Sprintf("Failed to insert transaction: %s", err.Error())
					panic(msg)
				}

				beego.Info("Will extract accounts from transaction: ", txHash)
				data.ExtractAcct(tx)
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
