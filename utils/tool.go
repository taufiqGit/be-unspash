package utils

import (
	"strconv"
)

func ParseFloat64(str string) float64 {
	if str == "" {
		return 0
	}
	price, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return price
}

func ParseInt(str string) int {
	if str == "" {
		return 0
	}
	price, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return price
}
