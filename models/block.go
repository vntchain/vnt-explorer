package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strings"
	"strconv"
	"fmt"
)

type ErrorBlockNumber struct {
	format string
	number string
}

func (e ErrorBlockNumber) Error() string {
	return fmt.Sprintf(e.format, e.number)
}

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
	Tps			 float32
	Witnesses    []*Node        `orm:"rel(m2m)"`
}

func (b *Block) Insert() error {
	o := orm.NewOrm()
	_, err := o.InsertOrUpdate(b)
	return err
}

func (b *Block) List(offset, limit int64, order string, fields ...string) ([]*Block, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(b)

	if order == "asc" {
		qs = qs.OrderBy("Number")
	} else {
		qs = qs.OrderBy("-Number")
	}

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
			e := ErrorBlockNumber {"Wrong block number: %s", nOrh}
			beego.Error(e.Error())
			return nil, e
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

func (b *Block) TopTpsBlock() (*Block, error) {
	o := orm.NewOrm()

	qs := o.QueryTable(b).OrderBy("-tps").Limit(1)

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