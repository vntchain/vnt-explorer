package models

import "github.com/astaxie/beego/orm"

type Hydrant struct {
	Address   string `orm:"pk"`
	TimeStamp int64
}

func (h *Hydrant) InsertOrUpdate() error {
	o := orm.NewOrm()
	_, err := o.InsertOrUpdate(h)
	return err
}

func (h *Hydrant) Get(address string) (*Hydrant, error) {
	o := orm.NewOrm()
	h.Address = address
	err := o.Read(h)
	return h, err
}
