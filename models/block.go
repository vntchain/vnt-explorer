package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type Block struct {
	Number       string `orm:"pk"`
	TimeStamp    int64
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

func (b *Block) Insert() (error, *Block) {
	o := orm.NewOrm()
	id, err := o.Insert(b)
	beego.Info("Created block with id: ", id)
	if err != nil {
		return err, nil
	}
	return nil, b
}

func (b *Block) List(offset, limit int) (error, []*Block) {
	o := orm.NewOrm()
	qs := o.QueryTable(b)

	var blocks []*Block
	_, err := qs.Offset(offset).Limit(limit).All(&blocks)
	return err, blocks
}