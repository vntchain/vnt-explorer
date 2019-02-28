package models

import (
	"github.com/astaxie/beego/orm"
)

type Node struct {
	Address string `orm:"pk"`
	Vname   string `orm:"unique"`
	Home    string
	Logo    string
	Ip      string
	Status  int `orm:"index"`
	Votes   string
	Block   []*Block `orm:"reverse(many)"`
}

func (n *Node) Insert() error {
	o := orm.NewOrm()
	_, err := o.Insert(n)
	return err
}

func (n *Node) List(offset, limit int) ([]*Node, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(n)

	var nodes []*Node
	_, err := qs.Offset(offset).Limit(limit).All(&nodes)
	return nodes, err
}

func (n *Node) Get(address string) (*Node, error) {
	o := orm.NewOrm()
	n.Address = address
	err := o.Read(n)
	return n, err
}
