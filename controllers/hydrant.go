package controllers

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/bluele/gcache"
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
	hydrantCache                 = gcache.New(10000).LRU().Build()
)

var prv *ecdsa.PrivateKey

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
		this.ReturnErrorMsg("Wrong format of Address: %s, it must be 0x[0-9|a-z|A-Z]* or [0-9|a-z|A-Z]*", addr.Address)
		return
	}
	address := vntCommon.HexToAddress(addr.Address)

	// get account from hydrant db
	hydrant := getHydrant(address.String())
	now := time.Now().Unix()
	if hydrant != nil {
		// check interval
		if now-hydrant.TimeStamp < int64(hydrantInterval) {
			this.ReturnErrorMsg("Too Frequently: the last time sent vnt to you is %s", strconv.FormatInt(hydrant.TimeStamp, 10))
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
		this.ReturnErrorMsg("System error: %s. Please contract developers.", err.Error())
		return
	}

	// build transaction
	amount := big.NewInt(1).Mul(big.NewInt(int64(hydrantCount)), big.NewInt(1e18))
	tx := vntTypes.NewTransaction(nonce, address, amount, common.DefaultGasLimit, big.NewInt(common.DefaultGasPrice), nil)
	transaction, err := vntTypes.SignTx(tx, vntTypes.HomesteadSigner{}, prv)
	data, err := vntRlp.EncodeToBytes(transaction)
	if err != nil {
		this.ReturnErrorMsg("System error: %s. Please contract developers.", err.Error())
		return
	}
	dataStr := fmt.Sprintf("0x%x", data)

	// call rpc
	rpc := common.NewRpc()
	rpc.Method = common.Rpc_SendRawTransaction
	rpc.Params = append(rpc.Params, dataStr)

	err, resp := utils.CallRpc(rpc)
	if err != nil {
		this.ReturnErrorMsg("System error: %s. Please contract developers.", err.Error())
		return
	}

	// TODO check transaction status
	txHash := resp.Result.(string)
	updateHydrant(address.String(), hydrant)
	this.ReturnData(txHash)
}

// get data from cache or db
func getHydrant(addr string) *models.Hydrant {
	addr = strings.ToLower(addr)
	if _type, err := hydrantCache.Get(addr); err == nil && _type != nil {
		beego.Info("Address hit in cache:", addr)
		return _type.(*models.Hydrant)
	} else {
		beego.Info("Address not hit in cache:", addr)
		a := &models.Hydrant{}
		a, err := a.Get(addr)
		if err != nil {
			beego.Info("Address not hit in db:", addr)
			return nil
		}
		beego.Info("Address hit in db:", addr)
		hydrantCache.Set(addr, a)
		return a
	}
}

// update db and cache
func updateHydrant(addr string, hydrant *models.Hydrant) {
	addr = strings.ToLower(addr)
	if err := hydrant.InsertOrUpdate(); err != nil {
		msg := fmt.Sprintf("Failed to update account: %s, error: %s", addr, err.Error())
		beego.Error(msg)
		panic(err)
	}
	hydrantCache.Set(addr, hydrant)
}

func getNonce(addr string) (uint64, error) {
	addr = vntCommon.HexToAddress(addr).String()
	rpc := common.NewRpc()
	rpc.Method = common.Rpc_GetTransactionCount
	rpc.Params = append(rpc.Params, addr)
	rpc.Params = append(rpc.Params, "latest")
	err, resp := utils.CallRpc(rpc)
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
