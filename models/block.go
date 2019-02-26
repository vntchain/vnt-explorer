package models

import "time"

type Block struct {
	Number       string `orm:"pk"`
	TimeStamp    time.Time
	TxCount      int
	Hash         string   `orm:"unique"`
	ParentHash   string   // FIXME 这里是否需要设置成外键
	Producer     *Account `orm:"rel(fk)"`
	Size         string
	GasUsed      uint64
	GasLimit     uint64
	BlockReard   string
	ExtraData    string
	Witnesses    []*Node        `orm:"rel(m2m)"`
	Transactions []*Transaction `orm:"reverse(many)"`
}
