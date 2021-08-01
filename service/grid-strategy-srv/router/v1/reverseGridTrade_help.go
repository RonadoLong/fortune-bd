package v1

import (
	"errors"
	"fmt"
	"wq-fotune-backend/service/grid-strategy-srv/model"
	"wq-fotune-backend/service/grid-strategy-srv/util/goex"
	"wq-fotune-backend/service/grid-strategy-srv/util/goex/binance"
	"wq-fotune-backend/service/grid-strategy-srv/util/grid"
	"wq-fotune-backend/service/grid-strategy-srv/util/huobi"

	"github.com/zhufuyi/logger"

	"time"


	"github.com/globalsign/mgo/bson"
)

type reverseGridForm struct {
	Exchange string `json:"exchange"`
	Symbol   string `json:"symbol"`
}

func (r *reverseGridForm) valid() error {
	switch "" {
	case r.Exchange:
		return errors.New("参数exchange为空")
	case r.Symbol:
		return errors.New("参数symbol为空")
	}

	return nil
}

func (r *reverseGridForm) getMinMoney() (float64, string) {
	el := model.GetExchangeLimitCache(model.GetKey(r.Exchange, r.Symbol))
	tradeCurrency := goex.GetTradeCurrency(r.Symbol)

	price, err := model.GetLatestPrice(r.Exchange, tradeCurrency+"usdt")
	if err != nil {
		logger.Error("model.GetLatestPrice error", logger.Err(err), logger.Any("form", r))
	}
	if price <= 0 {
		return FloatRound(el.VolumeLimit * 20), "usdt"
	}

	return FloatRound(55.0/price, el.QuantityPrecision), tradeCurrency
}

// 初始化反向网格
func (g *gridTradeForm) initReverseGrid() ([]*grid.Grid, error) {
	el := model.GetExchangeLimitCache(model.GetKey(g.Exchange, g.Symbol))

	grids, err := grid.Generate(
		g.GridIntervalType,
		g.MinPrice,
		g.MaxPrice,
		g.TotalSum*g.LatestPrice, // 转换为基准币
		g.GridNum,
		el.PricePrecision,
		el.QuantityPrecision,
	)
	if err != nil {
		return grids, err
	}

	err = g.checkGrids(grids) // 检查网格

	return grids, err
}

// 检测账号下持仓量是否能够满足网格
func (g *gridTradeForm) calculateTradeCurrencySize(grids []*grid.Grid) (float64, error) {
	basisPrice := g.BasisPrice
	gridBasisNO := 0
	needBuyCoin := 0.0
	needSellCoin := 0.0

	for k, v := range grids {
		if v.Price > basisPrice { // 大于基准线价格，需要卖出币的数量，也就是账号下必须已经持币数量
			gridBasisNO = k // 网格编号是有序的，最后一个大于basisPrice对应编号
			needBuyCoin += v.SellQuantity
		} else { // 小于等于基准线价格，统计委托挂单需要的金额
			needSellCoin += v.BuyQuantity
		}
	}

	needBuyCoin -= grids[gridBasisNO].SellQuantity // 去掉接近当前价格的卖单

	// 区分不同交易所
	switch g.Exchange {
	case model.ExchangeHuobi:
		needSellCoin = needSellCoin / (1 - huobi.FilledFees) // 扣除手续费之后的持仓数量
	case model.ExchangeBinance:
		needSellCoin = needSellCoin / (1 - binance.FilledFees) // 扣除手续费之后的持仓数量
	}

	if g.exchangeAccount == nil {
		return needSellCoin, errors.New("not found account")
	}

	el := model.GetExchangeLimitCache(model.GetKey(g.Exchange, g.Symbol))

	// 查询账户基准币的余额，判断余额是否满足网格所需的金额
	tradeCurrency := goex.GetTradeCurrency(g.Symbol)
	tradeCurrencyBalance, err := g.exchangeAccount.GetCurrencyBalance(tradeCurrency)
	if err != nil {
		return needSellCoin, fmt.Errorf("get currency balance error, err=%s", err.Error())
	}
	needTotalCoin := needSellCoin + needBuyCoin
	logger.Info("check account balance", logger.Float64("gridNeedTotalCoin", needTotalCoin), logger.Float64(fmt.Sprintf("%s Balance", tradeCurrency), tradeCurrencyBalance), logger.Int("quantityPrecision", el.QuantityPrecision))
	if tradeCurrencyBalance-needTotalCoin < 0.0 {
		return needSellCoin, fmt.Errorf("%s balance(%f) is less than grid need coin(%f)", tradeCurrency, tradeCurrencyBalance, needTotalCoin)
	}

	g.GridBasisNO = gridBasisNO
	g.sellCoinQuantity = FloatRound(needSellCoin, el.QuantityPrecision)

	return g.sellCoinQuantity, nil
}

// 计算计价币种的当前持仓量
func calculateAnchorCurrencyPosition(gsid string) float64 {
	buyQuantity, sellQuantity := 0.0, 0.0

	query := bson.M{"gsid": bson.ObjectIdHex(gsid), "orderState": huobi.OrderStateFilled}
	field := bson.M{"side": true, "quantity": true}

	count, _ := model.CountGridTradeRecords(query)
	limit := 100
	page := count / 100

	for i := 0; i <= page; i++ {
		gtr, err := model.FindGridTradeRecords(query, field, i, limit)
		if err != nil {
			continue
		}
		for _, v := range gtr {
			if v.Side == "buy" {
				buyQuantity += v.Quantity
			} else if v.Side == "sell" {
				sellQuantity += v.Quantity
			}
		}
	}

	return sellQuantity - buyQuantity
}

//  获取货币可撤单的仓位
func getAnchorCurrencyPosition(exchange string, symbol string, totalPosition float64, exchangeAccount goex.Accounter) float64 {
	needClosePosition := totalPosition
	feesRate := 0.0
	// 区分不同交易所
	switch exchange {
	case model.ExchangeHuobi:
		feesRate = huobi.FilledFees
	case model.ExchangeBinance:
		feesRate = binance.FilledFees
	}

	latestPrice, err := model.GetLatestPrice(exchange, symbol)
	if err != nil {
		logger.Warnf("获取交易所%s的交易对%s价格失败", exchange, symbol)
		return needClosePosition * (1 - feesRate)
	}

	currency := goex.GetAnchorCurrency(symbol)
	if currency != "" {
		time.Sleep(time.Millisecond * 200)
		balance, err := exchangeAccount.GetCurrencyBalance(currency)
		if err != nil {
			logger.Warn("GetCurrencyBalance error", logger.Err(err), logger.String("currency", currency))
		} else {
			currencyBalance := balance / latestPrice // 转换为币的数量
			if currencyBalance < totalPosition {
				needClosePosition = currencyBalance * (1 - feesRate)
			} else {
				needClosePosition = totalPosition * (1 - feesRate)
			}
		}
	}

	return needClosePosition * (1 - feesRate)
}

// 计算当前持仓量
func calculatePosition(exchange string, symbol string, gsid string, exchangeAccount goex.Accounter) (float64, error) {
	if exchangeAccount == nil {
		return 0.0, errors.New("not found account")
	}

	totalPosition := calculateAnchorCurrencyPosition(gsid)
	needClosePosition := getAnchorCurrencyPosition(exchange, symbol, totalPosition, exchangeAccount)

	el := model.GetExchangeLimitCache(model.GetKey(exchange, symbol))
	positionSize := FloatRound(needClosePosition, el.QuantityPrecision)
	logger.Infof("gsid=%s, calculate volume=%v, need close volume %v", gsid, totalPosition, positionSize)

	return positionSize, nil
}

// 反向网格平仓
func (c *cancelGridStrategyForm) closeReverseGridPosition() error {
	needClosePosition, err := calculatePosition(c.Exchange, c.Symbol, c.Gsid, c.exchangeAccount)
	if err != nil {
		logger.Warn("model.CalculatePosition error", logger.Err(err))
		return err
	}

	if needClosePosition > 0 {
		//latestPrice, err := model.GetLatestPrice(c.Exchange, c.Symbol)
		//if err != nil {
		//	return fmt.Errorf("获取交易所%s的交易对%s价格失败", c.Exchange, c.Symbol)
		//}
		//// 由volume转换为实际数量quantity
		//needClosePosition = needClosePosition / latestPrice

		// 平仓
		tradeRecord, err := c.placeMarketOrder("buy", needClosePosition)
		if err != nil {
			return err
		}

		if tradeRecord.Exchange != "" {
			err = tradeRecord.Insert()
			if err != nil {
				logger.Error("tradeRecord.Insert error", logger.Err(err), logger.Any("tradeRecord", tradeRecord))
				return err
			}
		}
	}
	return nil
}

func getReverseMinAndMaxTotalSum(exchange string, symbol string) (float64, float64, error) {
	el := model.GetExchangeLimitCache(model.GetKey(exchange, symbol))
	tradeCurrency := goex.GetTradeCurrency(symbol)

	price, err := model.GetLatestPrice(exchange, tradeCurrency+"usdt")
	if err != nil {
		logger.Error("model.GetLatestPrice error", logger.Err(err), logger.String("exchange.symbol", exchange+"."+tradeCurrency+"usdt"))
		return 0, 0, err
	}

	return FloatRound(55.0/price, el.QuantityPrecision), FloatRound(10000.0/price, el.QuantityPrecision), nil
}
