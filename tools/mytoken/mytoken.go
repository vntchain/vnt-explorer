package main

import (
	"encoding/json"
	"strconv"
	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/common/utils"
	"github.com/vntchain/vnt-explorer/models"
	"github.com/vntchain/vnt-explorer/common"
	"time"
	"fmt"
)

const MYTOKENAPI = "https://api.mytokenapi.com/currency/currencydetail"

var interval, intervalErr = beego.AppConfig.Int("market::interval")

var params = []utils.Param{
	{"timestamp", "1557902202527"},
	{"code", "e320c466e1775f4dd9848a6018cf5f0f"},
	{"platform", "web_pc"},
	{"v", "1.0.0"},
	{"language", "zh_CN"},
	{"legal_currency", "cny"},
	{"com_id", "vnt_CNY"},
}

type TokenInfoResp struct {
	Data      *mytokenInfo `json:"data"`
	TimeStamp int64        `json:"timestamp"`
	Code      int          `json:"code"`
	Message   string       `json:"message"`
}

type mytokenInfo struct {
	Name                 string  `json:"name"`
	Rank                 int     `json:"rank"`
	MarketCapDisplayCny  float64 `json:"market_cap_display_cny"`
	PriceDisplayCNY      float64 `json:"price_display_cny"`
	PercentChangeDisplay string  `json:"percent_change_display"`
	PercentChangeUtc0    float64 `json:"percent_change_utc0"`
	AvaliableSupply      float64 `json:"available_supply"`
	Volume24hCNY         string  `json:"volume_24h"`
	Volume24h            float64 `json:"volume_24h_from"`
}

func GetCoinAndInsertDB() {
	respBody, err := utils.CallApi(MYTOKENAPI, params)
	var res TokenInfoResp
	if err = json.Unmarshal(respBody, &res); err != nil {
		msg := fmt.Sprintf("json unmarshal err: %s, respBody: %s\n", err,string(respBody))
		beego.Error(msg)
		return
	}
	if res.Data == nil {
		beego.Error("Get token info is nil")
		return
	}
	volCny, err := strconv.ParseFloat(res.Data.Volume24hCNY, 10)
	if err != nil {
		volCny = 0.0
	}

	coin := models.MarketInfo{
		LastUpdated:      res.TimeStamp,
		PriceCny:         res.Data.PriceDisplayCNY,
		AvailableSupply:  res.Data.AvaliableSupply,
		Volume24h:        res.Data.Volume24h,
		Volume24hCny:     volCny,
		MarketCapCny:     res.Data.MarketCapDisplayCny,
		PercentChange24h: res.Data.PercentChangeUtc0,
	}
	coin.Insert()

}

func main() {
	if intervalErr != nil {
		interval = common.DefaultMarketInterval
	}
	go GetCoinAndInsertDB()
	t := time.Tick(time.Second * time.Duration(interval))
	for range t {
		go GetCoinAndInsertDB()
	}
}
