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
	sig := input[0:10]
	if _, ok := transferSig[sig]; ok {
		return true
	}
	return false
}

func GetTransferAddrs(tx *models.Transaction) (addrs []string) {

	input := tx.Input
	sig := input[0:10]

	input = input[10:]
	switch transferSig[sig] {
	case "transfer":
		type Input struct {
			To 		vntCommon.Address
			value 	big.Int
		}

		var _input Input
		err := Abi.UnpackInput(&_input, "transfer", []byte(input))

		if err != nil {
			msg := fmt.Sprintf("Failed to unpack input of method: transfer, input: %s, error: %s", input, err.Error())
			beego.Error()
			panic(msg)
		}

		addrs = append(addrs, tx.From)

		addrs = append(addrs, _input.To.String())
		break
	case "transferFrom":
		type Input struct {
			From	vntCommon.Address
			To 		vntCommon.Address
			value 	big.Int
		}

		var _input Input
		err := Abi.UnpackInput(&_input, "transferFrom", []byte(input))
		if err != nil {
			msg := fmt.Sprintf("Failed to unpack input of method: transferFrom, input: %s, error: %s", input, err.Error())
			beego.Error()
			panic(msg)
		}

		addrs = append(addrs, _input.From.String(), _input.To.String())
	}

	return
}

func call(token string, blockNumber uint64, data []byte) *common.Response {
	dataHex := utils.Encode(data)

	rpc := common.NewRpc()
	rpc.Method = common.Rpc_Call
	rpc.Params = append(rpc.Params, map[string]interface{}{"to": token,
		"gas": utils.EncodeUint64(3000000),
		"data": dataHex},
		utils.EncodeUint64(blockNumber))

	err, resp := utils.CallRpc(rpc)
	if err != nil {
		panic(err.Error())
	}
	return resp
}

func GetMount(token, addr string, blockNumber uint64) string {
	data, err := Abi.Pack("GetAmount", vntCommon.HexToAddress(addr))

	if err != nil {
		msg := fmt.Sprintf("Failed to pack input of method: GetAmount: %s", err.Error())
		beego.Error()
		panic(msg)
	}

	resp := call(token, blockNumber, data)

	var _out *big.Int

	outData, _ := utils.Decode(resp.Result.(string))
	beego.Info(outData)
	err = Abi.Unpack(&_out, "GetAmount",  outData)

	return _out.String()
}

func GetTotalSupply(token string, blockNumber uint64) *big.Int {
	data, err := Abi.Pack("GetTotalSupply")

	if err != nil {
		msg := fmt.Sprintf("Failed to pack input of method: GetTotalSupply: %s", err.Error())
		beego.Error()
		panic(msg)
	}

	resp := call(token, blockNumber, data)

	var _out *big.Int

	outData, _ := utils.Decode(resp.Result.(string))
	beego.Info(outData)
	err = Abi.Unpack(&_out, "GetTotalSupply",  outData)

	return _out
}

func GetDecimals(token string, blockNumber uint64) *big.Int {
	data, err := Abi.Pack("GetDecimals")

	if err != nil {
		msg := fmt.Sprintf("Failed to pack input of method: GetDecimals: %s", err.Error())
		beego.Error()
		panic(msg)
	}

	resp := call(token, blockNumber, data)

	var _out *big.Int

	outData, _ := utils.Decode(resp.Result.(string))
	beego.Info(outData)
	err = Abi.Unpack(&_out, "GetDecimals",  outData)

	return _out
}

func GetSymbol(token string, blockNumber uint64) string {
	data, err := Abi.Pack("GetSymbol")

	if err != nil {
		msg := fmt.Sprintf("Failed to pack input of method: GetSymbol: %s", err.Error())
		beego.Error()
		panic(msg)
	}

	resp := call(token, blockNumber, data)

	var _out string

	outData, _ := utils.Decode(resp.Result.(string))
	beego.Info(outData)
	err = Abi.Unpack(&_out, "GetSymbol",  outData)

	return _out
}

func GetTokenName(token string, blockNumber uint64) string {
	data, err := Abi.Pack("GetTokenName")

	if err != nil {
		msg := fmt.Sprintf("Failed to pack input of method: GetTokenName: %s", err.Error())
		beego.Error()
		panic(msg)
	}

	resp := call(token, blockNumber, data)

	var _out string

	outData, _ := utils.Decode(resp.Result.(string))
	beego.Info(outData)
	err = Abi.Unpack(&_out, "GetTokenName",  outData)

	return _out
}