package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strings"
	"strconv"
	"fmt"
)

type Block struct {
	Number       uint64 `orm:"pk"`
	TimeStamp    uint64
	TxCount      int
	Hash         string   `orm:"unique"`
	ParentHash   string
	Producer     string  `orm:"index"`
	Size         string
	GasUsed      uint64
	GasLimit     uint64
	BlockReward  string
	ExtraData    string
	Witnesses    []*Node        `orm:"rel(m2m)"`
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
		b.Number, err = strconv.ParseUint(nOrh, 10, 64)
		if err != nil {
			msg := fmt.Sprintf("Wrong block number: %s", nOrh)
			beego.Error(msg)
		}
		err = o.Read(b, "Number")
	}

	return b, err
}

func (b *Block) Last() (*Block, error) {
	o := orm.NewOrm()

	qs := o.QueryTable(b).OrderBy("-Number").Limit(1)

	var blocks []*Block
	_, err := qs.All(&blocks)

	if err != nil {
		return nil, err
	}

	if len(blocks) == 0 {
		return nil, nil
	}

	return blocks[0], nil
}

func (b *Block) Count() (int64, error) {
	o := orm.NewOrm()
	cnt, err := o.QueryTable(b).Count()
	return cnt, err
}
