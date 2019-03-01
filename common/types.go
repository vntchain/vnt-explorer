package common

import (
	"github.com/vntchain/vnt-explorer/common/utils"
	"fmt"
	"github.com/astaxie/beego"
)

type Hex string

func (hex Hex) ToUint64() uint64 {
	beego.Info("Will convert hex", hex)
	r,e := utils.DecodeUint64(string(hex))

	if e != nil {
		msg := fmt.Sprintf("Failed to decode hex to uint64: %s", e.Error())
		beego.Error(msg)
		panic(msg)
	}

	return r
}

func (hex Hex) ToString() string {
	b,e := utils.DecodeBig(string(hex))
	if e != nil {
		msg := fmt.Sprintf("Failed to decode hex to big: %s", e.Error())
		beego.Error(msg)
		panic(msg)
	}

	return b.String()
}

func (hex Hex) ToInt() int {
	return int(hex.ToUint64())
}

func (hex Hex) ToInt64() int64 {
	b,e := utils.DecodeBig(string(hex))
	if e != nil {
		msg := fmt.Sprintf("Failed to decode hex to big: %s", e.Error())
		beego.Error(msg)
		panic(msg)
	}
	return b.Int64()
}