package evaluation

type Evaluation interface {
	CalCommission() float64
	CalProfit() float64
	RateReturn() float64
}

//profit := decimal.NewFromFloat(0)
//if closePrice.Equal(decimal.NewFromInt(0)) || openPrice.Equal(decimal.NewFromInt(0)) {
//return profit
//}
//if _, ok := acMap[strings.ToUpper(symbol)]; ok { //币币交易利润
//if direction == Buy {
////  卖出价格减去买进价格乘以币数 除以当前价格 得出多少个币
//profit = closeVolume.Mul(closePrice).Sub(closeVolume.Mul(openPrice)).Div(closePrice)
//} else {
//profit = closeVolume.Mul(openPrice).Sub(closeVolume.Mul(closePrice)).Div(closePrice)
//
//}
//return profit
//}
//
//if _, ok := btcMapSwap[strings.ToUpper(symbol)]; ok {
//closeVolume = closeVolume.Mul(decimal.NewFromInt(100))
//}
//if _, ok := ethMapSwap[strings.ToUpper(symbol)]; ok {
//closeVolume = closeVolume.Mul(decimal.NewFromInt(10))
//}
//
//if direction == Buy { // 平多
//profit = closeVolume.Div(openPrice).Sub(closeVolume.Div(closePrice))
//} else { //平空
//profit = closeVolume.Div(closePrice).Sub(closeVolume.Div(openPrice))
//}
//return profit
