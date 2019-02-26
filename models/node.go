package models

type Node struct {
	Id      int
	Address *Account `orm:"rel(fk)"`
	Vname   string
	Home    string
	Logo    string
	Ip      string
	Status  int
	Votes   string
	Block   []*Block `orm:"reverse(many)"`
}
