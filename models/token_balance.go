package models

type TokenBalance struct {
	Id           int
	Address      *Account `orm:"rel(fk)"`
	TokenAddress *Account `orm:"rel(fk)"`
	Balance      string
}
