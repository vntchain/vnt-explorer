package models

import "time"

type Transaction struct {
	Hash        string `orm:"pk"`
	TimeStamp   time.Time
	From        *Account `orm:"rel(fk)"`
	To          *Account `orm:"rel(fk)"`
	Value       string
	GasLimit    uint64
	GasPrice    string
	GasUsed     uint64
	Nonce       uint64
	Index       int
	Input       string
	IsToken     bool
	TokenTo     *Account `orm:"rel(fk)"`
	TokenAmount string
	Block       *Block `orm:"rel(fk)"`
}
