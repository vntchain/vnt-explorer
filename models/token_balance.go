package models

import (
	"github.com/astaxie/beego/orm"
)

type TokenBalance struct {
	Id      int
	Account *Account `orm:"rel(fk)"`
	Token   *Account `orm:"rel(fk)"`
	Balance string   `orm:"index"`
	Percent float32
}

func (t *TokenBalance) TableUnique() [][]string {
	return [][]string{
		{"account_id", "token_id"},
	}
}

func (t *TokenBalance) Insert() error {
	o := orm.NewOrm()
	_, err := o.Insert(t)
	return err
}

func (t *TokenBalance) Update() error {
	o := orm.NewOrm()
	_, err := o.Update(t)
	return err
}

func (t *TokenBalance) List(account, token, order string, offset, limit int, fields []string) ([]*TokenBalance, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(t)
	cond := orm.NewCondition()
	if len(account) > 0 {
		cond = cond.And("account__address", account)
	} else if len(token) > 0 {
		cond = cond.And("token__address", token)
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
	t.Account = &Account{Address: account}
	t.Token = &Account{Address: token}
	err := o.Read(t, "Account", "Token")
	if err == nil && t.Account != nil && t.Account.Address != "" {
		err = o.Read(t.Account)
	}
	if err == nil && t.Token != nil && t.Token.Address != "" {
		err = o.Read(t.Token)
	}
	return t, err
}

func (t *TokenBalance) GetById(id int) (*TokenBalance, error) {
	o := orm.NewOrm()
	t.Id = id
	err := o.Read(t)
	if err == nil && t.Account != nil && t.Account.Address != "" {
		err = o.Read(t.Account)
	}
	if err == nil && t.Token != nil && t.Token.Address != "" {
		err = o.Read(t.Token)
	}
	return t, err
}

func (t *TokenBalance) Count(account, token string) (int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(t)

	cond := orm.NewCondition()
	if account != "" {
		cond = cond.And("Account", account)
	}

	if token != "" {
		cond = cond.And("Token", token)
	}

	qs = qs.SetCond(cond)

	return qs.Count()
}
