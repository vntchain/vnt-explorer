package utils

import (
	"strconv"
	"fmt"
	"github.com/astaxie/beego"
)

func GetBalancePercent(balance string, totalSupply string, decimal int) float32 {
	if len(balance) < decimal {
		return 0
	}

	balance = balance[:len(balance)-decimal]
	b, err := strconv.ParseFloat(balance, 64)
	if err != nil {
		msg := fmt.Sprintf("failed to parse balance: %s", err.Error())
		beego.Error(msg)
		panic(msg)
	}

	totalSupply = totalSupply[:len(totalSupply) - decimal]

	t, err := strconv.ParseFloat(totalSupply, 64)
	return float32(b/t) * 100
}

