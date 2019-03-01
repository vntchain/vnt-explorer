package data

import (
	"github.com/vntchain/vnt-explorer/models"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"github.com/vntchain/vnt-explorer/common"
	"math/big"
	"strings"
	"github.com/vntchain/vnt-explorer/common/utils"
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

	blockNumber := common.Hex(resp.Result.(string)).ToInt64()

	return blockNumber
}

func GetBlock(number int64) (*models.Block, []interface{}, []interface{}) {
	rpc := common.NewRpc()
	rpc.Method = common.Rpc_GetBlockByNumber

	hex := utils.Encode(big.NewInt(number).Bytes())
	if strings.HasPrefix(hex, "0x0") {
		hex = "0x" + hex[3:]
	}

	rpc.Params = append(rpc.Params, hex, false)

	resp := callRpc(rpc)

	blockMap := resp.Result.(map[string]interface{})

	beego.Info("BlockMap: ", blockMap)

	bNumber := common.Hex(blockMap["number"].(string)).ToUint64()

	timestamp := common.Hex(blockMap["timestamp"].(string)).ToUint64()

	size := common.Hex(blockMap["size"].(string)).ToUint64()

	gasUsed := common.Hex(blockMap["gasUsed"].(string)).ToUint64()

	gasLimit := common.Hex(blockMap["gasLimit"].(string)).ToUint64()

	b := &models.Block{
		Number: fmt.Sprintf("%d", bNumber),
		TimeStamp: timestamp,
		Hash: blockMap["hash"].(string),
		ParentHash: blockMap["parentHash"].(string),
		Producer: blockMap["producer"].(string),
		Size: fmt.Sprintf("%d", size),
		GasUsed: gasUsed,
		GasLimit: gasLimit,
		ExtraData: blockMap["extraData"].(string),
	}

	var txs, witnesses []interface{}
	var ok bool
	txIs := blockMap["transactions"].([]interface{})
	beego.Info("txs: ", txIs)
	if txs, ok = blockMap["transactions"].([]interface{}); !ok {
		txs = make([]interface{}, 0)
	}

	if witnesses, ok = blockMap["witnesses"].([]interface{}); !ok {
		witnesses = make([]interface{}, 0)
	}

	b.TxCount = len(txs)

	return b, txs, witnesses
}

func GetTx(txHash string) *models.Transaction {
	rpc := common.NewRpc()
	rpc.Method = common.Rpc_GetTxByHash

	rpc.Params = append(rpc.Params, txHash)

	resp := callRpc(rpc)

	txMap := resp.Result.(map[string]interface{})
	beego.Info("Transaction: ", txMap)

	rpc.Method = common.Rpc_GetTxReceipt

	resp = callRpc(rpc)
	receiptMap := resp.Result.(map[string]interface{})
	beego.Info("Transaction: ", receiptMap)

	tx := &models.Transaction{
		Hash: txMap["hash"].(string),
		From: txMap["from"].(string),
		Value: common.Hex(txMap["value"].(string)).ToString(),
		GasLimit: common.Hex(txMap["gas"].(string)).ToUint64(),
		GasPrice: common.Hex(txMap["gasPrice"].(string)).ToString(),
		GasUsed: common.Hex(receiptMap["gasUsed"].(string)).ToUint64(),
		Nonce: common.Hex(txMap["nonce"].(string)).ToUint64(),
		Index: common.Hex(txMap["transactionIndex"].(string)).ToInt(),
		Input: txMap["input"].(string),
		BlockNumber: common.Hex(txMap["blockNumber"].(string)).ToString(),
	}

	var to string
	var ok bool
	if to, ok = txMap["to"].(string); !ok {
		to = ""
	}

	tx.To = to

	return tx
}