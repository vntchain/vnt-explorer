package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	registerModel()

	dbuser := beego.AppConfig.String("mysql::user")
	dbpass := beego.AppConfig.String("mysql::pass")
	dbhost := beego.AppConfig.String("mysql::host")
	dbport := beego.AppConfig.String("mysql::port")
	dbname := beego.AppConfig.String("mysql::db")

	//dbUrl := fmt.Sprintf("%s:%s@/%s?charset=utf8", dbuser, dbpass, dbname)
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", dbuser, dbpass, dbhost, dbport, dbname)
	beego.Info("Will connect to mysql url", dbUrl)
	err := orm.RegisterDataBase("default", "mysql", dbUrl)
	if err != nil {
		beego.Error("failed to register database", err)
		panic(err.Error())
	}
}

func registerModel() {
	beego.Info("Will register models.")
	orm.RegisterModel(new(Account))
	orm.RegisterModel(new(Block))
	orm.RegisterModel(new(Node))
	orm.RegisterModel(new(TokenBalance))
	orm.RegisterModel(new(Transaction))
	orm.RegisterModel(new(Hydrant))
}
