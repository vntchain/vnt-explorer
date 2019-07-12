package controllers

import (
	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/models"
	"strconv"
	"time"
)

type MarketController struct {
	BaseController
}

func (this *MarketController) History() {
	days := 14
	var err error
	beego.Info("Will get history...days: ", days)
	daysStr := this.GetString("days")
	if daysStr != "" {
		days, err = strconv.Atoi(daysStr)
		if err != nil {
			this.ReturnErrorMsg("Wrong format of parameter days: %s", err.Error(), "")
			return
		}
	}

	if days > 100 {
		days = 100
	}

	type Item struct {
		TimeStamp int64
		Year      int
		Month     int
		Day       int
		PriceCny  float64
		PriceUsd  float64
		Volume    float64
	}

	history := make([]Item, 0)

	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()

	end := time.Date(year, month, day, 23, 59, 59, 0, now.Location())
	start := end.AddDate(0, 0, -days)
	during := time.Hour * 24

	beego.Info("Will get history...start: ", start, "end:", end)

	market := &models.MarketInfo{}

	for ; end.Unix() >= start.Unix(); start = start.Add(during) {
		marketInfo, err := market.Get(start.Unix(), "lte")
		if err != nil {
			this.ReturnErrorMsg("Failed to get market history: %s", err.Error(), "")
			return
		}
		if marketInfo == nil {
			continue
		}

		if marketInfo.LastUpdated < start.Add(-during).Unix() || marketInfo.LastUpdated > start.Add(during).Unix() {
			continue
		}

		item := Item{
			marketInfo.LastUpdated,
			start.Year(),
			int(start.Month()),
			start.Day(),
			marketInfo.PriceCny,
			marketInfo.PriceUsd,
			marketInfo.Volume24h,
		}

		history = append(history, item)

	}

	this.ReturnData(history, nil)
}

func (this *MarketController) Market() {
	now := time.Now()
	market := &models.MarketInfo{}
	marketInfo, err := market.Get(now.Unix(), "lte")
	if err != nil {
		this.ReturnErrorMsg("Failed to get market info: %s", err.Error(), "")
		return
	}
	this.ReturnData(marketInfo, nil)
}
