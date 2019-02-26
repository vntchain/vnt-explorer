package models

import "time"

type Block struct {
	Number       string `orm:"pk"`
	TimeStamp    time.Time `orm:"auto_now_add;type(datetime)"`
	TxCount      int
	Hash         string   `orm:"unique"`
	ParentHash   string   // FIXME 这里是否需要设置成外键
	Producer     string  `orm:"index"`
	Size         string
	GasUsed      uint64
	GasLimit     uint64
	BlockReward  string
	ExtraData    string
	Witnesses    []*Node        `orm:"rel(m2m)"`
	Transactions []*Transaction `orm:"reverse(many)"`
}
