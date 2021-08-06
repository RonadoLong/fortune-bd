package utils

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
)

func Keep2Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func Keep8Decimal(value float64) float64 {
	value, _ = decimal.NewFromFloat(value).Round(8).Float64()
	return value
}
