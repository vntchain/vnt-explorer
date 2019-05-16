package models

import "github.com/astaxie/beego/orm"

const TABLE = "market_info"

type MarketInfo struct {
	Id               int
	LastUpdated      int64 `orm:"index"`
	PriceCny         float64
	AvailableSupply  float64
	Volume24h        float64
	Volume24hCny     float64
	MarketCapCny     float64
	PercentChange24h float64
}

func (m *MarketInfo) Insert() error {
	o := orm.NewOrm()
	_, err := o.Insert(m)
	return err
}

// 获取大于等于t或者小于等于t的最近的数据
func (m *MarketInfo) Get(t int64, order string) (*MarketInfo, error) {
	var marketList []*MarketInfo
	o := orm.NewOrm()
	qs := o.QueryTable(TABLE)
	if order == "gte" {
		qs = qs.Filter("last_updated__gte", t).OrderBy("last_updated").Limit(1)
	} else {
		qs = qs.Filter("last_updated__lte", t).OrderBy("-last_updated").Limit(1)
	}

	_, err := qs.All(&marketList)
	if err != nil {
		return nil, err
	}
	if len(marketList) == 0 {
		return nil, nil
	}
	return marketList[0], nil
}
