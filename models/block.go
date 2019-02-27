package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"strings"
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

func (b *Block) Insert() error {
	o := orm.NewOrm()
	_, err := o.Insert(b)
	return err
}

func (b *Block) List(offset, limit int64, fields ...string) ([]*Block, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(b)

	var blocks []*Block
	_, err := qs.Offset(offset).Limit(limit).All(&blocks, fields...)
	return blocks, err
}

func (b *Block) Get(nOrh string, fields ...string) (*Block, error) {
	o := orm.NewOrm()

	var err error
	if strings.HasPrefix(nOrh, "0x") {
		beego.Info("Will read block by hash: ", nOrh)
		b.Hash = nOrh
		err = o.Read(b, "Hash")
	} else {
		beego.Info("Will read block by number: ", nOrh)
		b.Number = nOrh
		err = o.Read(b, "Number")
	}

	return b, err
}

func (b *Block) Count() (int64, error) {
	o := orm.NewOrm()
	cnt, err := o.QueryTable(b).Count()
	return cnt, err
}