package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/vntchain/vnt-explorer/models"
)

func main() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RunSyncdb("default", true, true)
	alterTable()
}

func registerModel() {
	beego.Info("Will register models.")
	orm.RegisterModel(new(models.Account))
	orm.RegisterModel(new(models.Block))
	orm.RegisterModel(new(models.Node))
	orm.RegisterModel(new(models.TokenBalance))
	orm.RegisterModel(new(models.Transaction))
	orm.RegisterModel(new(models.Hydrant))
}

func alterTable() {
	needAlterMap := make(map[string][]string)
	needAlterMap["account"] = []string{"balance", "token_amount", "token_acct_count"}
	needAlterMap["block"] = []string{"number"}
	needAlterMap["node"] = []string{"votes", "total_bounty", "extracted_bounty", "last_extract_time"}
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
