package data

import (
	"fmt"
	"github.com/vntchain/vnt-explorer/models"
	"runtime"
	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/tools/racer/pool"
	"strings"
)

var BlockPool = pool.New(runtime.NumCPU() * 3, 50)
var BlockInsertPool = pool.New(runtime.NumCPU() * 3, 50)
var TxPool = pool.New(runtime.NumCPU() * 3, 3000)
var AccountExtractPool = pool.New(runtime.NumCPU() * 3, 3000)
var AccountPool = pool.New(runtime.NumCPU() * 3, 10000)
var WitnessesPool = pool.New(runtime.NumCPU() * 3, 20)
var NodePool = pool.New(runtime.NumCPU() * 3, 10)

type BlockTask struct {
	pool.BasicTask
	BlockNumber	int64
}

func (this *BlockTask) DoWork(workRoutine int) {
	this.PreDoWork(workRoutine)
	PersistBlock(this.BlockNumber)
}

func NewBlockTask(BlockNumber int64) *BlockTask {
	return &BlockTask {
		BasicTask: pool.BasicTask{
			Name: fmt.Sprintf("Block-%d", BlockNumber),
			Pool: BlockPool,
		},
		BlockNumber: BlockNumber,
	}
}

type BlockInsertTask struct {
	pool.BasicTask
	Block	*models.Block
}

func (this *BlockInsertTask) DoWork(workRoutine int) {
	this.PreDoWork(workRoutine)
	beego.Info("Will insert block:", this.Block.Number)
	err := this.Block.Insert()
	if err != nil {
		msg := fmt.Sprintf("Failed to insert or update block: %v, error: %s,", this.Block, err.Error())
		panic(msg)
	}
}

func NewBlockInsertTask(Block *models.Block) *BlockInsertTask {
	return &BlockInsertTask {
		BasicTask: pool.BasicTask{
			Name: fmt.Sprintf("Block-Insert-%d", Block.Number),
			Pool: BlockPool,
		},
		Block: Block,
	}
}

type TxTask struct {
	pool.BasicTask
	Tx	*models.Transaction
}

func (this *TxTask) DoWork(workRoutine int) {
	this.PreDoWork(workRoutine)

	err := this.Tx.Insert()
	if err != nil {
		msg := fmt.Sprintf("Failed to insert transaction: %s", err.Error())
		panic(msg)
	}
}

func NewTxTask(Tx *models.Transaction) *TxTask {
	return &TxTask {
		BasicTask: pool.BasicTask{
			Name: fmt.Sprintf("Tx-%s", Tx.Hash),
			Pool: TxPool,
		},
		Tx: Tx,
	}
}

type ExtractAccountTask struct {
	pool.BasicTask
	Tx		*models.Transaction
}

func (this *ExtractAccountTask) DoWork(workRoutine int) {
	this.PreDoWork(workRoutine)

	ExtractAcct(this.Tx)
}

func NewExtractAccountTask(Tx *models.Transaction) *ExtractAccountTask {
	return &ExtractAccountTask {
		pool.BasicTask{
			fmt.Sprintf("ext-account-%s", Tx.Hash),
			AccountExtractPool,
		},
		Tx,
	}
}

const (
	ACTION_INSERT = 1
	ACTION_UPDATE = 2
)

type AccountTask struct {
	pool.BasicTask
	Account		*models.Account
	Action		int
}

func (this *AccountTask) DoWork(workRoutine int) {
	this.PreDoWork(workRoutine)

	switch this.Action {
	case ACTION_INSERT:
		if err := this.Account.Insert(); err != nil {
			msg := fmt.Sprintf("Failed to insert account: %v, error: %s", this.Account.Address, err.Error())
			beego.Error(msg)
			panic(msg)
		}
		acctCache.Set(this.Account.Address, this.Account)
		break
	case ACTION_UPDATE:
		if err := this.Account.Update(); err != nil {
			msg := fmt.Sprintf("Failed to update account: %s, error: %s", this.Account.Address, err.Error())
			beego.Error(msg)
			panic(err)
		}
		acctCache.Set(this.Account.Address, this.Account)
		break
	default:

	}
}

func NewAccountTask(Account *models.Account, Action int) *AccountTask {
	return &AccountTask {
		pool.BasicTask{
			fmt.Sprintf("account-%s", Account.Address),
			AccountPool,
		},
		Account,
		Action,
	}
}

type WitnessesTask struct {
	pool.BasicTask
	Witnesses		[]string
	BlockNumber		uint64
}

func (this *WitnessesTask) DoWork(workRoutine int) {
	this.PreDoWork(workRoutine)

	PersistWitnesses(this.Witnesses, this.BlockNumber)
}

func NewWitnessesTask(Witnesses []string, BlockNumber uint64) *WitnessesTask {
	return &WitnessesTask {
		pool.BasicTask{
			"witnesses",
			AccountPool,
		},
		Witnesses,
		BlockNumber,
	}
}


type NodesTask struct {
	pool.BasicTask
}

func (this *NodesTask) DoWork(workRoutine int) {
	this.PreDoWork(workRoutine)

	witnesses := GetWitnesses(-1)
	witMap := make(map[string]int)
	for _, w := range witnesses {
		addr := fmt.Sprintf("%v", w)
		witMap[strings.ToLower(addr)] = 1
	}

	nodes := GetNodes()
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

func NewNodesTask() *NodesTask {
	return &NodesTask {
		pool.BasicTask{
			"nodes",
			AccountPool,
		},
	}
}

func PostBlockTask(task *BlockTask) {
	err := BlockPool.PostWork("block", task)
	if err != nil {
		beego.Error("区块线程池满载！")
		panic("")
	}
}

func PostInsertBlockTask(task *BlockInsertTask) {
	err := BlockInsertPool.PostWork("block", task)
	if err != nil {
		beego.Error("区块插入线程池满载！")
		panic("")
	}
}

func PostTxTask(task *TxTask) {
	err := TxPool.PostWork("tx", task)
	if err != nil {
		beego.Error("交易线程池满载！")
		panic("")
	}
}

func PostExtractAccountTask(task *ExtractAccountTask) {
	err := AccountExtractPool.PostWork("ext-account", task)
	if err != nil {
		beego.Error("帐户线程池满载！")
		panic("")
	}
}

func PostAccountTask(task *AccountTask) {
	err := AccountPool.PostWork("account", task)
	if err != nil {
		beego.Error("帐户线程更新插入池满载！")
		panic("")
	}
}

func PostWitnessesTask(task *WitnessesTask) {
	err := AccountPool.PostWork("witnesses", task)
	if err != nil {
		beego.Error("Witnesses池满载！")
		panic("")
	}
}

func PostNodesTask(task *NodesTask) {
	err := NodePool.PostWork("nodes", task)
	if err != nil {
		beego.Error("Nodes池满载！")
		panic("")
	}
}