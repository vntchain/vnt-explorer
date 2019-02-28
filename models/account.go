package models

import (
	"github.com/astaxie/beego/orm"
)

type Account struct {
	Address        string `orm:"pk"`
	Vname          string `orm:"unique"`
	Balance        string
	TxCount        uint64
	IsContract     bool
	ContractName   string
	ContractOwner  string `orm:"index"`
	Code           string
	Abi            string
	Home           string
	InitTx         string
	IsToken        bool
	TokenType      int
	TokenSymbol    string
	TokenLogo      string
	TokenAmount    string
	TokenAcctCount string
	FirstBlock     string
	LastBlock      string
}

func (a *Account) Insert() error {
	o := orm.NewOrm()
	_, err := o.Insert(a)
	return err
}

func (a *Account) List(offset, limit int) ([]*Account, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(a)

	var accounts []*Account
	_, err := qs.Offset(offset).Limit(limit).All(&accounts)
	return accounts, err
}

func (a *Account) Get(address string) (*Account, error) {
	o := orm.NewOrm()
	a.Address = address
	err := o.Read(a)
	return a, err
}

func (a *Account) Count() (int64, error) {
	o := orm.NewOrm()
	cnt, err := o.QueryTable(a).Count()
	return cnt, err
}
