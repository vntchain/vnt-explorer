package models

import "time"

type Transaction struct {
	Hash        string `orm:"pk"`
	TimeStamp   time.Time `orm:"auto_now_add;type(datetime)"`
	From        string `orm:"index"`
	To          string `orm:"index"`
	Value       string
	GasLimit    uint64
	GasPrice    string
	GasUsed     uint64
	Nonce       uint64
	Index       int
	Input       string
	IsToken     bool
	TokenTo     string `orm:"index"`
	TokenAmount string
	BlockNumber *Block `orm:"rel(fk)";index`
}
