package cron

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/global"
	"wq-fotune-backend/libs/logger"
	api "wq-fotune-backend/libs/okex_client"
)

type Ticker struct {
	Symbol string  `json:"symbol"`
	Last   float64 `json:"last"`
	Buy    float64 `json:"buy"`
	Open   float64 `json:"open"`
	Sell   float64 `json:"sell"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Vol    float64 `json:"vol"`
	Change string  `json:"change"`
	Date   uint64  `json:"date"` // 单位:ms
}

func (tk Ticker) MarshalBinary() ([]byte, error) {
	return json.Marshal(tk)
}

var TickOKex = "tick:okex"
var TickOkexAll = "tick:okex:all"
var RateKey = "rate:usd-rmb"

func StoreOkexTick() {
	d := time.Second * 10
	t := time.NewTicker(d)
	defer t.Stop()
	for {
		<-t.C
		storeOkexTick()
		//logger.Infof("lens %d", len(OkexTickArray))
		if len(OkexTickMap) == 0 {
			continue
		}
		OkexTickMapAll["USDT-USDT"] = &Ticker{
			Last: 1,
			Buy:  1,
			Open: 1,
			High: 1,
			Low:  1,
		}
		if err := global.RedisCli.HMSet(TickOKex, OkexTickMap).Err(); err != nil {
			logger.Errorf("将行情存到redis失败 %v", err)
		}
		if err := global.RedisCli.HMSet(TickOkexAll, OkexTickMapAll).Err(); err != nil {
			logger.Errorf("将行情存到redis失败 %v", err)
		}
		storeRateRmb()
		//for k, v := range OkexTickMap {
		//    logger.Infof("%v--%v", k, v.Symbol)
		//}
	}
}

var OkexTickMap = make(map[string]interface{})
var OkexTickMapAll = make(map[string]interface{}) //all保存所有品种

var OkexTickArray = make([]Ticker, 0)
var OkexTickArrayAll = make([]Ticker, 0)
var OkexTickArrayBtc = make([]Ticker, 0)

func storeOkexTick() {
	client := api.InitClient("", "", "")
	tickers, err := client.APIClient.OKExSpot.GetAllTicker()
	if err != nil {
		logger.Infof("storeOkexTick OKExSpot.GetTicker has err %v", err)
		return
	}
	//重置为空
	OkexTickArray = OkexTickArray[:0]
	OkexTickArrayAll = OkexTickArrayAll[:0]
	OkexTickArrayBtc = OkexTickArrayBtc[:0]
	for _, v := range tickers {
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
		if strings.HasSuffix(v.Symbol, "-USDT") && v.Low >= 1 {
			//logger.Info(global.StructToJsonStr(tick))
			OkexTickMap[v.Symbol] = tick
			OkexTickArray = append(OkexTickArray, tick)
		}
		if strings.HasSuffix(v.Symbol, "-USDT") {
			OkexTickMapAll[v.Symbol] = tick
			OkexTickArrayAll = append(OkexTickArrayAll, tick)
		}
		if len(v.Symbol) >= 4 && strings.HasSuffix(v.Symbol, "BTC") {
			OkexTickArrayBtc = append(OkexTickArrayBtc, tick)
		}

	}
	//OkexTickArray = tickers
}

type QuoteRate struct {
	InstrumentID string `json:"instrument_id"`
	Rate         string `json:"rate"`
	Timestamp    string `json:"timestamp"`
}

func (rate QuoteRate) MarshalBinary() ([]byte, error) {
	return json.Marshal(rate)
}

func storeRateRmb() {
	request, err := http.NewRequest("get", "https://www.okex.com/api/futures/v3/rate", nil)
	if err != nil {
		logger.Warnf("获取法币汇率接口错误 %v", err)
		return
	}
	response, err := env.GetProxyHttpClient().Do(request)
	if err != nil {
		logger.Warnf("获取法币汇率接口错误 %v", err)
		return
	}
	defer response.Body.Close()
	all, _ := ioutil.ReadAll(response.Body)
	var rate QuoteRate
	if err := json.Unmarshal(all, &rate); err != nil {
		logger.Warnf("解析法币汇率数据错误 %v %s", err, string(all))
		return
	}
	if err := global.RedisCli.Set(RateKey, rate, 0).Err(); err != nil {
		logger.Warnf("保存汇率失败 %v", err)
		return
	}
}
