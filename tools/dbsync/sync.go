package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gookit/color"
	"github.com/vntchain/vnt-explorer/models"
	"strings"
)

func main() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	dbuser := beego.AppConfig.String("mysql::user")
	dbhost := beego.AppConfig.String("mysql::host")
	dbport := beego.AppConfig.String("mysql::port")
	dbname := beego.AppConfig.String("mysql::db")

	red := color.FgRed.Render
	green := color.FgGreen.Render

	fmt.Println(green("本操作将会对以下数据库进行重建："))
	fmt.Println(green("服务器："), dbhost)
	fmt.Println(green("端口："), dbport)
	fmt.Println(green("用户名："), dbuser)
	fmt.Println(green("数据库："), dbname)

	fmt.Println(red("危险！危险！危险！本操作会将当前数据库清空！"))
	fmt.Println(red("如果是线上数据库，建议您采用migration的方式进行数据库升级！\n"))

	var alert = []string{
		"你真的要这么做吗？", "一定要想清楚啊！",
		"不要冲动啊！", "你确定知道你在做什么吗？",
		"没有后悔药哦！", "好吧，看起来你很坚定！",
		"嗯，看好你哦！", "好吧，我们开始一个全新的世界吧！",
		"要不要再考虑一下？数据真的会清空哦。", "好的，let's go!",
	}

	var command string
	for _, console := range alert {
		fmt.Print(red(console), green("Y(继续)/N(放弃):"))
		fmt.Scan(&command)
		if strings.ToUpper(command) == "Y" {
			continue
		} else {
			fmt.Println("放弃！")
			break
		}
	}

	if strings.ToUpper(command) == "Y" {
		orm.RunSyncdb("default", true, true)
		alterTable()
	}
}

func registerModel() {
	beego.Info("Will register models.")
	orm.RegisterModel(new(models.Account))
	orm.RegisterModel(new(models.Block))
	orm.RegisterModel(new(models.Node))
	orm.RegisterModel(new(models.TokenBalance))
	orm.RegisterModel(new(models.Transaction))
	orm.RegisterModel(new(models.Hydrant))
	orm.RegisterModel(new(models.MarketInfo))
}

func alterTable() {
	needAlterMap := make(map[string][]string)
	needAlterMap["account"] = []string{"balance", "token_amount", "token_acct_count"}
	needAlterMap["block"] = []string{"number"}
	needAlterMap["node"] = []string{"votes"}
	needAlterMap["token_balance"] = []string{"balance"}
	for tableName, columns := range needAlterMap {
		for _, col := range columns {
			if err := alterColumn(tableName, col, "decimal(64,0)"); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func alterColumn(tableName, column, dataType string) error {
	o := orm.NewOrm()
	alterString := fmt.Sprintf("ALTER TABLE %s MODIFY %s %s", tableName, column, dataType)
	_, err := o.Raw(alterString).Exec()
	if err != nil {
		return fmt.Errorf("ALTER TABLE %s error: ", tableName, err)
	}
	return nil
}
