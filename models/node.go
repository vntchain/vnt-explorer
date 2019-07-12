package models

import (
	"github.com/astaxie/beego/orm"
)

type Node struct {
	Address         string `orm:"pk"`
	Vname           string `orm:"unique"`
	Home            string
	Logo            string
	Ip              string
	IsSuper         int
	IsAlive         int
	Status          int `orm:"index"`
	Votes           string
	VotesFloat      float64 `orm:"-"`
	VotesPercent    float32
	Longitude       float64
	Latitude        float64
	Block           []*Block `orm:"reverse(many)"`
	City            string
	NodeUrl         string
}

func (n *Node) Insert() error {
	o := orm.NewOrm()
	_, err := o.InsertOrUpdate(n)
	return err
}

func (n *Node) List(order string, offset, limit int, fields []string) ([]*Node, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(n)

	if order == "asc" {
		qs = qs.OrderBy("Votes")
	} else {
		qs = qs.OrderBy("-Votes")
	}

	cond := orm.NewCondition()

	cond = cond.And("status", 1)

	qs = qs.SetCond(cond)

	var nodes []*Node
	_, err := qs.Offset(offset).Limit(limit).All(&nodes, fields...)
	return nodes, err
}

func (n *Node) Get(address string) (*Node, error) {
	o := orm.NewOrm()
	n.Address = address
	err := o.Read(n)
	return n, err
}

func (n *Node) Count(isSuper int) (int64, error) {
	o := orm.NewOrm()

	qs := o.QueryTable(n)

	cond := orm.NewCondition()

	cond = cond.And("status", 1)

	if isSuper != -1 {
		cond = cond.And("isSuper", isSuper)
	}

	qs = qs.SetCond(cond)
	return qs.Count()
}

func (n *Node) All() ([]*Node, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(n)
	cond := orm.NewCondition()
	cond = cond.And("status", 1)
	qs = qs.SetCond(cond)

	var nodes []*Node
	_, err := qs.All(&nodes)
	return nodes, err
}
