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
	m.SQL("ALTER TABLE `block` " +
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

	m.SQL("ALTER TABLE `node` " +
		"ADD COLUMN `city` varchar(255) NOT NULL DEFAULT '', " +
		"ADD COLUMN `node_url` varchar(255) NOT NULL DEFAULT ''," +
		"DROP COLUMN `total_bounty`, " +
		"DROP COLUMN `extracted_bounty`, " +
		"DROP COLUMN `last_extract_time`;")

	m.SQL("CREATE TABLE IF NOT EXISTS `subscription` (" +
		"`email` varchar(255) NOT NULL PRIMARY KEY," +
		"`time_stamp` bigint(20) unsigned NOT NULL DEFAULT 0" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8")
}

// Reverse the migrations
func (m *V12_20190528_172518) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("ALTER TABLE `block` " +
		"DROP COLUMN `reward`, " +
		"DROP COLUMN `fee`;")
	m.SQL("DROP TABLE if exists `market_info`")
	m.SQL("ALTER TABLE `node` " +
		"DROP COLUMN `city`, " +
		"DROP COLUMN `node_url`," +
		"ADD COLUMN `total_bounty` decimal(64,0), " +
		"ADD COLUMN `extracted_bounty` decimal(64,0), " +
		"ADD COLUMN `last_extract_time` decimal(64,0);")
	m.SQL("DROP TABLE if exists `subscription`")

}
