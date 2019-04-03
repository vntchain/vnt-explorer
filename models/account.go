package models

import (
	"github.com/astaxie/beego/orm"
)

type Account struct {
	Address        string `orm:"pk"`
	Vname          string `orm:"unique"`
	Balance        string `orm:"index"`
	TxCount        uint64
	IsContract     bool   `orm:"index"`
	ContractName   string
	ContractOwner  string `orm:"index"`
	Code           string `orm:"type(text)"`
	Abi            string `orm:"type(text)"`
	Home           string
	InitTx         string
	LastTx		   string
	IsToken        bool	  `orm:"index"`
	TokenType      int
	TokenSymbol    string
	TokenLogo      string
	TokenAmount    string
	TokenDecimals  uint64
	TokenAcctCount string
	FirstBlock     uint64
	LastBlock      uint64
	Percent 	   float32
}

func (a *Account) Insert() error {
	o := orm.NewOrm()
	_, err := o.InsertOrUpdate(a)
	return err
}

func (a *Account) Update() error {
	o := orm.NewOrm()
	_, err := o.Update(a)
	return err
}

func (a *Account) List(isContract, isToken int, order string, offset, limit int, fields []string) ([]*Account, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(a)
	cond := a.makeCond(isContract, isToken)
	qs = qs.SetCond(cond)

	if order == "asc" {
		qs = qs.OrderBy("Balance")
	} else {
		qs = qs.OrderBy("-Balance")
	}
	var accounts []*Account
	_, err := qs.Offset(offset).Limit(limit).All(&accounts, fields...)
	return accounts, err
}

func (a *Account) Get(address string) (*Account, error) {
	o := orm.NewOrm()
	a.Address = address
	err := o.Read(a)
	return a, err
}

func (a *Account) Count(isContract, isToken int) (int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(a)
	cond := a.makeCond(isContract, isToken)
	qs = qs.SetCond(cond)
	cnt, err := qs.Count()
	return cnt, err
}

func (a *Account) makeCond(isContract, isToken int) *orm.Condition {
	cond := orm.NewCondition()
	if isToken >= 0 {
		cond = cond.And("IsToken", isToken == 1)
	} else if isContract >= 0 {
		cond = cond.And("IsContract", isContract == 1)
	}

	return cond
}