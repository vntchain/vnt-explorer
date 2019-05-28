package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type V12_20190528_172518 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &V12_20190528_172518{}
	m.Created = "20190528_172518"

	migration.Register("V12_20190528_172518", m)
}

// Run the migrations
func (m *V12_20190528_172518) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("ALTER TABLE BLOCK " +
		"ADD COLUMN `reward` double NOT NULL DEFAULT '0', " +
		"ADD COLUMN `fee` double NOT NULL DEFAULT '0';")
	m.SQL("CREATE TABLE if not exists `market_info` (" +
		"`id` int(11) NOT NULL AUTO_INCREMENT, " +
		"`last_updated` bigint(20) NOT NULL DEFAULT '0'," +
		"`price_cny` double NOT NULL DEFAULT '0'," +
		"`price_usd` double NOT NULL DEFAULT '0'," +
		"`available_supply` double NOT NULL DEFAULT '0'," +
		"`volume24h` double NOT NULL DEFAULT '0'," +
		"`volume24h_cny` double NOT NULL DEFAULT '0'," +
		"`volume24h_usd` double NOT NULL DEFAULT '0'," +
		"`market_cap_cny` double NOT NULL DEFAULT '0'," +
		"`market_cap_usd` double NOT NULL DEFAULT '0'," +
		"`percent_change24h` double NOT NULL DEFAULT '0'," +
		"PRIMARY KEY (`id`)," +
		"KEY `market_info_last_updated` (`last_updated`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8")
}

// Reverse the migrations
func (m *V12_20190528_172518) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("ALTER TABLE `BLOCK` " +
		"DROP COLUMN `reward`, " +
		"DROP COLUMN `fee`;")
	m.SQL("DROP TABLE if exists `market_info`")
}
