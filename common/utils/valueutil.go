package utils

import (
	"fmt"
	"strconv"
)

func FromWei(value string) string {
	return FormatValue(value, 18)
}

func FormatValue(value string, decimals int) string {
	result := ""
	if decimals < 0 {
		return result
	}

	if len(value) <= decimals {
		result += "0"
		valueFormat := "%0" + strconv.Itoa(decimals) + "s"
		value = fmt.Sprintf(valueFormat, value)
	} else {
		result += value[:len(value)-decimals]
		value = value[len(value)-decimals:]
	}
	var notZeroPos int
	for notZeroPos = len(value) - 1; notZeroPos >= 0; notZeroPos-- {
		if value[notZeroPos] != '0' {
			break
		}
	}
	if notZeroPos >= 0 {
		result += "."
		result += value[:notZeroPos+1]
	}
	return result
}
