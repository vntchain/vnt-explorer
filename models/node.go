package models

type Node struct {
	Id      int
	Address string `orm:"index"`
	Vname   string `orm:"unique"`
	Home    string
	Logo    string
	Ip      string
	Status  int `orm:"index"`
	Votes   string
	Block   []*Block `orm:"reverse(many)"`
}
