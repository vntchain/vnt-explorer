package models

import "github.com/astaxie/beego/orm"

type Test struct {
	Id	int
	Name	string
}

func init() {
	orm.RegisterModel(new(Test))
}