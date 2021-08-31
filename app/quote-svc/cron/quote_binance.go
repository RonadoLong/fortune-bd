package cron

import (
	"fmt"
	"fortune-bd/libs/cache"
	"fortune-bd/libs/exchangeclient"
	"fortune-bd/libs/logger"
	"strings"
	"time"

)

var TickBinanceAll = "tick:binance:all"
var BinanceTickArrayAll = make([]Ticker, 0)
var BinanceTickMapAll = make(map[string]interface{}) //all保存所有品种
var BinaceTickArrayBtc = make([]Ticker, 0)

func StoreBinanceTick() {
	d := time.Second * 8
	t := time.NewTicker(d)
	defer t.Stop()
	for {
		<-t.C
		storeBinanceTick()
		if len(BinanceTickMapAll) == 0 {
			continue
		}
		if err := cache.Redis().HMSet(TickBinanceAll, BinanceTickMapAll).Err(); err != nil {
			logger.Errorf("binance将行情存到redis失败 %v", err)
		}
	}
}

func storeBinanceTick() {
	client := exchangeclient.InitBinance("", "")
	tickers, err := client.ApiClient.GetTickers()
	if err != nil {
		logger.Infof("storeBinanceTick GetTickers has err %v", err)
		return
	}
	//重置为空
	BinanceTickArrayAll = BinanceTickArrayAll[:0]
	BinaceTickArrayBtc = BinaceTickArrayBtc[:0]
	for _, v := range tickers {
		if v.Open <= 0 {
			continue
		}
		change := (v.Last - v.Open) / v.Open * 100
		tick := Ticker{
			Symbol: v.Symbol,
			Last:   v.Last,
			Buy:    v.Buy,
			Open:   v.Open,
			Sell:   v.Sell,
			High:   v.High,
			Low:    v.Low,
			Vol:    v.Vol,
			Change: fmt.Sprintf("%.2f%v", change, "%"),
			Date:   v.Date,
		}
		if change > 0.0 {
			tick.Change = "+" + tick.Change
		}
		if strings.HasSuffix(v.Symbol, "USDT") {
			tick.Symbol = strings.ReplaceAll(v.Symbol, "USDT", "-USDT")
			BinanceTickMapAll[tick.Symbol] = tick
			BinanceTickArrayAll = append(BinanceTickArrayAll, tick)
		}
		if len(v.Symbol) >= 4 && strings.HasSuffix(v.Symbol, "BTC") {
			tick.Symbol = strings.ReplaceAll(v.Symbol, "BTC", "-BTC")
			BinaceTickArrayBtc = append(BinaceTickArrayBtc, tick)
		}
	}
	logger.Infof("获取USDT行情数据数量: %d", len(BinanceTickArrayAll))
}
