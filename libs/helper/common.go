package helper

import (
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

func GetTimeNowOfUinx() int64 {
	return time.Now().UTC().Unix()
}

func GetTimeNow() time.Time {
	return time.Now().UTC()
}

func StringToFloat64(s string) float64 {
	d, _ := decimal.NewFromString(s)
	f, _ := d.Float64()
	return f
}

func Float64ToString(f float64) string {
	return decimal.NewFromFloat(f).String()
}

func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func Float32ToString(f float32) string {
	return decimal.NewFromFloat32(f).String()
}
