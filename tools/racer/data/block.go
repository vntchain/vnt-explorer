package data

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/astaxie/beego"
	"github.com/bluele/gcache"
	"github.com/vntchain/vnt-explorer/common"
	"github.com/vntchain/vnt-explorer/common/utils"
	"github.com/vntchain/vnt-explorer/models"
	"github.com/vntchain/vnt-explorer/tools/racer/token"
	"sync"
	"strconv"
	"github.com/astaxie/beego/orm"
	"os"
	"encoding/csv"
)
const BatchSize = 30

var AccountMap = sync.Map{}

type TX struct {
	BlockHash	string		`json:"blockHash"`
	BlockNumber string		`json:"blockNumber"`
	From		string		`json:"from"`
	Gas			string		`json:"gas"`
	GasPrice	string		`json:"gasPrice"`
	Hash		string		`json:"hash"`
	Input		string		`json:"input"`
	Nonce		string		`json:"Nonce"`
	To			string		`json:"to"`
	Index		string		`json:"transactionIndex"`
	Value		string		`json:"value"`
}

type Receipt struct {
	GasUsed		string		`json:"gasUsed"`
	Status		string		`json:"status"`
	ContractAddr string		`json:"contractAddress"`
}

var acctCache = gcache.New(10000).LRU().Build()

const (
	ACC_TYPE_NULL     = 0
	ACC_TYPE_NORMAL   = 1
	ACC_TYPE_CONTRACT = 2
	ACC_TYPE_TOKEN    = 3
)

func GetLocalHeight() (int64, *models.Block) {
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
		bNumber = -1
	} else {
		bNumber = int64(block.Number)
	}

	if bNumber + 1 != count {
		msg := fmt.Sprintf("Block data in db not matched! count %d not match lastest block number %d.", count, bNumber)
		beego.Warn(msg)
		c, b := SearchValidHeight(count, bNumber)
		beego.Info(">>> Using", c, "as local height.")
		return c, b
	}

	return count-1, block
}

func SearchValidHeight(currCount int64, topNumber int64) (int64, *models.Block) {
	beego.Info("Search for a valid block...")
	left := uint64(topNumber - BatchSize)
	right := uint64(currCount-1)
	block := &models.Block{}
	for left < right {
		mid := (left + right) / 2
		b, err := block.GetByNumber(mid)
		if err != nil {
			if err == orm.ErrNoRows {
				right = block.Number - 1
				continue
			} else {
				msg := fmt.Sprintf("Failed to get block by number: %d, err: %s", mid, err.Error())
				beego.Error(msg)
				panic(msg)
			}
		}

		c, err := block.CountBellow(mid)

		beego.Info(">>>> Left:", left, "Right:", right, "Mid:", mid, "Count: ", c)
		if err != nil {
			msg := fmt.Sprintf("Failed to get count bellow number: %d", mid)
			beego.Error(msg)
			panic(msg)
		}

		if int64(b.Number) + 1 == c {
			beego.Info(">>>> A good block:", b.Number)
			if right - b.Number <= 1 {

				beego.Info(">>>> A bingo block:", b.Number)
				return int64(b.Number), b
			} else {
				beego.Info(">>>> Use as left: block:", b.Number)
				left = b.Number
			}
		} else {
			beego.Info(">>>> A bad block:", b.Number, ", using as right")
			right = b.Number - 1
		}
	}
	return int64(left), block
}

func GetRemoteHeight() int64 {
	rpc := common.NewRpc()
	rpc.Method = common.Rpc_BlockNumber

	err, resp, _ := utils.CallRpc(rpc)
	if err != nil {
		panic(err.Error())
	}

	beego.Info("Response body", resp)

	blockNumber := utils.Hex(resp.Result.(string)).ToInt64()

	return blockNumber
}

func GetBlockMap(number int64) map[string]interface{} {
	rpc := common.NewRpc()
	rpc.Method = common.Rpc_GetBlockByNumber

	if number >= 0 {
		hex := utils.Encode(big.NewInt(number).Bytes())
		if strings.HasPrefix(hex, "0x0") {
			hex = "0x" + hex[3:]
		}

		if hex == "0x" {
			hex = "0x0"
		}

		rpc.Params = append(rpc.Params, hex, true)
	} else {
		rpc.Params = append(rpc.Params, "latest", true)
	}


	err, resp, _ := utils.CallRpc(rpc)
	if err != nil {
		panic(err.Error())
	}

	blockMap := resp.Result.(map[string]interface{})
	return blockMap
}

func GetBlock(number int64) (*models.Block, []map[string]interface{}, []interface{}) {

	blockMap := GetBlockMap(number)

	beego.Info("BlockMap: ", blockMap)

	bNumber := utils.Hex(blockMap["number"].(string)).ToUint64()

	timestamp := utils.Hex(blockMap["timestamp"].(string)).ToUint64()

	size := utils.Hex(blockMap["size"].(string)).ToUint64()

	gasUsed := utils.Hex(blockMap["gasUsed"].(string)).ToUint64()

	gasLimit := utils.Hex(blockMap["gasLimit"].(string)).ToUint64()

	b := &models.Block{
		Number:     bNumber,
		TimeStamp:  timestamp,
		Hash:       blockMap["hash"].(string),
		ParentHash: blockMap["parentHash"].(string),
		Producer:   strings.ToLower(blockMap["producer"].(string)),
		Size:       fmt.Sprintf("%d", size),
		GasUsed:    gasUsed,
		GasLimit:   gasLimit,
		ExtraData:  blockMap["extraData"].(string),
	}

	var witnesses []interface{}
	var txs []map[string]interface{}
	var ok bool
	txIs := blockMap["transactions"].([]interface{})
	beego.Info("txs: ", txIs)

	if len(txIs) > 0 {
		for _, tx := range txIs {
			txs = append(txs, tx.(map[string]interface{}))
		}
	}

	//txs = blockMap["transactions"].([]map[string]interface{})

	//if txs, ok = blockMap["transactions"].([]map[string]interface{}); !ok {
	//	beego.Info("Failed to to get txs", txs)
	//
	//	txs = make([]map[string]interface{}, 0)
	//}

	if witnesses, ok = blockMap["witnesses"].([]interface{}); !ok {
		witnesses = make([]interface{}, 0)
	}

	return b, txs, witnesses
}

func GetWitnesses(number int64) ([]interface{}) {

	blockMap := GetBlockMap(number)

	var ok bool
	var witnesses []interface{}

	if witnesses, ok = blockMap["witnesses"].([]interface{}); !ok {
		witnesses = make([]interface{}, 0)
	}

	return witnesses
}

func GetLastBlock(number int64) (*models.Block) {
	blockMap := GetBlockMap(number)

	beego.Info("BlockMap: ", blockMap)

	bNumber := utils.Hex(blockMap["number"].(string)).ToUint64()

	timestamp := utils.Hex(blockMap["timestamp"].(string)).ToUint64()

	size := utils.Hex(blockMap["size"].(string)).ToUint64()

	gasUsed := utils.Hex(blockMap["gasUsed"].(string)).ToUint64()

	gasLimit := utils.Hex(blockMap["gasLimit"].(string)).ToUint64()

	b := &models.Block{
		Number:     bNumber,
		TimeStamp:  timestamp,
		Hash:       blockMap["hash"].(string),
		ParentHash: blockMap["parentHash"].(string),
		Producer:   strings.ToLower(blockMap["producer"].(string)),
		Size:       fmt.Sprintf("%d", size),
		GasUsed:    gasUsed,
		GasLimit:   gasLimit,
		ExtraData:  blockMap["extraData"].(string),
	}

	return b
}

func PersistBlock(number int64) {
	beego.Info("Will get block: ", number)
	block, txs, witnesses := GetBlock(number)
	var lastBlock *models.Block = nil
	if number > 0 {
		lastBlock = GetLastBlock(number-1)
	}

	beego.Info("Block:", block)
	beego.Info("txs:", txs)
	beego.Info("witness:", witnesses)

	leftAddrs := make([]string, 0)
	for _, w := range witnesses {
		leftAddrs = append(leftAddrs, fmt.Sprintf("%v", w))
	}

	var dynamicReward float64

	for _, tx := range txs {
		tx := GetTx(tx)
		tx.TimeStamp = block.TimeStamp
		beego.Info("Got transaction: ", tx)

		//PostTxTask(NewTxTask(tx))

		beego.Info("Will extract accounts from transaction: ", tx.Hash)
		PostExtractAccountTask(NewExtractAccountTask(tx))
		//data.ExtractAcct(tx)

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
	//data.PersistWitnesses(leftAddrs, block.Number)
	PostWitnessesTask(NewWitnessesTask(leftAddrs, block.Number))

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
	block.Reward = staticReward
	block.Fee = dynamicReward
	// beego.Info("---------> block.BlockReward == ", block.BlockReward)

	PostInsertBlockTask(NewBlockInsertTask(block))
}


func GetTx(txMap map[string]interface{}) *models.Transaction {
	rpc := common.NewRpc()
	//rpc.Method = common.Rpc_GetTxByHash
	//
	//rpc.Params = append(rpc.Params, txHash)
	//
	//err, resp, _ := utils.CallRpc(rpc)
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//txMap := resp.Result.(map[string]interface{})
	//beego.Info("Transaction: ", txMap)

	txHash := txMap["hash"].(string)

	rpc.Method = common.Rpc_GetTxReceipt
	rpc.Params = append(rpc.Params, txHash)

	err, resp, _ := utils.CallRpc(rpc)
	if err != nil {
		panic(err.Error())
	}

	receiptMap := resp.Result.(map[string]interface{})
	beego.Info("Transaction: ", receiptMap)

	tx := &models.Transaction{
		Hash:        txHash,
		From:        strings.ToLower(txMap["from"].(string)),
		Value:       utils.Hex(txMap["value"].(string)).ToString(),
		GasLimit:    utils.Hex(txMap["gas"].(string)).ToUint64(),
		GasPrice:    utils.Hex(txMap["gasPrice"].(string)).ToString(),
		GasUsed:     utils.Hex(receiptMap["gasUsed"].(string)).ToUint64(),
		Nonce:       utils.Hex(txMap["nonce"].(string)).ToUint64(),
		Index:       utils.Hex(txMap["transactionIndex"].(string)).ToInt(),
		Input:       txMap["input"].(string),
		Status:      utils.Hex(receiptMap["status"].(string)).ToInt(),
		BlockNumber: utils.Hex(txMap["blockNumber"].(string)).ToUint64(),
	}

	var to string
	var ok bool
	if to, ok = txMap["to"].(string); !ok {
		to = ""

		beego.Info("This is a transaction of contract creation.")
		if contractAddr, ok := receiptMap["contractAddress"].(string); ok {
			tx.ContractAddr = strings.ToLower(contractAddr)
		}
		tx.To = nil
	} else {
		tx.To = &models.Account{Address: strings.ToLower(to)}
	}

	return tx
}

// Extract Account from a transaction
func ExtractAcct(tx *models.Transaction) {
	from := tx.From
	to := tx.To
	contractAddr := tx.ContractAddr

	if _, ok := AccountMap.Load(from); !ok {
		AccountMap.Store(from, 0)
		if a := GetAccount(from); a == nil {
			beego.Info("Block:", tx.BlockNumber, ", will insert normal account:", from)
			NewAccount(from, tx, ACC_TYPE_NORMAL, 1)
		} else {
			beego.Info("Block:", tx.BlockNumber, ", will update normal account:", from)
			UpdateAccount(a, tx, ACC_TYPE_NORMAL, 1)
		}
	}

	if to != nil && to.Address != "" {
		if a := GetAccount(to.Address); a == nil {
			beego.Info("Block:", tx.BlockNumber, ", will insert normal account:", to)
			NewAccount(to.Address, tx, ACC_TYPE_NORMAL, 1)
		} else {
			if a.IsToken {
				beego.Info("Block:", tx.BlockNumber, ", will update token account:", to)
				UpdateAccount(a, tx, ACC_TYPE_TOKEN, 1)
			} else if a.IsContract {
				beego.Info("Block:", tx.BlockNumber, ", will update contract account:", to)
				UpdateAccount(a, tx, ACC_TYPE_CONTRACT, 1)
			} else {
				if _, ok := AccountMap.Load(to.Address); !ok {
					AccountMap.Store(to.Address, 0)
					beego.Info("Block:", tx.BlockNumber, ", will update normal account:", from)
					UpdateAccount(a, tx, ACC_TYPE_NORMAL, 1)
				}
			}
		}
	} else if contractAddr != "" { // this case is for contract creation
		if a := GetAccount(contractAddr); a == nil {
			// new contract account
			beego.Info("Block:", tx.BlockNumber, ", will insert contract account:", contractAddr)
			NewAccount(contractAddr, tx, ACC_TYPE_CONTRACT, 0)
		} else if !a.IsContract {
			// this account already exists as a normal account,
			// will change it to a contract account
			//a.IsContract = true
			beego.Info("Block:", tx.BlockNumber, ", will update contract account:", contractAddr)
			UpdateAccount(a, tx, ACC_TYPE_CONTRACT, 0)
		}
	}

	PostTxTask(NewTxTask(tx))
	return
}

func GetBalance(addr string) string {
	rpc := common.NewRpc()
	rpc.Method = common.Rpc_GetBlance

	rpc.Params = append(rpc.Params, addr)
	rpc.Params = append(rpc.Params, "latest")

	err, resp, rpcError := utils.CallRpc(rpc)
	if err != nil && rpcError == nil {
		panic(err.Error())
	} else if rpcError != nil && rpcError.Code == -32000 {
		return "0"
	}

	balance := utils.Hex(resp.Result.(string)).ToString()
	return balance
}

func IsToken(addr string, tx *models.Transaction) (bool, *token.Erc20) {
	totalSupply := token.GetTotalSupply(addr)
	tokenName := token.GetTokenName(addr)
	decimals := token.GetDecimals(addr)
	symbol := token.GetSymbol(addr)

	if totalSupply != nil && decimals != nil && symbol != "" && tokenName != "" {
		erc20 := &token.Erc20{
			Address:     addr,
			TokenName:   tokenName,
			TotalSupply: totalSupply,
			Symbol:      symbol,
			Decimals:    decimals,
		}

		//// Update the tx
		//tx.IsToken = true
		//err := tx.Update()
		//if err != nil {
		//	msg := fmt.Sprintf("Failed to update transaction: %s, error: %s", tx.Hash, err.Error())
		//	beego.Error(msg)
		//	panic(msg)
		//}

		return true, erc20
	}

	return false, nil
}

// Insert a new Account, in this case, tye _type only could be "normal" or "contract"
func NewAccount(addr string, tx *models.Transaction, _type int, txCount uint64) {
	a := &models.Account{
		Address:        addr,
		Vname:          addr, //todo: get vname
		Balance:        "0",
		TxCount:        txCount,
		FirstBlock:     tx.BlockNumber,
		LastBlock:      tx.BlockNumber,
		TokenAmount:    "0",
		TokenAcctCount: "0",
		InitTx:         tx.Hash,
		LastTx:         tx.Hash,
	}

	if _type == ACC_TYPE_CONTRACT {
		a.IsContract = true
		a.ContractName = "" // TODO: extract contract name from contract code
		a.ContractOwner = tx.From

		//beego.Info("######### a.ContractOwner:", a.ContractOwner)

		if ok, erc20 := IsToken(addr, tx); ok {
			a.IsToken = true

			// TODO: get token detail by calling the contract
			a.TokenType = token.TOKEN_ERC20
			a.ContractName = erc20.TokenName
			a.TokenAmount = erc20.TotalSupply.String()
			a.TokenSymbol = erc20.Symbol
			a.TokenDecimals = erc20.Decimals.Uint64()
			a.TokenAcctCount = "0"
			a.TokenLogo = ""
		}
	}

	a.Balance = GetBalance(addr)
	a.Percent = utils.GetBalancePercent(a.Balance, common.VNT_TOTAL, common.VNT_DECIMAL)
	insertAcc(a)
}

func UpdateAccount(account *models.Account, tx *models.Transaction, _type int, txInc uint64) {

	account.Balance = GetBalance(account.Address)
	account.Percent = utils.GetBalancePercent(account.Balance, common.VNT_TOTAL, common.VNT_DECIMAL)
	account.LastBlock = tx.BlockNumber

	if account.LastTx != tx.Hash {
		account.LastTx = tx.Hash
		account.TxCount += txInc
	}

	retAddrs := make([]string, 0)

	if _type == ACC_TYPE_CONTRACT {
		// if already exists as a normal account,
		// then now it turns out a new contract account
		if !account.IsContract {
			if ok, erc20 := IsToken(account.Address, tx); ok {
				account.IsToken = true
				account.TokenType = token.TOKEN_ERC20
				account.ContractName = erc20.TokenName
				account.TokenAmount = erc20.TotalSupply.String()
				account.TokenSymbol = erc20.Symbol
				account.TokenDecimals = erc20.Decimals.Uint64()
				account.TokenAcctCount = "0"
				account.TokenLogo = ""
			}

			account.IsContract = true
			account.ContractOwner = tx.From
		}
	} else if _type == ACC_TYPE_TOKEN {
		//tx.IsToken = true
		retAddrs = token.UpdateTokenBalance(account, tx)
	}

	updateAcc(account)
	// Save the accounts in token transfer
	for _, a := range retAddrs {
		if acct := GetAccount(a); acct != nil {
			acct.Balance = GetBalance(a)
			acct.Percent = utils.GetBalancePercent(acct.Balance, common.VNT_TOTAL, common.VNT_DECIMAL)
			acct.LastBlock = tx.BlockNumber
			if acct.LastTx != tx.Hash {
				acct.LastTx = tx.Hash
			}
			updateAcc(acct)
		} else {
			NewAccount(a, tx, ACC_TYPE_NORMAL, 0)
			beego.Info("Inserted accounts: ", a)
		}
	}
}

func PersistWitnesses(accts []string, blockNumber uint64) {
	beego.Info("Will persist witnesses accounts: ", accts)
	for _, a := range accts {
		if _, ok := AccountMap.Load(a); ok {
			continue
		}
		AccountMap.Store(a, 0)
		if acct := GetAccount(a); acct != nil {
			acct.Balance = GetBalance(a)
			acct.Percent = utils.GetBalancePercent(acct.Balance, common.VNT_TOTAL, common.VNT_DECIMAL)
			acct.LastBlock = blockNumber
			updateAcc(acct)
		} else {
			NewAccount(a, &models.Transaction{BlockNumber: blockNumber}, ACC_TYPE_NORMAL, 0)
			beego.Info("Inserted witness account: ", a)
		}
	}
}

func GetAccount(addr string) *models.Account {
	addr = strings.ToLower(addr)
	if _type, err := acctCache.Get(addr); err == nil && _type != nil {
		beego.Info("Address hit in cache:", addr)
		return _type.(*models.Account)
	} else {
		beego.Info("Address not hit in cache:", addr)
		a := &models.Account{}
		a, err := a.Get(addr)
		if err != nil {
			beego.Info("Address not hit in db:", addr)
			return nil
		}
		beego.Info("Address hit in db:", addr)
		acctCache.Set(addr, a)
		return a
	}
}

// insert into db and cache
func insertAcc(acct *models.Account) {
	//if err := acct.Insert(); err != nil {
	//	msg := fmt.Sprintf("Failed to insert account: %v, error: %s", acct, err.Error())
	//	beego.Error(msg)
	//	panic(msg)
	//}
	//acctCache.Set(acct.Address, acct)

	PostAccountTask(NewAccountTask(acct, ACTION_INSERT))
}

// update db and cache
func updateAcc(acct *models.Account) {
	//if err := acct.Update(); err != nil {
	//	msg := fmt.Sprintf("Failed to update account: %s, error: %s", addr, err.Error())
	//	beego.Error(msg)
	//	panic(err)
	//}
	//acctCache.Set(acct.Address, acct)
	PostAccountTask(NewAccountTask(acct, ACTION_UPDATE))
}

func InsertGenius() {
	genius, _, _ := GetBlock(0)
	var accounts = make([]*models.Account, 0)
	var txs = make([]*models.Transaction, 0)

	geniusAlloc := beego.AppConfig.String("node::genius_alloc")
	currDir, err := os.Getwd()
	if err != nil {
		msg := fmt.Sprintf("Failed to get current directory, err: %s", err.Error())
		beego.Error(msg)
		panic(err)
	}
	geniusAlloc = currDir + "/tools/racer/data/alloc.csv"
	f, err := os.Open(geniusAlloc)
	if err != nil {
		msg := fmt.Sprintf("Failed to read genius allocation snapshot file %s, err: %s", geniusAlloc, err.Error())
		beego.Error(msg)
		panic(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	rows, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	for i, row := range rows {
		fmt.Printf("Row %d, content: %s\n", i, row)
		tx, account := genTxAndAccount(row, genius)
		accounts = append(accounts, account)
		txs = append(txs, tx)
	}

	fmt.Println("Transactions: ", txs)
	fmt.Println("Accounts: ", accounts)

	for _, tx := range txs {
		fmt.Println("Inserting tx: ", tx.Hash)
		if err := tx.Insert(); err != nil {
			fmt.Println("Failed to insert tx: ", tx.Hash, " ,err", err)
			panic(err)
		}
	}

	for _, account := range accounts {
		account.Balance = GetBalance(account.Address)
		account.Percent = utils.GetBalancePercent(account.Balance, common.VNT_TOTAL, common.VNT_DECIMAL)

		a := &models.Account{}
		a, err := a.Get(account.Address)
		if err != nil {
			if err == orm.ErrNoRows {
				fmt.Println("Inserting account: ", account.Address)
				account = &models.Account{
					Address:        account.Address,
					Vname:          account.Address, //todo: get vname
					Balance:        account.Balance,
					TxCount:        1,
					FirstBlock:     0,
					LastBlock:      0,
					TokenAmount:    "0",
					TokenAcctCount: "0",
					InitTx:         "Genius_"+account.Address,
					LastTx:         "Genius_"+account.Address,
				}
				if err = account.Insert(); err != nil {
					fmt.Println("Failed to insert account: ", account.Address, " ,err", err)
					panic(err)
				}
			} else {
				fmt.Println("Failed to get account: ", account.Address, " ,err", err)
				panic(err)
			}
		} else {
			fmt.Println("Updating account: ", account.Address)
			a.Balance = account.Balance
			a.TxCount += 1
			if err := a.Update(); err != nil {
				fmt.Println("Failed to update account: ", account.Address, " ,err", err)
				panic(err)
			}
		}
	}

	fmt.Println("Updating genius block...")
	block := &models.Block{}
	if block, err = block.GetByNumber(0); err != nil {
		if err == orm.ErrNoRows {
			block = genius
		} else {
			fmt.Println("Failed to get block 0, err:", err)
			panic(err)
		}
	}

	block.TxCount = len(txs)

	err = block.Insert()
	if err != nil {
		fmt.Println("Failed to update genius block, err", err)
		panic(err)
	}

	fmt.Println("Done!")
}

func genTxAndAccount(snapshot []string, genius *models.Block) (*models.Transaction, *models.Account) {
	if len(snapshot) < 2 {
		msg := fmt.Sprintf("Invalide snapshot: %v", snapshot)
		beego.Error(msg)
		panic(msg)
	}

	address := snapshot[0]
	balance := snapshot[1]

	var account = &models.Account{
		Address: address,
	}

	var tx = &models.Transaction{
		Hash: "Genius_" + address,
		TimeStamp: genius.TimeStamp,
		From: "Genius",
		To: account,
		Value: balance,
		GasLimit: 0,
		GasPrice: "0",
		GasUsed: 0,
		Nonce: 0,
		Index: 0,
		Input: "",
		Status: 1,
		ContractAddr: "",
		IsToken: false,
		TokenFrom: "",
		TokenTo: "",
		TokenAmount: "",
		BlockNumber: 0,
	}

	return tx, account
}
