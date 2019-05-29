package models

import "github.com/astaxie/beego/orm"

type Subscription struct {
	Email     string `orm:"pk"`
	TimeStamp uint64
}

func (s *Subscription) Insert() error {
	o := orm.NewOrm()
	_, err := o.Insert(s)
	return err
}

func (s *Subscription) Get(email string) (*Subscription, error) {
	o := orm.NewOrm()
	s.Email = email
	err := o.Read(s)
	return s, err
}
