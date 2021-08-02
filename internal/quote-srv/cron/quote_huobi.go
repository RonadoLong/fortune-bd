package cron

import (
	"fmt"
	"log"
	"strings"
	"time"
	global2 "wq-fotune-backend/libs/global"
	api "wq-fotune-backend/libs/huobi_client"
	"wq-fotune-backend/libs/logger"
)

var TickHuobiAll = "tick:huobi:all"
var HuobiTickArrayAll = make([]Ticker, 0)
var HuobiTickMapAll = make(map[string]interface{}) //all保存所有品种
var HuobiTickArrayBtc = make([]Ticker, 0)

func StoreHuobiTick() {
	d := time.Second * 9
	t := time.NewTicker(d)
	defer t.Stop()
	for {
		<-t.C
		storeHuobiTick()
		if len(HuobiTickMapAll) == 0 {
			log.Println("为0")
			continue
		}
		if err := global2.RedisCli.HMSet(TickHuobiAll, HuobiTickMapAll).Err(); err != nil {
			logger.Errorf("huobi将行情存到redis失败 %v", err)
		}
	}
}

func storeHuobiTick() {
	client := api.InitClient("", "", false)
	tickers, err := client.APIClient.GetTickers()
	if err != nil {
		logger.Infof("storeHuobiTick GetTickers has err %v", err)
		return
	}
	//重置为空
	HuobiTickArrayAll = HuobiTickArrayAll[:0]
	HuobiTickArrayBtc = HuobiTickArrayBtc[:0]
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
			HuobiTickMapAll[tick.Symbol] = tick
			HuobiTickArrayAll = append(HuobiTickArrayAll, tick)
		}
		if len(v.Symbol) >= 4 && strings.HasSuffix(v.Symbol, "BTC") {
			tick.Symbol = strings.ReplaceAll(v.Symbol, "BTC", "-BTC")
			HuobiTickArrayBtc = append(HuobiTickArrayBtc, tick)
		}
	}
}
