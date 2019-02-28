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

func (t *TokenBalance) List(account string, token string, offset, limit int) ([]*TokenBalance, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(t)
	cond := orm.NewCondition()
	if len(account) > 0 {
		cond = cond.And("account", account)
	} else if len(token) > 0 {
		cond = cond.And("token", token)
	}

	qs = qs.SetCond(cond)
	var tokens []*TokenBalance
	_, err := qs.Offset(offset).Limit(limit).All(&tokens)
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

func (t *TokenBalance) CountByToken() (int64, error) {
	o := orm.NewOrm()
	cnt, err := o.QueryTable(t).GroupBy("Token").Count()
	return cnt, err
}
