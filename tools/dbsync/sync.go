package main

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/models"
)

func main() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RunSyncdb("default", true, true)
}

func registerModel() {
	beego.Info("Will register models.")
	orm.RegisterModel(new(models.Account))
	orm.RegisterModel(new(models.Block))
	orm.RegisterModel(new(models.Node))
	orm.RegisterModel(new(models.TokenBalance))
	orm.RegisterModel(new(models.Transaction))
}