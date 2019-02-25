package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"fmt"
	"github.com/vntchain/go-vnt/log"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	dbuser := beego.AppConfig.String("mysql::user")
	dbpass := beego.AppConfig.String("mysql::pass")
	//dbhost := beego.AppConfig.String("mysql::host")
	//dbport := beego.AppConfig.String("mysql::port")
	dbname := beego.AppConfig.String("mysql::db")

	dbUrl := fmt.Sprintf("%s:%s@/%s?charset=utf8", dbuser, dbpass, dbname)
	beego.Info("Will connect to mysql url", dbUrl)
	err := orm.RegisterDataBase("default", "mysql", dbUrl)
	if err != nil {

		log.Error("failed to register database", err)
		panic(err.Error())
	}

	orm.RunSyncdb("default", false, true)
}