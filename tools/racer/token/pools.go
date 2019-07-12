package token

import (
	"github.com/vntchain/vnt-explorer/tools/racer/pool"
	"runtime"
	"github.com/vntchain/vnt-explorer/models"
	"github.com/astaxie/beego/orm"
	"github.com/vntchain/vnt-explorer/common/utils"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"sync"
)

var TokenMap = sync.Map{}

var TokenPool = pool.New(runtime.NumCPU() * 3, 3000)

type TokenTask struct {
	pool.BasicTask
	Token		*models.Account
	Holder		string
}

func (this *TokenTask) DoWork(workRoutine int) {
	this.PreDoWork(workRoutine)
	var token = this.Token
	var addr = this.Holder
	if _, ok := TokenMap.Load(token.Address+"_"+addr); ok{
		return
	}

	TokenMap.Store(token.Address + "_" + addr,1)

	tokenBalance := &models.TokenBalance{}
	tokenBalance, err := tokenBalance.GetByAddr(addr, token.Address)
	if err != nil && err == orm.ErrNoRows {
		tokenBalance.Balance = GetAmount(token.Address, addr)
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

		beego.Debug("Success to insert token balance:", tokenBalance.Token, tokenBalance.Account, tokenBalance.Balance)
	} else if err != nil {
		msg := fmt.Sprintf("Failed to get token balance, token:%s, address:%s",
			tokenBalance.Token, tokenBalance.Account)
		beego.Error(msg)
		panic(msg)
	} else if tokenBalance.Id > 0 {
		tokenBalance.Balance = GetAmount(token.Address, addr)
		tokenBalance.Percent = utils.GetBalancePercent(tokenBalance.Balance, token.TokenAmount, int(token.TokenDecimals))
		if err := tokenBalance.Update(); err != nil {
			msg := fmt.Sprintf("Failed to update token balance, token:%s, address:%s, token:%s",
				tokenBalance.Token, tokenBalance.Account, tokenBalance.Account)
			beego.Error(msg)
			panic(msg)
		}
		beego.Debug("Success to update token balance:", tokenBalance.Token, tokenBalance.Account, tokenBalance.Balance)
	}
}

func NewTokenTask(Token *models.Account, Holder string) *TokenTask {
	return &TokenTask {
		pool.BasicTask{
			fmt.Sprintf("token-%s-%s", Token.Address, Holder),
			TokenPool,
		},
		Token,
		Holder,
	}
}

func PostTokenTask(task *TokenTask) {
	err := TokenPool.PostWork("token", task)
	if err != nil {
		beego.Error("帐户线程更新插入池满载！")
		panic("")
	}
}