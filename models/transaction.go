package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Transaction struct {
	Hash         string   `orm:"pk"`
	TimeStamp    uint64   `orm:"index"`
	From         string   `orm:"index"`
	To           *Account `orm:"rel(fk);null;index"`
	Value        string
	GasLimit     uint64
	GasPrice     string
	GasUsed      uint64
	Nonce        uint64
	Index        int
	Input        string `orm:"type(text)"`
	Status       int
	ContractAddr string // when transaction is a contract creation
	IsToken      bool	`orm:"index"`
	TokenFrom    string `orm:"index"`
	TokenTo      string `orm:"index"`
	TokenAmount  string
	BlockNumber  uint64 `orm:"index"`
}

func makeCond(block string, account string, isToken int, from string, start, end int64) *orm.Condition {
	cond := orm.NewCondition()
	if len(block) > 0 {
		cond = cond.And("blockNumber", block)
	}

	if len(account) > 0 && isToken != 1 {
		cond2 := orm.NewCondition()
		cond = cond.AndCond(cond2.Or("from", account).Or("to", account))
	}

	if isToken == 0 {
		cond = cond.And("is_token", false)
	} else if isToken == 1 {
		cond = cond.And("is_token", true)

		if len(account) > 0 {
			if from == "" {
				from = "account"
			}
			// only returns txs that token_to or token_from indicates the address
			if from == "account" {
				cond2 := orm.NewCondition()
				cond = cond.AndCond(cond2.Or("token_to", account).
					Or("token_from", account))
			} else {
				cond = cond.And("to", account)
			}
		}
	}

	if start >= 0 {
		cond = cond.And("time_stamp__gt", start)
	}

	if end >= 0 {
		cond = cond.And("time_stamp__lte", end)
	}
	return cond
}

func (t *Transaction) Insert() error {
	o := orm.NewOrm()
	//_, err := o.Insert(t)
	_, err := o.InsertOrUpdate(t)
	return err
}

func (t *Transaction) Update() error {
	o := orm.NewOrm()
	_, err := o.Update(t)
	return err
}

func (t *Transaction) List(offset, limit int64, order, block string, account string, isToken int, from string, start, end int64, fields ...string) ([]*Transaction, error) {
	o := orm.NewOrm()

	beego.Info("block:", block, "account:", account, "istoken:", isToken)

	qs := o.QueryTable(t).Offset(offset).Limit(limit)

	cond := makeCond(block, account, isToken, from, start, end)

	qs = qs.SetCond(cond)

	if order == "asc" {
		qs = qs.OrderBy("TimeStamp")
	} else {
		qs = qs.OrderBy("-TimeStamp")
	}

	var txs []*Transaction
	_, err := qs.All(&txs, fields...)
	for _, tx := range txs {
		if tx.To != nil && tx.To.Address != "" {
			o.Read(tx.To)
		}
	}
	return txs, err
}

func (t *Transaction) Get(hash string, fields ...string) (*Transaction, error) {
	o := orm.NewOrm()

	var err error

	t.Hash = hash
	err = o.Read(t, "Hash")
	if err == nil && t.To != nil && t.To.Address != "" {
		o.Read(t.To)
	}
	return t, err
}

func (t *Transaction) Count(block string, account string, isToken int, from string, start, end int64) (int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(t)
	cond := makeCond(block, account, isToken, from, start, end)

	qs = qs.SetCond(cond)

	return qs.Count()
}
