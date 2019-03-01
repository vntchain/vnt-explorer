package data

import (
	"github.com/vntchain/vnt-explorer/models"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"github.com/vntchain/vnt-explorer/common"
	"math/big"
	"strings"
)

func GetLocalHeight() int64 {
	b := &models.Block{}
	count, err := b.Count()
	if err != nil {
		msg := fmt.Sprintf("Failed to get block count: %s", err.Error())
		beego.Error(msg)
		panic(msg)
	}

	block, err := b.Last()
	if err != nil {
		msg := fmt.Sprintf("Failed to get last block: %s", err.Error())
		beego.Error(msg)
		panic(msg)
	}

	if block == nil && count > 0 {
		msg := fmt.Sprintf("Block data in db not matched! count %d not equal to lastest block number %d, please check you local database.", count, 0)
		beego.Error(msg)
		panic(msg)
	}

	var bNumber int64

	if block == nil {
		bNumber = 0
	} else {
		bNumber, err = strconv.ParseInt(block.Number, 10, 64)
		if err != nil {
			msg := fmt.Sprintf("Failed to parse block number: %s", err.Error())
			beego.Error(msg)
			panic(msg)
		}
	}

	if bNumber != count {
		msg := fmt.Sprintf("Block data in db not matched! count %d not equal to lastest block number %d, please check you local database.", count, bNumber)
		beego.Error(msg)
		panic(msg)
	}

	return count
}

func GetRemoteHeight() int64 {
	rpc := common.NewRpc()
	rpc.Method = common.Rpc_BlockNumber

	resp := callRpc(rpc)

	beego.Info("Response body", resp)

	blockNumber, err := DecodeBig(resp.Result.(string))
	if err != nil {
		msg := fmt.Sprintf("Failed to decode block number: %s", err.Error())
		beego.Error(msg)
		panic(msg)
	}

	return blockNumber.Int64()
}

func GetBlock(number int64) {
	rpc := common.NewRpc()
	rpc.Method = common.Rpc_GetBlockByNumber

	hex := Encode(big.NewInt(number).Bytes())
	if strings.HasPrefix(hex, "0x0") {
		hex = "0x" + hex[3:]
	}

	rpc.Params = append(rpc.Params, hex, false)

	resp := callRpc(rpc)

	blockMap := resp.Result.(map[string]interface{})
	beego.Info("Block: ", blockMap)

	bNumber, err := DecodeUint64(blockMap["number"].(string))
	if err != nil {
		msg := fmt.Sprintf("Failed to decode block number: %s", err.Error())
		beego.Error(msg)
		panic(msg)
	}

	timestamp, err := DecodeUint64(blockMap["timestamp"].(string))
	if err != nil {
		msg := fmt.Sprintf("Failed to decode timestamp: %s", err.Error())
		beego.Error(msg)
		panic(msg)
	}

	size, err := DecodeUint64(blockMap["size"].(string))
	if err != nil {
		msg := fmt.Sprintf("Failed to decode size: %s", err.Error())
		beego.Error(msg)
		panic(msg)
	}

	gasUsed, err := DecodeUint64(blockMap["gasUsed"].(string))
	if err != nil {
		msg := fmt.Sprintf("Failed to decode gasUsed: %s", err.Error())
		beego.Error(msg)
		panic(msg)
	}

	gasLimit, err := DecodeUint64(blockMap["gasLimit"].(string))
	if err != nil {
		msg := fmt.Sprintf("Failed to decode gasLimit: %s", err.Error())
		beego.Error(msg)
		panic(msg)
	}

	b := &models.Block{
		Number: string(bNumber),
		TimeStamp: timestamp,
		Hash: blockMap["hash"].(string),
		ParentHash: blockMap["parentHash"].(string),
		Producer: blockMap["producer"].(string),
		Size: string(size),
		GasUsed: gasUsed,
		GasLimit: gasLimit,
		ExtraData: blockMap["extraData"].(string),
	}

	var txs, witnesses []string
	var ok bool
	if txs, ok = blockMap["transactions"].([]string); !ok {
		txs = make([]string, 0)
	}

	if witnesses, ok = blockMap["witnesses"].([]string); !ok {
		witnesses = make([]string, 0)
	}

	b.TxCount = len(txs)

	beego.Info("Block: ", b)
	beego.Info("Transactions: ", txs)
	beego.Info("witnesses: ", witnesses)
}