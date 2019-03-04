package models

import (
	"github.com/astaxie/beego/orm"
)

type TokenBalance struct {
	Id      int
	Account string `orm:"index"`
	Token   string `orm:"index"`
	Balance string
}

func (t *TokenBalance) TableUnique() [][]string {
	return [][]string{
		{"Account", "Token"},
	}
}

func (t *TokenBalance) Insert() error {
	o := orm.NewOrm()
	_, err := o.Insert(t)
	return err
}

func (t *TokenBalance) List(account, token, order string, offset, limit int, fields []string) ([]*TokenBalance, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(t)
	cond := orm.NewCondition()
	if len(account) > 0 {
		cond = cond.And("Account", account)
	} else if len(token) > 0 {
		cond = cond.And("Token", token)
	}
	qs = qs.SetCond(cond)
	if order == "asc" {
		qs = qs.OrderBy("Balance")
	} else {
		qs = qs.OrderBy("-Balance")
	}

	var tokens []*TokenBalance
	_, err := qs.Offset(offset).Limit(limit).All(&tokens, fields...)
	return tokens, err
}

func (t *TokenBalance) GetByAddr(account string, token string) (*TokenBalance, error) {
	o := orm.NewOrm()
	t.Account = account
	t.Token = token
	err := o.Read(t, "Account", "Token")
	return t, err
}

func (t *TokenBalance) GetById(id int) (*TokenBalance, error) {
	o := orm.NewOrm()
	t.Id = id
	err := o.Read(t)
	return t, err
}
