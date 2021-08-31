package grid

import (
	"fmt"
	"math"
	"testing"
	"wq-fotune-backend/app/grid-strategy-srv/util/goex/binance"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi"
)

func TestGenerate(t *testing.T) {
	totalSum := 100.0
	gridNum := 16
	limitVolume := 5.0

	minPrice, maxPrice := getMinMax(12000, gridNum)

	grids, _ := Generate(GSGrid, minPrice, maxPrice, totalSum, gridNum, 2, 6)
	if err := IsValidGrids(grids, limitVolume); err != nil {
		fmt.Println(err.Error())
		return
	}

	fees := 0.002
	averageProfit, averageProfitRate := CalculateProfit(grids, fees)
	fmt.Printf("averageProfit=%v, averageProfitRate=%v%%\n\n", averageProfit, FloatRound(averageProfitRate*100))

	PrintFormat(grids)
}

func TestGetMinMaxPrice2(t *testing.T) {
	feesMap := map[string]float64{"huobi": 0.002, "binance": 0.001}
	limitVolumeMap := map[string]float64{"huobi": 5, "binance": 10}

	exchange := "huobi"
	symbol := "ethusdt"
	minPriceDifferenceMap := map[string]float64{"btcusdt": 150, "ethusdt": 20}
	fees := feesMap[exchange]

	latestPrice := 0.0
	// 区分不同交易所
	switch exchange {
	case "huobi":
		latestPrice, _ = huobi.GetLatestPrice(symbol)
	case "binance":
		latestPrice, _ = binance.GetLatestPrice(symbol, "socks5://127.0.0.1:10808")
	}

	totalSum := 10000.0 // 投资金额
	gridNum := 30       // 网格数量

	//minPriceDifference := 150.0 // 网格最小和最大间距差的最小值
	minPriceDifference := minPriceDifferenceMap[symbol]
	//latestPrice = 320

	//fmt.Println(gridNum, latestPrice)
	for intervalPrice := 1.0; intervalPrice <= 10; intervalPrice += 0.5 { // 网格价格间隔轮询
		minPrice, maxPrice := GetMinMax(latestPrice, intervalPrice, gridNum, false)
		priceDifference := maxPrice - minPrice
		if priceDifference < minPriceDifference {
			continue
		}

		grids, _ := Generate(GSGrid, minPrice, maxPrice, totalSum, gridNum, 8, 6)
		if err := IsValidGrids(grids, limitVolumeMap[exchange]); err != nil {
			continue
		}

		averageProfit, averageProfitRate := CalculateProfit(grids, fees)
		if averageProfit < 0.0 {
			continue
		}

		if averageProfitRate > 0.001 {
			fmt.Println(totalSum, gridNum, intervalPrice, minPrice, latestPrice, maxPrice, priceDifference, averageProfitRate, averageProfit)
			break
		}
	}
}

func TestGrid(t *testing.T) {
	totalSum := 1000.0
	gridNum := 10
	limitVolume := 10.0

	minPrice, maxPrice := 300.0, 400.0

	grids, _ := Generate(GSGrid, minPrice, maxPrice, totalSum, gridNum, 0, 3)
	if err := IsValidGrids(grids, limitVolume); err != nil {
		fmt.Println(err.Error())
		return
	}

	PrintFormat(grids)
}

func TestGeneratePQ(t *testing.T) {
	minPrice := 280.0 // 最近三个月最小价格
	maxPrice := minPrice * 2
	q := 1.006
	totalSum := 2000.0

	num := int(math.Log10(maxPrice/350) / math.Log10(q))

	grids := GenerateGS(350, q, totalSum, num, 2, 6)
	//grids := geometricSequenceGrid3(350, q, totalSum, num, 2, 6)
	PrintFormat(grids)
}

func TestFloatRound(t *testing.T) {
	pi := 3.1415926535
	for i := 0; i < 12; i++ {
		fmt.Println(FloatRound(pi, i))
	}
}

func BenchmarkFloatRound(b *testing.B) {
	pi := 3.1415926535
	for i := 0; i < b.N; i++ {
		FloatRound(pi, 6)
	}
}
