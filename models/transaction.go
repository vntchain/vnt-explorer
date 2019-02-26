package models

import "time"

type Transaction struct {
	Id          int
	Hash        string
	BlockNumber string
	TimeStamp   time.Time
	From        string
	To          string
	Value       string
	GasLimit    uint64
	GasPrice    string
	GasUsed     uint64
	Nonce       uint64
	Index       int
	Input       string
	IsToken     bool
	TokenTo     string
	TokenAmount string
}
