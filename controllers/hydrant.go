package controllers

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego"
	vntCommon "github.com/vntchain/go-vnt/common"
	vntTypes "github.com/vntchain/go-vnt/core/types"
	vntCrypto "github.com/vntchain/go-vnt/crypto"
	vntRlp "github.com/vntchain/go-vnt/rlp"
	"github.com/vntchain/vnt-explorer/common"
	"github.com/vntchain/vnt-explorer/common/utils"
	"github.com/vntchain/vnt-explorer/models"
)

var (
	hydrantInterval, intervalErr = beego.AppConfig.Int("hydrant::interval")
	hydrantCount, countErr       = beego.AppConfig.Int("hydrant::count")
	hydrantFrom                  = beego.AppConfig.String("hydrant::from")
	hydrantPrivteKey             = beego.AppConfig.String("hydrant::privateKey")
	hydrantChainId, chainIdErr   = beego.AppConfig.Int("hydrant::chainId")
	addrMap                      = make(map[string]interface{})
)

var prv *ecdsa.PrivateKey
var mutex sync.Mutex

type HydrantController struct {
	BaseController
}

func getConfig() {
	if intervalErr != nil {
		hydrantInterval = common.DefaultHydrantInterval
		intervalErr = nil
	}

	if countErr != nil {
		hydrantCount = common.DefaultHydrantCount
		countErr = nil
	}

	if chainIdErr != nil{
		hydrantChainId = common.DefaultHydrantChainId
		chainIdErr = nil
	}
	var err error
	prv, err = vntCrypto.HexToECDSA(hydrantPrivteKey)
	if err != nil {
		beego.Error("PrivteKey  conf of system account is error: ", err)
	}
}

func (this *HydrantController) SendVnt() {
	// get config
	getConfig()
	if prv == nil {
		this.ReturnErrorMsg("System error!", "")
		return
	}

	type Addr struct {
		Address string
	}
	var addr Addr
	body := this.Ctx.Input.RequestBody
	err := json.Unmarshal(body, &addr)
	if err != nil {
		this.ReturnErrorMsg("Wrong format of Address: %s", err.Error())
		return
	}
	if !isHex(addr.Address) {
		this.ReturnErrorMsg("Wrong format of Address: %s, it must be 0x[0-9|a-f|A-F]* or [0-9|a-f|A-F]*", addr.Address)
		return
	}
	address := vntCommon.HexToAddress(addr.Address)

	// 如果有交易正在向该账户发送vnt，则返回
	mutex.Lock()
	if _, exists := addrMap[address.String()]; exists {
		defer mutex.Unlock()
		this.ReturnErrorMsg("Sending vnt to %s! Please wait for a moment! ", addr.Address)
		return
	}
	addrMap[address.String()] = nil
	mutex.Unlock()

	// get account from hydrant db
	hydrant := getHydrant(address.String())
	now := time.Now().Unix()
	if hydrant != nil {
		// check interval
		if now-hydrant.TimeStamp < int64(hydrantInterval) {
			lastTime := time.Unix(hydrant.TimeStamp, 0)
			this.ReturnErrorMsg("Too Frequently: the last time sent vnt to you is %s", lastTime.Format("2006-01-02 15:04:05"))
			deleteAddrMap(address.String())
			return
		}
	}
	hydrant = &models.Hydrant{
		Address:   address.String(),
		TimeStamp: now,
	}

	// get nonce
	nonce, err := getNonce(hydrantFrom)
	if err != nil {
		this.ReturnErrorMsg("System error, get nonce: %s. Please contract developers.", err.Error())
		deleteAddrMap(address.String())
		return
	}

	// build transaction
	amount := big.NewInt(1).Mul(big.NewInt(int64(hydrantCount)), big.NewInt(1e18))
	tx := vntTypes.NewTransaction(nonce, address, amount, common.DefaultGasLimit, big.NewInt(common.DefaultGasPrice), nil)
	transaction, err := vntTypes.SignTx(tx, vntTypes.NewHubbleSigner(big.NewInt(int64(hydrantChainId))), prv)
	data, err := vntRlp.EncodeToBytes(transaction)
	if err != nil {
		this.ReturnErrorMsg("System error, signTx: %s. Please contract developers.", err.Error())
		deleteAddrMap(address.String())
		return
	}
	dataStr := fmt.Sprintf("0x%x", data)

	// call rpc
	rpc := common.NewRpc()
	rpc.Method = common.Rpc_SendRawTransaction
	rpc.Params = append(rpc.Params, dataStr)

	err, resp, _ := utils.CallRpc(rpc)
	if err != nil {
		this.ReturnErrorMsg("System error, sendRawTransaction %s. Please contract developers.", err.Error())
		deleteAddrMap(address.String())
		return
	}

	txHash := resp.Result.(string)
	// check transaction status
	go checkTxReceipt(txHash, address.String(), hydrant)
	this.ReturnData(txHash, nil)
}

// get data from cache or db
func getHydrant(addr string) *models.Hydrant {
	addr = strings.ToLower(addr)
	a := &models.Hydrant{}
	a, err := a.Get(addr)
	if err != nil {
		return nil
	}
	return a
}

// update db and cache
func updateHydrant(addr string, hydrant *models.Hydrant) {
	addr = strings.ToLower(addr)
	if err := hydrant.InsertOrUpdate(); err != nil {
		msg := fmt.Sprintf("Failed to update account: %s, error: %s", addr, err.Error())
		beego.Error(msg)
		panic(err)
	}
}

func getNonce(addr string) (uint64, error) {
	addr = vntCommon.HexToAddress(addr).String()
	rpc := common.NewRpc()
	rpc.Method = common.Rpc_GetTransactionCount
	rpc.Params = append(rpc.Params, addr)
	rpc.Params = append(rpc.Params, "latest")
	err, resp, _ := utils.CallRpc(rpc)
	if err != nil {
		return 0, err
	}
	nonce := utils.Hex(resp.Result.(string)).ToUint64()
	return nonce, nil
}

// isHexCharacter returns bool of c being a valid hexadecimal.
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// hasHexPrefix validates str begins with '0x' or '0X'.
func hasHexPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// isHex validates whether each byte is valid hexadecimal string.
func isHex(str string) bool {
	if hasHexPrefix(str) {
		str = str[2:]
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}

func checkTxReceipt(txHash string, addr string, hydrant *models.Hydrant) bool {
	defer deleteAddrMap(addr)
	// query 20 times at most
	for i := 0; i < 20; i++ {
		rpc := common.NewRpc()
		rpc.Method = common.Rpc_GetTxReceipt
		rpc.Params = append(rpc.Params, txHash)
		err, resp, _ := utils.CallRpc(rpc)
		if err != nil || resp.Result == nil {
			time.Sleep(time.Second)
			continue
		}

		receiptMap := resp.Result.(map[string]interface{})
		if utils.Hex(receiptMap["status"].(string)).ToUint64() == 1 {
			updateHydrant(addr, hydrant)
			return true
		} else {
			return false
		}
	}
	return false
}

func deleteAddrMap(addr string) {
	mutex.Lock()
	delete(addrMap, addr)
	mutex.Unlock()
}
