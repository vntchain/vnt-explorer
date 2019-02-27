package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type Transaction struct {
	Hash        string `orm:"pk"`
	TimeStamp   int64
	From        string `orm:"index"`
	To          string `orm:"index"`
	Value       string
	GasLimit    uint64
	GasPrice    string
	GasUsed     uint64
	Nonce       uint64
	Index       int
	Input       string
	IsToken     bool
	TokenTo     string `orm:"index"`
	TokenAmount string
	BlockNumber *Block `orm:"rel(fk)";index`
}

func (t *Transaction) Insert() error {
	o := orm.NewOrm()
	_, err := o.Insert(t)
	return err
}

func (t *Transaction) List(offset, limit int64, block string, account string, isToken int, fields ...string) ([]*Transaction, error) {
	o := orm.NewOrm()

	beego.Info("block:", block, "account:", account, "istoken:", isToken)

	qs := o.QueryTable(t).Offset(offset).Limit(limit);

	cond := orm.NewCondition()
	if len(block) > 0 {
		cond = cond.And("block_number_id", block)
	}

	if len(account) > 0 {
		cond2 := orm.NewCondition()
		cond = cond.AndCond(cond2.Or("from", account).Or("to", account).Or("token_to", account))
	}

	if isToken == 0 {
		cond = cond.And("is_token", false)
	} else if isToken == 1 {
		cond = cond.And("is_token", true)
	}

	qs = qs.SetCond(cond)

	var txs []*Transaction
	_, err := qs.All(&txs, fields...)
	return txs, err
}