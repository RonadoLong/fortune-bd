package global

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
)

func ObjectToString(val interface{}) string {
	marshal, _ := json.Marshal(val)
	return string(marshal)
}

func StringJoinString(val ...string) string {
	builder := strings.Builder{}
	for _, v := range val {
		builder.WriteString(v)
	}
	return builder.String()
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

func StringToFloat32(s string) float32 {
	d := StringToFloat64(s)
	return float32(d)
}
