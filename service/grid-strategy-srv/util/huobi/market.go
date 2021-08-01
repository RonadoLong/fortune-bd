package huobi

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
	"wq-fotune-backend/service/grid-strategy-srv/util/huobi/pkg/client"
	"wq-fotune-backend/service/grid-strategy-srv/util/huobi/pkg/getrequest"

	"github.com/zhufuyi/pkg/gohttp"
)

const (
	MIN1  = "1min"
	MIN5  = "5min"
	MIN15 = "15min"
	MIN30 = "30min"
	MIN60 = "60min"
	HOUR4 = "4hour"
	DAY1  = "1day"
	MON1  = "1mon"
	WEEK1 = "1week"
	YEAR1 = "1year"
)

// SymbolInfo 品种信息
type SymbolInfo struct {
	Symbol           string  // 品种
	MinOrderValue    float64 // 买卖最小交易额(单位是品种后面的货币)
	MinOrderAmt      float64 // 买卖最小交易所数量
	PricePrecision   int     // 价格精度
	AccountPrecision int     // 数量精度
	LeverageRatio    float64 // 最大杠杆比例
}

// GetSymbols 获取品种信息
func GetSelectSymbols(symbols []string) ([]*SymbolInfo, error) {
	symbolInfors, err := GetSymbols()
	if err != nil {
		return symbolInfors, err
	}

	sis := []*SymbolInfo{}
	for _, v := range symbolInfors {
		for _, symbol := range symbols {
			if v.Symbol == symbol {
				sis = append(sis, v)
			}
		}
	}

	return sis, nil
}

// GetAnchorCurrencySymbols 获取某个币本位下的所有品种信息
func GetAnchorCurrencySymbols(anchorCurrency string) ([]*SymbolInfo, error) {
	limitSymbols := []*SymbolInfo{}
	if anchorCurrency == "" {
		return nil, errors.New("anchorCurrency is empty")
	}

	symbolInfors, err := GetSymbols()
	if err != nil {
		return symbolInfors, err
	}

	acLen := len(anchorCurrency)
	for _, v := range symbolInfors {
		symbol := strings.ToLower(v.Symbol)
		if len(symbol) > acLen {
			if symbol[len(symbol)-acLen:] == anchorCurrency {
				limitSymbols = append(limitSymbols, v)
			}
		}
	}

	return limitSymbols, nil
}

// GetSymbols 获取所有品种信息
func GetSymbols() ([]*SymbolInfo, error) {
	symbols := []*SymbolInfo{}
	cli := new(client.CommonClient).Init(host)

	resp, err := cli.GetSymbols()
	if err != nil {
		return symbols, err
	}

	for _, v := range resp {
		if v.State == "online" {
			minOrderValue, _ := v.MinOrderValue.Float64()
			minOrderAmt, _ := v.MinOrderAmt.Float64()
			leverageRatio, _ := v.LeverageRatio.Float64()

			symbols = append(symbols, &SymbolInfo{
				Symbol:           v.Symbol,
				MinOrderAmt:      minOrderAmt,
				MinOrderValue:    minOrderValue,
				PricePrecision:   v.PricePrecision,
				AccountPrecision: v.AmountPrecision,
				LeverageRatio:    leverageRatio,
			})
		}
	}

	return symbols, nil
}

// -------------------------------------------------------------------------------------------------

// LatestTradeInfo 最新交易数据
type LatestTradeInfo struct {
	Status string `json:"status"`
	Tick   *struct {
		Data []struct {
			Amount    float64 `json:"amount"`
			Price     float64 `json:"price"`
			Direction string  `json:"direction"`
		}
	} `json:"tick"`
}

// GetLatestPrice 获取品种的最新价格
func GetLatestPrice(symbol string) (float64, error) {
	var price float64

	resp := &LatestTradeInfo{}
	url := fmt.Sprintf("https://%s/market/trade", host)
	params := map[string]interface{}{"symbol": symbol}

	err := gohttp.GetJSON(resp, url, params)
	if err != nil {
		return price, err
	}

	if resp.Status == "ok" && resp.Tick != nil {
		if len(resp.Tick.Data) > 0 {
			price = resp.Tick.Data[0].Price
		}
	}

	if price <= 0.0 {
		return price, fmt.Errorf("price is zero, resp=%+v", resp)
	}

	return price, nil
}

// Kline K线数据
type Kline struct {
	Time   string  `json:"time"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

// GetKlines 获取品种K线数据，period值为1min, 5min, 15min, 30min, 60min, 4hour, 1day, 1mon, 1week, 1year
func GetKlines(symbol string, period string, size int) ([]*Kline, error) {
	klines := []*Kline{}

	cli := new(client.MarketClient).Init(host)
	optionalRequest := getrequest.GetCandlestickOptionalRequest{Period: period, Size: size}
	resp, err := cli.GetCandlestick(symbol, optionalRequest)
	if err != nil {
		return klines, err
	}

	var atTime string
	var open, high, low, closePrice, volume float64
	for _, v := range resp {
		t := time.Unix(v.Id, 0)
		atTime = t.UTC().Format("2006-01-02T15:04:05.000Z")
		open, _ = v.Open.Float64()
		high, _ = v.High.Float64()
		low, _ = v.Low.Float64()
		closePrice, _ = v.Close.Float64()
		volume, _ = v.Vol.Float64()
		klines = append(klines, &Kline{atTime, open, high, low, closePrice, volume})
	}

	return klines, nil
}

// Get3CyclePriceRange 获取最近3个周期最低价格
func Get3CyclePriceRange(symbol string, period string) (float64, float64, error) {
	klines, err := GetKlines(symbol, period, 3)
	if err != nil {
		return 0.0, 0.0, err
	}

	lowerPrices := []float64{}
	highPrices := []float64{}
	for _, kl := range klines {
		lowerPrices = append(lowerPrices, kl.Low)
		highPrices = append(highPrices, kl.High)
	}
	sort.Float64s(lowerPrices)
	sort.Float64s(highPrices)

	lowerPrice, highPrice := 0.0, 0.0
	lowLen := len(lowerPrices)
	if lowLen > 1 {
		lowerPrice = (lowerPrices[0] + lowerPrices[1]) / 2
	} else if lowLen == 1 {
		lowerPrice = lowerPrices[0]
	} else {
		return 0.0, 0.0, errors.New("lowPrices is nil")
	}

	highLen := len(highPrices)
	if highLen == 3 {
		highPrice = (highPrices[1] + highPrices[2]) / 2
	} else if highLen == 2 {
		highPrice = (highPrices[0] + highPrices[1]) / 2
	} else if highLen == 1 {
		highPrice = highPrices[0]
	} else {
		return 0.0, 0.0, errors.New("highPrices is nil")
	}

	return lowerPrice, highPrice, nil
}
