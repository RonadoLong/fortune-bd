package evaluation

import (
	"github.com/shopspring/decimal"
	"strings"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/pkg/utils"
)

type Okex struct {
	Symbol      string
	ShortVolume float64
	OpenPrice   float64
	Commission  float64
	ClosePrice  float64
	Principal   float64 //本金
	Profit      float64
	Direction   string
	Type        string
}

func NewOKex(symbol string, pos, openPrice, commission, closePrice, balance float64) *Okex {
	return &Okex{
		Symbol:      symbol,
		ShortVolume: pos,
		OpenPrice:   openPrice,
		Commission:  commission,
		ClosePrice:  closePrice,
		Principal:   balance,
		Profit:      0,
		Direction:   "",
		Type:        "",
	}
}

const (
	BTC = "BTC"
	ETH = "ETH"
	LTC = "LTC"
	EOS = "EOS"
	BCH = "BCH"
)

func getSymbol(symbol string) string {
	if strings.Contains(symbol, BTC) {
		return BTC
	}
	if strings.Contains(symbol, ETH) {
		return ETH
	}
	if strings.Contains(symbol, LTC) {
		return LTC
	}
	if strings.Contains(symbol, EOS) {
		return EOS
	}
	if strings.Contains(symbol, BCH) {
		return BCH
	}
	return ""
}

func getContractUsdtSize(symbol string) float64 {
	switch symbol {
	case BTC:
		return 0.01
	case ETH:
		return 0.1
	case LTC:
		return 1
	case EOS:
		return 10
	case BCH:
		return 0.1
	}
	return 0
}

func (o *Okex) CalCommission() float64 {
	return o.Commission
}

//
//USDT合约的收益公式
//多仓：收益=（平仓价-开仓价）*面值*张数
//空仓：收益=（开仓价-平仓价）*面值*张数

//USDT合约面值：0.01BTC，10EOS，0.1ETH，1LTC，0.1BCH，100XRP，10ETC，1BSV，1000TRX
//币本位合约面值：BTC=100USD 其他币种=10USD
func (o *Okex) CalProfit() float64 {
	logger.Infof("CalProfit %+v", o)
	shortVolume := decimal.NewFromFloat(o.ShortVolume) //1 假设 1手
	closePrice := decimal.NewFromFloat(o.ClosePrice)   //9450
	openPrice := decimal.NewFromFloat(o.OpenPrice)     //9350
	profit := 0.0
	defer func() {
		o.Profit = profit
	}()

	if o.ClosePrice == 0 || o.OpenPrice == 0 {
		return 0
	}
	//合约交易
	// btc * 100 USDT 币本位 合约
	if o.Type == ContractTradingCoin {
		if strings.Contains(strings.ToUpper(o.Symbol), BTC) {
			shortVolume = shortVolume.Mul(decimal.NewFromInt(100))
		} else {
			shortVolume = shortVolume.Mul(decimal.NewFromInt(10))
		}

		if o.Direction == CloseBuy {
			profit, _ = shortVolume.Div(openPrice).Sub(shortVolume.Div(closePrice)).Float64()
			return profit
		}
		// o.Direction == CloseSell
		profit, _ = shortVolume.Div(closePrice).Sub(shortVolume.Div(openPrice)).Float64()
		return profit
	}
	// usdt 本位 合约
	if o.Type == ContractTradingUsdt {
		symbol := getSymbol(o.Symbol)
		size := getContractUsdtSize(symbol)
		sizeDecimal := decimal.NewFromFloat(size)
		if o.Direction == CloseBuy {
			profit, _ = closePrice.Sub(openPrice).Mul(sizeDecimal).Mul(shortVolume).Float64()
			return profit
		}
		profit, _ = openPrice.Sub(closePrice).Mul(sizeDecimal).Mul(shortVolume).Float64()
		return profit
	}
	return profit
}

func (o *Okex) RateReturn() float64 {
	principal := decimal.NewFromFloat(o.Principal) //20000 usdt
	profit := decimal.NewFromFloat(o.Profit)
	rateReturn, _ := profit.Div(principal).Float64() // 1000 / 20000 *100
	return utils.Keep2Decimal(rateReturn * 100)
}
