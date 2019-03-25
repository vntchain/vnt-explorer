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
	"github.com/astaxie/beego/orm"
	"strconv"
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
		for _, addr := range addrs {
			tokenBalance := &models.TokenBalance{}
			tokenBalance, err := tokenBalance.GetByAddr(addr, token.Address)
			if err != nil && err == orm.ErrNoRows {
				tokenBalance.Balance = GetAmount(token.Address, addr, tx.BlockNumber)
				tokenBalance.Percent = utils.GetBalancePercent(tokenBalance.Balance, token.TokenAmount, int(token.TokenDecimals))
				 if err = tokenBalance.Insert(); err != nil {
				 	msg := fmt.Sprintf("Failed to insert token balance, token:%s, address:%s, balance:%s",
				 		tokenBalance.Token, tokenBalance.Account, tokenBalance.Balance)
				 	beego.Error(msg)
				 	panic(msg)
				 }

				 currCount, _ := strconv.ParseUint(token.TokenAcctCount, 10, 64)
				 currCount ++
				 token.TokenAcctCount = fmt.Sprintf("%d", currCount)
				 err = token.Update()
				 if err != nil {
				 	msg := fmt.Sprintf("Failed to update token account count, error: %s", err.Error())
				 	beego.Error(msg)
				 	panic(msg)
				 }

				 beego.Info("Success to insert token balance:", tokenBalance.Token, tokenBalance.Account, tokenBalance.Balance)
			} else if err != nil {
				msg := fmt.Sprintf("Failed to get token balance, token:%s, address:%s",
					tokenBalance.Token, tokenBalance.Account)
				beego.Error(msg)
				panic(msg)
			} else if tokenBalance.Id > 0 {
				tokenBalance.Balance = GetAmount(token.Address, addr, tx.BlockNumber)
				tokenBalance.Percent = utils.GetBalancePercent(tokenBalance.Balance, token.TokenAmount, int(token.TokenDecimals))
				if err := tokenBalance.Update(); err != nil {
					msg := fmt.Sprintf("Failed to update token balance, token:%s, address:%s, token:%s",
						tokenBalance.Token, tokenBalance.Account, tokenBalance.Account)
					beego.Error(msg)
					panic(msg)
				}
				beego.Info("Success to update token balance:", tokenBalance.Token, tokenBalance.Account, tokenBalance.Balance)
			}
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

		addrs = append(addrs, tx.From, strings.ToLower(_input.To.String()))

		tx.TokenFrom = tx.From
		tx.TokenTo = _input.To.String()
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

		addrs = append(addrs, tx.From, strings.ToLower(_input.From.String()), strings.ToLower(_input.To.String()))

		tx.TokenFrom = _input.From.String()
		tx.TokenTo = _input.To.String()
		tx.TokenAmount = _input.Value.String()
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

func GetAmount(token, addr string, blockNumber uint64) string {
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