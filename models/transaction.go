package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type Transaction struct {
	Hash        	string `orm:"pk"`
	TimeStamp   	uint64
	From        	string `orm:"index"`
	To          	string `orm:"index"`
	Value       	string
	GasLimit    	uint64
	GasPrice    	string
	GasUsed     	uint64
	Nonce       	uint64
	Index       	int
	Input       	string	`orm:"type(text)"`
	Status			int
	ContractAddr	string // when transaction is a contract creation
	IsToken     	bool
	TokenTo     	string `orm:"index"`
	TokenAmount 	string
 	BlockNumber 	uint64 `orm:"index"`
}

func makeCond(block string, account string, isToken int) *orm.Condition {
	cond := orm.NewCondition()
	if len(block) > 0 {
		cond = cond.And("blockNumber", block)
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
	return cond
}

func (t *Transaction) Insert() error {
	o := orm.NewOrm()
	_, err := o.Insert(t)
	return err
}

func (t *Transaction) Update() error {
	o := orm.NewOrm()
	_, err := o.Update(t)
	return err
}

func (t *Transaction) List(offset, limit int64, order, block string, account string, isToken int, fields ...string) ([]*Transaction, error) {
	o := orm.NewOrm()

	beego.Info("block:", block, "account:", account, "istoken:", isToken)

	qs := o.QueryTable(t).Offset(offset).Limit(limit);

	cond := makeCond(block, account, isToken)

	qs = qs.SetCond(cond)

	if order == "asc" {
		qs = qs.OrderBy("TimeStamp")
	} else {
		qs = qs.OrderBy("-TimeStamp")
	}

	var txs []*Transaction
	_, err := qs.All(&txs, fields...)
	return txs, err
}

func (t *Transaction) Get(hash string, fields ...string) (*Transaction, error) {
	o := orm.NewOrm()

	var err error

	t.Hash = hash
	err = o.Read(t, "Hash")

	return t, err
}

func (t *Transaction) Count(block string, account string, isToken int) (int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(t)
	cond := makeCond(block, account, isToken)

	qs = qs.SetCond(cond)

	return qs.Count()
}