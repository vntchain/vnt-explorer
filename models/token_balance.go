package models

type TokenBalance struct {
	Id		int
	Account	string `orm:"index"`
	Token	string `orm:"index"`
	Balance	string
}