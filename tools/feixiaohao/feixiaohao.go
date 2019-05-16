package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/common"
	"github.com/vntchain/vnt-explorer/common/utils"
	"github.com/vntchain/vnt-explorer/models"
	"strconv"
	"time"
)

const (
	COINCODE      = "vntchain"
	FEIXIAOHAOURL = "http://dncapi.bqiapp.com/api/coin/"
)

var interval, intervalErr = beego.AppConfig.Int("market::interval")

type CoinInfoResp struct {
	Data   *coinInfo `json:"data"`
	Status string    `json:"status"`
	Code   string    `json:"code"`
	Msg    string    `json:"msg"`
}

type coinInfo struct {
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	Rank          int     `json:"rank"`
	PriceCNY      float64 `json:"price_cny"`
	ChangePercent float64 `json:"change_percent"`
	Supply        float64 `json:"supply"`
	MarketCap     float64 `json:"marketcap"`
	Volume24hCNY  float64 `json:"vol"`
	UpdateTime    int64   `json:"updatetime"`
}

type MarketInfoResp struct {
	Data     []*marketInfo `json:"data"`
	MaxPage  int           `json:"maxpage"`
	CurrPage int           `json:"currpage"`
	Code     int           `json:"code"`
	Msg      string        `json:"msg"`
}

type marketInfo struct {
	Vol        float64 `json:"vol"`
	Price      float64 `json:"price"`
	Accounting float64 `json:"accounting"`
}

func getCoinInfo(code string) (*coinInfo, error) {
	coinInfoUrl := FEIXIAOHAOURL + "coininfo/"
	params := []utils.Param{
		utils.Param{"code", code},
	}
	respBody, err := utils.CallApi(coinInfoUrl, params)
	if err != nil {
		return nil, err
	}

	var res CoinInfoResp
	if err = json.Unmarshal(respBody, &res); err != nil {
		return nil, fmt.Errorf("json unmarshal err: %s", err)
	}

	return res.Data, nil

}

func getVolInfo(code string) (float64, error) {
	marketUrl := FEIXIAOHAOURL + "market_ticker/"
	vol := 0.0
	for curPage := 1; ; curPage++ {
		params := []utils.Param{
			utils.Param{"code", code},
			utils.Param{"page", strconv.Itoa(curPage)},
			utils.Param{"pagesize", "100"},
		}
		respBody, err := utils.CallApi(marketUrl, params)
		if err != nil {
			return 0.0, err
		}

		var res MarketInfoResp
		if err = json.Unmarshal(respBody, &res); err != nil {
			return 0.0, fmt.Errorf("json unmarshal err: %s", err)
		}
		for _, market := range res.Data {
			vol += market.Vol
		}
		if res.CurrPage == res.MaxPage {
			break
		}
	}
	return vol, nil
}

func GetCoinAndInsertDB() {
	token, err1 := getCoinInfo(COINCODE)
	vol, err2 := getVolInfo(COINCODE)
	if err1 == nil && err2 == nil {
		coin := models.MarketInfo{
			LastUpdated:      token.UpdateTime,
			PriceCny:         token.PriceCNY,
			AvailableSupply:  token.Supply,
			Volume24h:        vol,
			Volume24hCny:     token.Volume24hCNY,
			MarketCapCny:     token.MarketCap,
			PercentChange24h: token.ChangePercent,
		}
		coin.Insert()
	} else {
		beego.Error("Get info from feixiaohao error: ", err1, err2)
	}
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
