package token

import (
	"github.com/vntchain/vnt-explorer/models"
	"io/ioutil"
	"github.com/astaxie/beego"
	"github.com/vntchain/go-vnt/accounts/abi"
	"github.com/vntchain/go-vnt/core/wavm"
	"math/big"
	vntCommon "github.com/vntchain/go-vnt/common"
	"fmt"
	"github.com/vntchain/vnt-explorer/common"
	"github.com/vntchain/vnt-explorer/common/utils"
	"strings"
)

var transferSig = map[string]string {
	"0xa9059cbb": "transfer",
	"0x23b872dd": "transferFrom",
}

var abiPath = "./tools/racer/token/erc20.json"

var Abi = readAbi(abiPath)

const (
	TOKEN_ERC20 = 0
)

type Erc20 struct {
	Address 		string
	TokenName		string
	TotalSupply		*big.Int
	Symbol			string
	Decimals		*big.Int
}

func readAbi(abiPath string) abi.ABI {
	beego.Info("Will read abi:", abiPath)
	abiData, err := ioutil.ReadFile(abiPath)
	if err != nil {
		beego.Error("could not read abi: ", "error", err)
		panic(err)
	}

	abi, err := wavm.GetAbi(abiData)
	if err != nil {
		beego.Error("could not read abi: ", "error", err)
	}

	return abi
}


func IsTransfer(tx *models.Transaction) bool {
	input := tx.Input
	if len(input) < 10 {
		return false
	}
	sig := input[0:10]
	if _, ok := transferSig[sig]; ok {
		return true
	}
	return false
}

func UpdateTokenBalance(token *models.Account, tx *models.Transaction) []string {
	if IsTransfer(tx) {

		beego.Info("This is a token transfer transaction:", tx.Hash)
		addrs := GetTransferAddrs(tx)

		tx.IsToken = true
		err := tx.Update()
		if err != nil {
			msg := fmt.Sprintf("Failed to update transaction: %s, error: %s", tx.Hash, err.Error())
			beego.Error(msg)
			panic(msg)
		}
		for _, addr := range addrs {
			PostTokenTask(NewTokenTask(token, addr))
		}
		return addrs[1:]
	}
	return nil
}

func GetTransferAddrs(tx *models.Transaction) (addrs []string) {

	input := tx.Input
	sig := input[0:10]

	//input = input[10:]
	data, err := utils.Decode(input)
	if err != nil {
		msg := fmt.Sprintf("Failed to decode transfer input: %s, error: %s", input, err.Error())
		beego.Error(msg)
		panic(msg)
	}
	data = data[4:]
	switch transferSig[sig] {
	case "transfer":
		type Input struct {
			To 		vntCommon.Address
			Value 	*big.Int
		}

		var _input Input

		err = Abi.UnpackInput(&_input, "transfer", data)

		if err != nil {
			msg := fmt.Sprintf("Failed to unpack input of method: transfer, input: %s, error: %s", data, err.Error())
			beego.Error()
			panic(msg)
		}

		tokenTo := strings.ToLower(_input.To.String())

		addrs = append(addrs, tx.From, tokenTo)

		tx.TokenFrom = tx.From
		tx.TokenTo = tokenTo
		tx.TokenAmount = _input.Value.String()
		break
	case "transferFrom":
		type Input struct {
			From	vntCommon.Address
			To 		vntCommon.Address
			Value 	*big.Int
		}

		var _input Input
		err := Abi.UnpackInput(&_input, "transferFrom", data)
		if err != nil {
			msg := fmt.Sprintf("Failed to unpack input of method: transferFrom, input: %s, error: %s", data, err.Error())
			beego.Error()
			panic(msg)
		}

		tokenFrom := strings.ToLower(_input.From.String())
		tokenTo := strings.ToLower(_input.To.String())
		addrs = append(addrs, tx.From, tokenFrom, tokenTo)

		tx.TokenFrom = tokenFrom
		tx.TokenTo = tokenTo
		tx.TokenAmount = _input.Value.String()
	}

	return
}

func call(token string, data []byte) (*common.Response, *common.Error) {
	dataHex := utils.Encode(data)

	rpc := common.NewRpc()
	rpc.Method = common.Rpc_Call
	rpc.Params = append(rpc.Params, map[string]interface{}{"to": token,
		"gas": utils.EncodeUint64(3000000),
		"data": dataHex}, "latest")

	err, resp, rpcError := utils.CallRpc(rpc)
	if err != nil && rpcError == nil {
		panic(err.Error())
	}
	return resp, rpcError
}

func GetAmount(token, addr string) string {
	data, err := Abi.Pack("GetAmount", vntCommon.HexToAddress(addr))

	if err != nil {
		msg := fmt.Sprintf("Failed to pack input of method: GetAmount: %s", err.Error())
		beego.Error()
		panic(msg)
	}

	resp, rpcErr := call(token, data)

	if rpcErr != nil && rpcErr.Code == -32000 {
		return "0"
	}

	var _out *big.Int

	outData, _ := utils.Decode(resp.Result.(string))
	beego.Info(outData)
	err = Abi.Unpack(&_out, "GetAmount",  outData)

	return _out.String()
}

func GetTotalSupply(token string) *big.Int {
	data, err := Abi.Pack("GetTotalSupply")

	if err != nil {
		msg := fmt.Sprintf("Failed to pack input of method: GetTotalSupply: %s", err.Error())
		beego.Error()
		panic(msg)
	}

	resp, _ := call(token, data)

	var _out *big.Int

	outData, _ := utils.Decode(resp.Result.(string))
	beego.Info(outData)
	err = Abi.Unpack(&_out, "GetTotalSupply",  outData)

	return _out
}

func GetDecimals(token string) *big.Int {
	data, err := Abi.Pack("GetDecimals")

	if err != nil {
		msg := fmt.Sprintf("Failed to pack input of method: GetDecimals: %s", err.Error())
		beego.Error()
		panic(msg)
	}

	resp, _ := call(token, data)

	var _out *big.Int

	outData, _ := utils.Decode(resp.Result.(string))
	beego.Info(outData)
	err = Abi.Unpack(&_out, "GetDecimals",  outData)

	return _out
}

func GetSymbol(token string) string {
	data, err := Abi.Pack("GetSymbol")

	if err != nil {
		msg := fmt.Sprintf("Failed to pack input of method: GetSymbol: %s", err.Error())
		beego.Error()
		panic(msg)
	}

	resp, _ := call(token, data)

	var _out string

	outData, _ := utils.Decode(resp.Result.(string))
	beego.Info(outData)
	err = Abi.Unpack(&_out, "GetSymbol",  outData)

	return _out
}

func GetTokenName(token string) string {
	data, err := Abi.Pack("GetTokenName")

	if err != nil {
		msg := fmt.Sprintf("Failed to pack input of method: GetTokenName: %s", err.Error())
		beego.Error()
		panic(msg)
	}

	resp, _ := call(token, data)

	var _out string

	outData, _ := utils.Decode(resp.Result.(string))
	beego.Info(outData)
	err = Abi.Unpack(&_out, "GetTokenName",  outData)

	return _out
}