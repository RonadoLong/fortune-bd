package model

import (
	"fmt"
	"fortune-bd/app/grid-strategy-svc/util/goex/binance"
	"fortune-bd/app/grid-strategy-svc/util/huobi"
	"strconv"
	"strings"

)

const (
	// 火币交易所
	ExchangeHuobi = "huobi"
	// okex交易所
	ExchangeOkex = "okex"
	// 币安交易所
	ExchangeBinance = "binance"

	// PrefixIDGlob 用户自定义id前缀，挂单买入
	PrefixIDGlob = "glob"
	// PrefixIDGlos 用户自定义id前缀，挂单卖出
	PrefixIDGlos = "glos"
	// PrefixIDMob 用户自定义id前缀，市价单买入
	PrefixIDMob = "mob"
	// PrefixIDMos 用户自定义id前缀，市价单卖出
	PrefixIDMos = "mos"

	// GridTypeNormal 普通网格
	GridTypeNormal = 0
	// GridTypeTrend 趋势网格
	GridTypeTrend = 1
	// GridTypeInfinite 无限网格
	GridTypeInfinite = 2
	// GridTypeReverse 反向网格
	GridTypeReverse = 3

	// 加解密过程填充的数据
	SecretSalt = "yoEa05cCxBw4QFY6"

	// DateTimeUTC 时间格式
	DateTimeUTC = "2006-01-02 15:04:05"
)

// FloatRound 截取小数位数，默认保留2位，四舍五入
func FloatRound(f float64, points ...int) float64 {
	size := 2
	if len(points) > 0 {
		size = points[0]
	}

	format := "%." + strconv.Itoa(size) + "f"
	if size < 1 {
		format = "%." + "f"
	}

	res, err := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	if err != nil {
		return 0
	}
	return res
}

// FloatRoundOff 截取小数位数，默认保留2位，非四舍五入
func FloatRoundOff(f float64, points ...int) float64 {
	val, _ := strconv.ParseFloat(Float64ToStr(f, points...), 64)
	return val
}

// Float64ToStr 截取小数位数，默认保留2位，不支持四舍五入，去掉后面的0
func Float64ToStr(f float64, points ...int) string {
	size := 2
	if len(points) > 0 {
		size = points[0]
	}

	result := strconv.FormatFloat(f, 'f', -1, 64)
	if result == "" {
		return ""
	}

	return ensurePointStr(result, size)
}

func ensurePointStr(str string, pointSize int) string {
	ss := strings.Split(str, ".")
	if len(ss) == 2 {
		if pointSize <= 0 {
			return ss[0]
		} else if len(ss[1]) <= pointSize {
			return str
		}
		return ss[0] + "." + ss[1][:pointSize]
	}

	return str
}

// RoundOffToStr 截取小数位数，默认保留2位，不支持四舍五入
func RoundOffToStr(f float64, points ...int) string {
	size := 2
	if len(points) > 0 {
		size = points[0]
	}

	result := fmt.Sprintf("%v", f)
	if result == "" {
		return ""
	}

	return ensurePointStr(result, size)
}

// GetExchangeFees 获取交易所手续费率
func GetExchangeFees(exchange string) float64 {
	fees := 0.001
	// 区分不同交易所
	switch exchange {
	case ExchangeHuobi:
		fees = huobi.FilledFees
	case ExchangeBinance:
		fees = binance.FilledFees
	}

	return fees
}
