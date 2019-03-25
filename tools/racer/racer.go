package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/tools/racer/data"
	"github.com/vntchain/vnt-explorer/models"
	"strings"
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
	//tb := models.TokenBalance{}
	//ta, err := tb.GetByAddr("hello", "there")
	//beego.Info(ta, err == orm.ErrNoRows)
	//return

	for {
		doSync()
	}
}

func doSync() {
	defer func() {
		if r := recover(); r != nil {
			beego.Error("Error happened:", r)
			time.Sleep(2 * time.Second)
		}
	}()

	rmtHgt, localHgt, lastBlock := checkHeight()

	// localHgt = 14
	//rmtHgt = 25913
	beego.Info(fmt.Sprintf("Local height: %d, rmtHeight: %d", localHgt, rmtHgt))
	if localHgt >= rmtHgt {
		beego.Info("here!")
		time.Sleep(2 * time.Second)
		return
	}

	var block *models.Block
	var txs, witnesses []interface{}
	var leftAddrs []string

	// Set the block sync batch to 1000
	if rmtHgt - localHgt > 1000 {
		rmtHgt = localHgt + 1000
	}

	for localHgt < rmtHgt {
		block, txs, witnesses = data.GetBlock(localHgt + 1)

		beego.Info("Block:", block)
		beego.Info("txs:", txs)
		beego.Info("witness:", witnesses)

		leftAddrs = make([]string, 0)
		for _, w := range witnesses {
			leftAddrs = append(leftAddrs, fmt.Sprintf("%v", w))
		}

		var dynamicReward float64

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

			tmp, err := strconv.Atoi(tx.GasPrice)
			if err != nil {
				msg := fmt.Sprintf("Failed to convert gasPrice: %s", err.Error())
				panic(msg)
			}
			dynamicReward += float64(tx.GasUsed) * (float64(tmp) / 1e18)
			// beego.Info("---------> tx.GasUsed == ", tx.GasUsed, "tx.GasPrice ==", tx.GasPrice)
		}

		block.TxCount = len(txs)
		time := float32(2.0)
		if lastBlock != nil {
			time = float32(block.TimeStamp - lastBlock.TimeStamp)
		}
		block.Tps = float32(block.TxCount) / time

		// Persist witness accounts and other unknown accounts from token transfer
		data.PersistWitnesses(leftAddrs, block.Number)

		// compute blockReward
		// 区块奖励，0-47304000都是6个vnt，47304001-94608000是3个，再之后是1.5个
		var staticReward float64
		if block.Number >= 0 && block.Number <= 47304000 {
			staticReward = 6
		} else if block.Number >= 47304001 && block.Number <= 94608000 {
			staticReward = 3
		} else if block.Number >= 94608001 {
			staticReward = 1.5
		}

		result := staticReward + dynamicReward
		block.BlockReward = fmt.Sprintf("%f VNT (%f + %f)", result, staticReward, dynamicReward)
		// beego.Info("---------> block.BlockReward == ", block.BlockReward)
		err := block.Insert()
		if err != nil {
			msg := fmt.Sprintf("Failed to insert block: %s", err.Error())
			panic(msg)
		}

		localHgt = localHgt + 1
	}

	witMap := make(map[string]int)
	//fmt.Println("witnesses: %v", leftAddrs)
	for _, addr := range leftAddrs {
		witMap[strings.ToLower(addr)] = 1
	}

	nodes := data.GetNodes()
	for _, node := range nodes {
		//fmt.Println("node address: %s", node.Address)
		if witMap[node.Address] == 1 {
			node.IsSuper = 1
		} else {
			node.IsSuper = 0
		}
		if err := node.Insert(); err != nil {
			msg := fmt.Sprintf("Failed to insert node: %s", err.Error())
			panic(msg)
		}
	}

}

func checkHeight() (int64, int64, *models.Block) {
	rmtHgt := data.GetRemoteHeight()
	localHgt, lastBlock := data.GetLocalHeight()

	if localHgt > rmtHgt {
		msg := fmt.Sprintf("Local height %d is bigger than remote height: %d, please check your remote node", localHgt, rmtHgt)
		beego.Error(msg)
		panic(msg)
	}

	return rmtHgt, localHgt, lastBlock
}
