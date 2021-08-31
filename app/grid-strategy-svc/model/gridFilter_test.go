package model

import (
	"fmt"
	"testing"
	"wq-fotune-backend/app/grid-strategy-srv/util/goex/binance"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi"

	"github.com/globalsign/mgo/bson"
	"github.com/json-iterator/go"
	"github.com/k0kubun/pp"
)

func TestGridFilter_Insert(t *testing.T) {
	cg := &CalculateGrid{
		Exchange:             "binance",
		Symbol:               "ethusdt",
		TargetProfitRate:     0.001,
		PriceDifferenceLimit: 50,

		ParamsRange: &SymbolParams{
			IntervalPrice: &ValueRange{
				Start: 1,
				End:   10,
				Step:  1,
			},
			LatestPrice: &ValueRange{
				Start: 350,
				End:   400,
				Step:  2,
			},
		},
	}

	file := "C:\\Work\\Golang\\Project\\src\\wq-grid-strategy\\config\\" + fmt.Sprintf("%s-%s.csv", cg.Exchange, cg.Symbol)
	data := make(chan *GridFilter)
	go Save2CSV(file, data)
	cg.DoneAndSave(data)
	close(data)
}

func TestMgo2CSV(t *testing.T) {
	cg := CalculateGrid{
		ParamsRange: &SymbolParams{
			IntervalPrice: &ValueRange{},
			LatestPrice:   &ValueRange{},
		},
	}
	js, _ := jsoniter.Marshal(cg)
	fmt.Printf("%s\n", js)
	return

	file := "C:\\Work\\Golang\\Project\\src\\wq-grid-strategy\\config\\eth.csv"
	err := mgo2csv(file)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestFindGridFilters(t *testing.T) {
	exchange := "binance"
	symbol := "ethusdt"
	totalSum := 1001.0

	if totalSum > 10000 {
		totalSum = 10000
	}

	latestPrice := 0.0
	// 区分不同交易所
	switch exchange {
	case "huobi":
		latestPrice, _ = huobi.GetLatestPrice(symbol)
	case "binance":
		latestPrice, _ = binance.GetLatestPrice(symbol, "socks5://127.0.0.1:10808")
	}

	fmt.Println(totalSum, latestPrice)

	query := bson.M{
		//"exchange":     exchange,
		//"totalSum":     int(totalSum/10) * 10,  // 投资金额取10的倍数
		//"currentPrice": int(latestPrice/5) * 5, // 根据品种获取值
		"uniqueVal": fmt.Sprintf("%s.%v.%v", exchange, int(totalSum/10)*10, int(latestPrice/5)*5),
	}

	//intervalPrice:=1.0

	gf, err := FindGridFilter(query, bson.M{})
	if err != nil {
		t.Error(query, err)
		return
	}

	pp.Println(gf)
}

func TestFilterBestGrid(t *testing.T) {
	gfs := gridFilters{
		{
			TotalSum:             80,
			CurrentPrice:         200,
			GridNum:              5,
			AverageIntervalPrice: 3.8,
		},
		{
			TotalSum:             80,
			CurrentPrice:         200,
			GridNum:              6,
			AverageIntervalPrice: 3.8,
		},
		{
			TotalSum:             80,
			CurrentPrice:         200,
			GridNum:              7,
			AverageIntervalPrice: 3.8,
		},
		{
			TotalSum:             80,
			CurrentPrice:         200,
			GridNum:              8,
			AverageIntervalPrice: 3.8,
		},
	}

	v := getBestGrid(gfs)
	fmt.Println(v.TotalSum, v.CurrentPrice, v.GridNum, v.AverageIntervalPrice)
}
