package model

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/app/grid-strategy-srv/util/goex"
	"wq-fotune-backend/app/grid-strategy-srv/util/goex/binance"
	"wq-fotune-backend/app/grid-strategy-srv/util/grid"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi"

	"strconv"

	"github.com/globalsign/mgo/bson"
	"github.com/zhufuyi/pkg/krand"
	"github.com/zhufuyi/pkg/logger"
)

// GetLatestPrice 获取交易所品种的最新价格
func GetLatestPrice(exchange string, symbol string) (float64, error) {
	latestPrice := 0.0
	err := error(nil)
	// 区分不同交易所
	switch exchange {
	case ExchangeHuobi:
		latestPrice, err = huobi.GetLatestPrice(symbol)
		if err != nil {
			logger.Warn("huobi.GetLatestPrice error", logger.Err(err))
			return latestPrice, fmt.Errorf("从交易所%s获取%s最新价格失败", exchange, symbol)
		}
	case ExchangeBinance:
		latestPrice, err = binance.GetLatestPrice(symbol, env.ProxyAddr)
		if err != nil {
			logger.Warn("binance.GetLatestPrice error", logger.Err(err))
			return latestPrice, fmt.Errorf("从交易所%s获取%s最新价格失败", exchange, symbol)
		}

	default:
		return latestPrice, fmt.Errorf("暂时不支持%s交易所", exchange)
	}

	return latestPrice, nil
}

// CalculateProfit 计算利润
func CalculateProfit(exchange string, grids []*grid.Grid) (float64, float64) {
	// 区分不同交易所
	fees := 0.0
	switch exchange {
	case ExchangeHuobi:
		fees = huobi.FilledFees
	case ExchangeBinance:
		fees = binance.FilledFees
	}
	return grid.CalculateProfit(grids, fees)
}

// GetExchangeRule 获取交易所规则参数
func GetExchangeRule(exchange string, symbol string) (float64, float64) {
	limitVolume := 10.0
	el := GetExchangeLimitCache(GetKey(exchange, symbol))
	if el.VolumeLimit > 0.0 {
		limitVolume = el.VolumeLimit
	}

	fees := 0.001
	// 区分不同交易所
	switch exchange {
	case ExchangeBinance:
		fees = binance.FilledFees
	case ExchangeHuobi:
		fees = huobi.FilledFees
	}

	return limitVolume, fees
}

// CalculateBestGrid 计算出给定参数中最好的网格策略参数
func CalculateBestGrid(exchange string, symbol string, totalSum float64, latestPrice float64, minProfitRate float64) (*GridFilter, error) {
	if minProfitRate < 0.0005 {
		minProfitRate = 0.0005
	}

	cg := &CalculateGrid{
		Exchange:         exchange,
		Symbol:           symbol,
		TargetProfitRate: minProfitRate,
		ParamsRange: &SymbolParams{
			TotalSum:    &ValueRange{totalSum, totalSum, 1},
			LatestPrice: &ValueRange{latestPrice, latestPrice, 1},
		},
	}

	limitVolume, fees := GetExchangeRule(exchange, symbol)
	gf := cg.Done(limitVolume, fees)
	if gf == nil {
		return nil, errors.New("not found")
	}

	return gf, nil
}

// GenNewGrid 生成新的网格
func GenNewGrid(exchange string, symbol string, totalSum float64, gridNum int, minPrice float64, maxPrice float64) ([]*grid.Grid, error) {
	key := GetKey(exchange, symbol)
	elc := GetExchangeLimitCache(key)

	grids, _ := grid.Generate(
		grid.GSGrid,
		minPrice,
		maxPrice,
		totalSum,
		gridNum,
		elc.PricePrecision,
		elc.QuantityPrecision,
	)
	if err := grid.IsValidGrids(grids, elc.VolumeLimit); err != nil {
		return []*grid.Grid{}, fmt.Errorf("网格价格范围太小，建议价格范围在%v以上", int(elc.VolumeLimit)*gridNum+5)
	}

	return grids, nil
}

//  撤单
func cancelLimitOrder(gs *GridStrategy, exchangeAccount goex.Accounter) ([]string, error) {
	if exchangeAccount == nil {
		return []string{}, errors.New("not found account")
	}

	// 获取所有委托订单
	query := bson.M{"gsid": gs.ID, "orderState": huobi.OrderStateSubmitted}
	field := bson.M{"orderID": true, "_id": true, "symbol": true}
	gtrs, err := FindGridTradeRecords(query, field, 0, 100)
	if err != nil {
		return nil, err
	}

	successOrderIDs := []string{}
	failedOrderIDs := []string{}
	for _, v := range gtrs {
		err = exchangeAccount.CancelOrder(v.OrderID, v.Symbol)
		if err != nil {
			isCancelel := false
			// 区分不同交易所
			switch gs.Exchange {
			case ExchangeHuobi:
				if strings.Contains(err.Error(), "order-orderstate-error") {
					isCancelel = true
				}
			case ExchangeBinance:
				if strings.Contains(err.Error(), "Unknown order sent") {
					isCancelel = true
				}
			}

			if isCancelel {
				query = bson.M{"_id": v.ID}
				update := bson.M{"$set": bson.M{"orderState": huobi.OrderStateFilled, "stateTime": time.Now()}}
				UpdateGridTradeRecord(query, update)
			} else {
				failedOrderIDs = append(failedOrderIDs, v.OrderID)
				logger.Warn("cancel order failed", logger.Err(err), logger.String("clientOrderID", v.ClientOrderID), logger.String("orderID", v.OrderID))
				continue
			}
		}

		successOrderIDs = append(successOrderIDs, v.OrderID)
		time.Sleep(time.Millisecond * 50)
	}

	if len(failedOrderIDs) > 0 {
		logger.Error("取消失败的订单", logger.String("uid", gs.UID), logger.String("gsid", gs.ID.Hex()), logger.Any("failedOrderIDs", failedOrderIDs))
	}

	return successOrderIDs, nil
}

// 更新订单状态
func updateLocalOrderStatus(orderIDs []string, status string) []string {
	failedOrderIDs := []string{}

	update := bson.M{
		"$set": bson.M{"orderState": status, "stateTime": time.Now()},
	}
	for _, orderID := range orderIDs {
		query := bson.M{"orderID": orderID}
		err := UpdateGridTradeRecord(query, update)
		if err != nil {
			failedOrderIDs = append(failedOrderIDs, orderID)
			logger.Error("updateGridTradeRecord error", logger.Err(err), logger.String("orderID", orderID), logger.Any("update", update))
			continue
		}
	}

	return failedOrderIDs
}

// 计算交易币种的当前持仓量
func calculateCurrencyPosition(gsid string) float64 {
	buyQuantity, selQuantity := 0.0, 0.0

	query := bson.M{"gsid": bson.ObjectIdHex(gsid), "orderState": huobi.OrderStateFilled}
	field := bson.M{"side": true, "quantity": true}

	count, _ := CountGridTradeRecords(query)
	limit := 100
	page := count / 100

	for i := 0; i <= page; i++ {
		gtr, err := FindGridTradeRecords(query, field, i, limit)
		if err != nil {
			continue
		}
		for _, v := range gtr {
			if v.Side == "buy" {
				buyQuantity += v.Quantity
			} else if v.Side == "sell" {
				selQuantity += v.Quantity
			}
		}
	}

	return buyQuantity - selQuantity
}

// CalculatePosition 计算当前持仓量
func CalculatePosition(exchange string, symbol string, gsid string, exchangeAccount goex.Accounter) (float64, error) {
	if exchangeAccount == nil {
		return 0.0, errors.New("not found account")
	}

	totalPosition := calculateCurrencyPosition(gsid)
	needClosePosition := getCurrencyPosition(exchange, symbol, totalPosition, exchangeAccount)

	el := GetExchangeLimitCache(GetKey(exchange, symbol))
	positionSize := FloatRoundOff(needClosePosition, el.QuantityPrecision)

	logger.Infof("gsid=%s, calculate position=%v, need close position %v, quantityPrecision=%d", gsid, totalPosition, positionSize, el.QuantityPrecision)

	return positionSize, nil
}

//  获取货币可撤单的仓位
func getCurrencyPosition(exchange string, symbol string, totalPosition float64, exchangeAccount goex.Accounter) float64 {
	needClosePosition := totalPosition
	feesRate := 0.0
	// 区分不同交易所
	switch exchange {
	case ExchangeHuobi:
		feesRate = huobi.FilledFees
	case ExchangeBinance:
		feesRate = binance.FilledFees
	}

	currency, _ := goex.SplitSymbol(symbol)
	if currency != "" {
		time.Sleep(time.Millisecond * 200)
		currencyBalance, err := exchangeAccount.GetCurrencyBalance(currency)
		if err != nil {
			logger.Warn("GetCurrencyBalance error", logger.Err(err), logger.String("currency", currency))
		} else {
			if currencyBalance < totalPosition {
				needClosePosition = currencyBalance * (1 - feesRate)
			} else {
				needClosePosition = totalPosition * (1 - feesRate)
			}
		}
	}

	return needClosePosition
}

// 检测账号下持仓量和资金余额是否能够满足网格
func checkAccountBalance(exchange string, symbol string, latestPrice float64, exchangeAccount goex.Accounter, grids []*grid.Grid, currencyPosition float64) (int, float64, error) {
	if exchangeAccount == nil {
		return 0, 0.0, errors.New("not found account")
	}

	basisGridNO := 0
	needBuyCoinQuantity := 0.0
	needMoney := 0.0
	err := error(nil)

	buyCoin := 0.0
	for k, v := range grids {
		if v.Price > latestPrice { // 大于基准线价格，需要卖出币的数量，也就是账号下必须已经持币数量
			basisGridNO = k // 网格编号是有序的，最后一个大于basisPrice对应编号
			needBuyCoinQuantity += v.SellQuantity
		} else { // 小于等于基准线价格，统计委托挂单需要的金额
			needMoney += v.Price * v.BuyQuantity
			buyCoin += v.BuyQuantity
		}
	}
	needBuyCoinQuantity -= grids[basisGridNO].SellQuantity // 减去和当前价格相近的卖单数量，不需要挂卖单

	// 区分不同交易所
	switch exchange {
	case ExchangeHuobi:
		needBuyCoinQuantity = needBuyCoinQuantity/(1-huobi.FilledFees) + buyCoin*huobi.FilledFees // 扣除手续费之后的持仓数量
	case ExchangeBinance:
		needBuyCoinQuantity = needBuyCoinQuantity/(1-binance.FilledFees) + buyCoin*binance.FilledFees // 扣除手续费之后的持仓数量
	}

	// 查询账户基准币的余额，判断余额是否满足网格所需的金额
	_, anchorSymbol := goex.SplitSymbol(symbol)
	time.Sleep(time.Millisecond * 200)
	anchorCurrencyBalance, err := exchangeAccount.GetCurrencyBalance(anchorSymbol)
	if err != nil {
		return basisGridNO, needBuyCoinQuantity, fmt.Errorf("get currency balance error, err=%s", err.Error())
	}
	needTotalMoney := needMoney + needBuyCoinQuantity*latestPrice
	currencyBalance := currencyPosition * latestPrice // 把持仓数量转换为余额
	logger.Info("check account balance", logger.Float64("gridNeedTotalMoney", needTotalMoney), logger.Float64(fmt.Sprintf("%s Balance", anchorSymbol), anchorCurrencyBalance), logger.Float64("currencyBalance", currencyBalance))
	if anchorCurrencyBalance+currencyBalance-needTotalMoney < 0.0 {
		//return basisGridNO, needBuyCoinQuantity, fmt.Errorf("%s balance(%v) is less than requisite money(%v)", anchorSymbol, anchorCurrencyBalance+currencyBalance, needTotalMoney)
		return basisGridNO, needBuyCoinQuantity, fmt.Errorf("%s余额%v少于网格需要的金额%v", anchorSymbol, anchorCurrencyBalance+currencyBalance, needTotalMoney)
	}

	key := GetKey(exchange, symbol)
	el := GetExchangeLimitCache(key)
	needBuyCoinQuantity = FloatRound(needBuyCoinQuantity, el.QuantityPrecision)

	return basisGridNO, needBuyCoinQuantity, nil
}

// BuyMarketOrder 买入市价单
func BuyMarketOrder(exchange string, symbol string, gsid string, latestPrice float64, needBuyCoinQuantity float64, exchangeAccount goex.Accounter) (*TradeRecord, error) {
	marketOrderRecord := &TradeRecord{}

	key := GetKey(exchange, symbol)
	el := GetExchangeLimitCache(key)
	if needBuyCoinQuantity*latestPrice < el.VolumeLimit {
		return &TradeRecord{}, nil
	}

	if needBuyCoinQuantity > 0.0 {
		if exchangeAccount == nil {
			return marketOrderRecord, errors.New("not found account")
		}

		side := "buy"
		clientOrderID := GenerateClientOrderID(exchange, PrefixIDMob, gsid)
		volume := 0.0
		// 区分不同交易所
		switch exchange {
		case ExchangeHuobi:
			volume = FloatRound(needBuyCoinQuantity*latestPrice, el.PricePrecision) // 火币需要转换为总额
		case ExchangeBinance:
			volume = needBuyCoinQuantity
		}
		//coinQuantity := fmt.Sprintf("%v", volume)
		coinQuantity := Float64ToStr(volume, el.QuantityPrecision)

		orderID, err := exchangeAccount.PlaceMarketOrder(side, symbol, coinQuantity, clientOrderID)
		if err != nil {
			return marketOrderRecord, err
		}

		logger.Info("place market order success",
			logger.String("side", side),
			logger.String("symbol", symbol),
			logger.Float64("price", latestPrice),
			logger.String("coinQuantity", coinQuantity),
			logger.String("clientOrderID", clientOrderID),
			logger.String("orderID", orderID),
		)

		marketOrderRecord = &TradeRecord{
			OrderID:       orderID,
			ClientOrderID: clientOrderID,
			OrderType:     "market",
			Side:          side,
			Price:         latestPrice,
			Quantity:      needBuyCoinQuantity,
			Volume:        volume,
			Unit:          goex.GetAnchorCurrency(symbol),
			OrderState:    "filled",
		}
	} else {
		logger.Info("no need to place market order")
	}

	return marketOrderRecord, nil
}

// SellMarketOrder 卖出市价单
func SellMarketOrder(exchange string, symbol string, gsid string, latestPrice float64, needSellCoinQuantity float64, exchangeAccount goex.Accounter) (*GridTradeRecord, error) {
	marketOrderRecord := &GridTradeRecord{}

	key := GetKey(exchange, symbol)
	el := GetExchangeLimitCache(key)
	if needSellCoinQuantity*latestPrice < el.VolumeLimit {
		return &GridTradeRecord{}, nil
	}

	if needSellCoinQuantity > 0.0 {
		if exchangeAccount == nil {
			return marketOrderRecord, errors.New("not found account")
		}

		side := "sell"
		clientOrderID := GenerateClientOrderID(exchange, PrefixIDMos, string(krand.String(krand.R_All, 6)))
		//coinQuantity := fmt.Sprintf("%v", needSellCoinQuantity)
		coinQuantity := Float64ToStr(needSellCoinQuantity, el.QuantityPrecision)
		volume := needSellCoinQuantity
		latestPrice := 0.0
		feesRate := 0.0
		err := error(nil)

		// 区分不同交易所
		switch exchange {
		case ExchangeHuobi:
			latestPrice, err = huobi.GetLatestPrice(symbol)
			if err != nil {
				logger.Warn("huobi.GetLatestPrice error", logger.Err(err), logger.String("symbol", symbol))
			}
			feesRate = huobi.FilledFees

		case ExchangeBinance:
			latestPrice, err = binance.GetLatestPrice(symbol, env.ProxyAddr)
			if err != nil {
				logger.Warn("binance.GetLatestPrice error", logger.Err(err), logger.String("symbol", symbol))
			}
			feesRate = binance.FilledFees
		}

		orderID, err := exchangeAccount.PlaceMarketOrder(side, symbol, coinQuantity, clientOrderID)
		if err != nil {
			logger.Warn("PlaceMarketOrder error", logger.Err(err), logger.String("params", fmt.Sprintf("%s, %s, %v, %s", side, symbol, coinQuantity, clientOrderID)))
			return marketOrderRecord, err
		}

		logger.Info("place market order success",
			logger.String("side", side),
			logger.String("symbol", symbol),
			logger.Float64("price", FloatRound(latestPrice, el.PricePrecision)),
			logger.String("coinQuantity", coinQuantity),
			logger.String("clientOrderID", clientOrderID),
			logger.String("orderID", orderID),
		)

		//if volume == needSellCoinQuantity {
		//	volume = FloatRound(needSellCoinQuantity*latestPrice, el.GetVolumePrecision())
		//}
		volume = FloatRound(needSellCoinQuantity*latestPrice, el.PricePrecision)

		_, anchorSymbol := goex.SplitSymbol(symbol)

		marketOrderRecord = &GridTradeRecord{
			GSID: bson.ObjectIdHex(gsid),
			GID:  0,

			OrderID:       orderID,
			ClientOrderID: clientOrderID,
			OrderType:     "market",
			Side:          side,
			Price:         latestPrice,
			Quantity:      needSellCoinQuantity,
			Volume:        volume,
			Unit:          anchorSymbol,
			Fees:          volume * feesRate,

			OrderState: "filled",
			StateTime:  time.Now(),

			IsStartUpOrder: false,

			Exchange: exchange,
			Symbol:   symbol,
		}

	} else {
		logger.Info("no need to place market order")
	}

	return marketOrderRecord, nil
}

// 买入、卖出网格限价单
func placeGridLimitOrder(gridType int, exchange string, symbol string, gsid string, exchangeAccount goex.Accounter, grids []*grid.Grid, basisGridNO int) ([]*TradeRecord, error) {
	records := []*TradeRecord{}

	if exchangeAccount == nil {
		return records, errors.New("not found account")
	}

	var side, price, amount, clientOrderID, orderID string
	var err error
	for k, grid := range grids {
		if k == basisGridNO { // 忽略接近当前价格的卖单
			continue
		}
		quantity, volume := 0.0, 0.0
		if k > basisGridNO {
			// 买入限价单
			side = "buy"
			clientOrderID = NewGridClientOrderID(exchange, side, gsid, k)
			amount = fmt.Sprintf("%v", grid.BuyQuantity)
			quantity = grid.BuyQuantity
			volume = grid.Price * grid.BuyQuantity
		} else {
			// 卖出限价单
			side = "sell"
			clientOrderID = NewGridClientOrderID(exchange, side, gsid, k)
			amount = fmt.Sprintf("%v", grid.SellQuantity)
			quantity = grid.SellQuantity
			volume = grid.Price * grid.SellQuantity
		}

		//price = fmt.Sprintf("%v", grid.Price)
		price = strconv.FormatFloat(grid.Price, 'f', -1, 64)

		orderID, err = exchangeAccount.PlaceLimitOrder(side, symbol, price, amount, clientOrderID)
		if err != nil {
			logger.Error("placeLimitOrder error", logger.Err(err), logger.String("param", fmt.Sprintf("%v, %v, %v, %v, %v", side, symbol, price, amount, clientOrderID)))
			continue
		}

		grid.OrderID = orderID
		grid.Side = side
		logger.Info("place limit order success",
			logger.String("side", side),
			logger.String("symbol", symbol),
			logger.String("price", price),
			logger.String("amount", amount),
			logger.String("clientOrderID", clientOrderID),
			logger.String("orderID", orderID),
		)

		records = append(records, &TradeRecord{
			GID:           k,
			OrderID:       orderID,
			ClientOrderID: clientOrderID,
			OrderType:     "limit",
			Side:          side,
			Price:         grid.Price,
			Quantity:      quantity,
			Volume:        volume,
			Unit:          goex.GetAnchorCurrency(symbol),
			OrderState:    huobi.OrderStateSubmitted,
		})

		if gridType == GridTypeInfinite {
			time.Sleep(time.Millisecond * 300) // 延时下单
		} else {
			time.Sleep(time.Millisecond * 50)
		}
	}

	return records, nil
}

// 网格策略更新的字段
func updateGridStrategyToDB(exchange string, gsid string, minPrice float64, maxPrice float64, latestPrice float64, buyCoinQuantity float64, grids []*grid.Grid, basisGridNO int) error {
	averageProfit, averageProfitRate := CalculateProfit(exchange, grids)

	query := bson.M{"_id": bson.ObjectIdHex(gsid)}
	update := bson.M{
		"$set": bson.M{
			"minPrice":          minPrice,
			"maxPrice":          maxPrice,
			"basisPrice":        latestPrice,
			"buyCoinQuantity":   buyCoinQuantity,
			"gridNum":           len(grids) - 1,
			"gridBaseNO":        basisGridNO,
			"averageProfit":     averageProfit,
			"averageProfitRate": averageProfitRate,
		},
		"$inc": bson.M{"resetPriceCount": 1},
	}

	return UpdateGridStrategy(query, update)
}

// 更新网格数据
func updateGridPendingOrderToDB(gsid string, grids []*grid.Grid, basisGridNO int) error {
	query := bson.M{"gsid": bson.ObjectIdHex(gsid)}
	update := bson.M{"$set": bson.M{
		"grids":       grids,
		"basisGridNO": basisGridNO,
	}}

	return UpdateGridPendingOrder(query, update)
}

// 保存订单记录
func saveNewOrderRecords(exchange string, symbol string, gsid string, records []*TradeRecord) error {
	errStr := ""
	now := time.Now()
	isStartUpOrder := false

	for _, v := range records {
		if v.OrderID == "" {
			continue
		}
		if v.Side == "sell" {
			isStartUpOrder = true
		}

		gtr := &GridTradeRecord{
			GSID: bson.ObjectIdHex(gsid),
			GID:  v.GID,

			OrderID:       v.OrderID,
			ClientOrderID: v.ClientOrderID,
			OrderType:     v.OrderType,
			Side:          v.Side,
			Price:         v.Price,
			Quantity:      v.Quantity,
			Volume:        v.Volume,
			Unit:          v.Unit,
			//Fees : ,

			OrderState: v.OrderState,
			StateTime:  now,

			IsStartUpOrder: isStartUpOrder,
			//BuyOrderID: ,
			//GridPeerBuyOrderID  : ,

			Exchange: exchange,
			Symbol:   symbol,
		}

		gtr.ID = bson.NewObjectId()

		err := gtr.Insert()
		if err != nil {
			errStr += fmt.Sprintf("%s, orderID=%v ||", err.Error(), v.OrderID)
		}
	}

	if errStr != "" {
		return errors.New(errStr)
	}

	return nil
}

// 更新缓存
func updateGridProcessCache(uid string, exchange string, symbol string, gsid string, grids []*grid.Grid, basisGridNO int) {
	key := StrategyCacheKey(uid, exchange, symbol, gsid)
	gp, ok := GetStrategyCache(key)
	if ok {
		gp.Grids = grids
		gp.BasisGridNO = basisGridNO
		SetStrategyCache(key, gp)
		logger.Info("update gridProcess cache success", logger.String("key", key))
		return
	}

	logger.Error("not found gridProcess cache", logger.String("key", key))
}

// 更新网格数据和缓存
func updateGridAndCache(uid string, exchange string, symbol string, gsid string, minPrice float64, maxPrice float64, latestPrice float64, buyCoinQuantity float64, grids []*grid.Grid, basisGridNO int, records []*TradeRecord) {
	err := updateGridStrategyToDB(exchange, gsid, minPrice, maxPrice, latestPrice, buyCoinQuantity, grids, basisGridNO)
	if err != nil {
		logger.Warn("updateGridStrategyToDB error", logger.Err(err), logger.String("gsid", gsid))
	}

	err = updateGridPendingOrderToDB(gsid, grids, basisGridNO)
	if err != nil {
		logger.Warn("updateGridPendingOrderToDB error", logger.Err(err), logger.String("gsid", gsid))
	}

	err = saveNewOrderRecords(exchange, symbol, gsid, records)
	if err != nil {
		logger.Warn("saveNewOrderRecords error", logger.Err(err), logger.String("gsid", gsid))
	}

	updateGridProcessCache(uid, exchange, symbol, gsid, grids, basisGridNO)
}

// UpdateGridStrategy 更新运行中网格策略
func UpdateRunningGridStrategy(gs *GridStrategy, exchangeAccount goex.Accounter, grids []*grid.Grid, minPrice float64, maxPrice float64, latestPrice float64, isClosePosition bool) error {
	uid, exchange, symbol, gsid := gs.UID, gs.Exchange, gs.Symbol, gs.ID.Hex()

	// 撤单
	successOrderIDs, err := cancelLimitOrder(gs, exchangeAccount)
	if err != nil {
		logger.Error("cancelLimitOrder error", logger.Err(err), logger.String("gsid", gsid))
		return fmt.Errorf("%s 网格策略撤单失败", gsid)
	}
	if len(successOrderIDs) > 0 {
		logger.Info("已取消的委托订单", logger.Any("successOrderIDs", successOrderIDs))
	}

	// 更新本地订单记录状态
	failedUpdateOrderIDs := updateLocalOrderStatus(successOrderIDs, huobi.OrderStateCanceled)
	if len(failedUpdateOrderIDs) > 0 {
		logger.Error("本地更新状态失败的订单", logger.String("gsid", gsid), logger.Any("failedUpdateOrderIDs", failedUpdateOrderIDs))
	} else {
		logger.Info("本地更新所有订单状态成功", logger.String("gsid", gsid), logger.Any("orderIDs", successOrderIDs))
	}

	// 计算当前持仓
	currencyPosition, err := CalculatePosition(exchange, symbol, gsid, exchangeAccount)
	if err != nil {
		logger.Error("CalculatePosition error", logger.Err(err), logger.String("gsid", gsid))
	}

	// 判断是否平仓
	if isClosePosition {
		logger.Info("start to close position", logger.String("params", fmt.Sprintf("%v,%v,%v,%v,%v", exchange, symbol, gsid, latestPrice, currencyPosition)))
		gtr, err := SellMarketOrder(exchange, symbol, gsid, latestPrice, currencyPosition, exchangeAccount)
		if err != nil {
			logger.Error("SellMarketOrder error", logger.Err(err), logger.String("gsid", gsid), logger.Float64("latestPrice", latestPrice), logger.Float64("currencyPosition", currencyPosition))
		} else {
			currencyPosition = 0.0
			gtr.Insert()
		}
	}

	// 检测账号下持仓量和资金余额是否能够满足网格所需金额
	basisGridNO, needBuyCoinQuantity, err := checkAccountBalance(exchange, symbol, latestPrice, exchangeAccount, grids, currencyPosition)
	if err != nil {
		logger.Warn("账号余额不足", logger.Err(err), logger.String("form", fmt.Sprintf("%v,%v,%v,%v", uid, exchange, symbol, latestPrice)))
		return err
	}
	buyCoinQuantity := needBuyCoinQuantity
	needBuyCoinQuantity -= currencyPosition // 减去网格策略原有的持仓数量
	logger.Info("账号余额满足网格策略", logger.String("uid", uid), logger.Float64("needBuyCoinQuantity", needBuyCoinQuantity))

	// 买入市价单，如果返回的needBuyCoinQuantity为0，则忽略
	marketRecord, err := BuyMarketOrder(exchange, symbol, gsid, latestPrice, needBuyCoinQuantity, exchangeAccount)
	if err != nil {
		logger.Warn("买入市价单失败", logger.Err(err), logger.Float64("needBuyCoinQuantity", needBuyCoinQuantity), logger.String("gsid", gsid))
		return errors.New("买入市价单失败")
	}

	// 买入、卖出网格限价单
	limitRecords, err := placeGridLimitOrder(gs.Type, exchange, symbol, gsid, exchangeAccount, grids, basisGridNO)
	if err != nil {
		logger.Warn("买入、卖出网格限价单失败", logger.Err(err), logger.Any("grids", grids), logger.String("gsid", gsid))
		return errors.New("委托买卖订单失败")
	}
	if len(limitRecords) == len(grids)-1 {
		logger.Info("买入、卖出网格所有限价单成功", logger.String("gsid", gsid), logger.Int("total", len(limitRecords)))
	} else {
		logger.Warn("有部分买入、卖出网格限价单失败", logger.String("gsid", gsid), logger.String("success/total", fmt.Sprintf("%d/%d", len(limitRecords), len(grids)-1)))
	}
	records := append(limitRecords, marketRecord)

	// 更新网格数据和缓存
	updateGridAndCache(uid, exchange, symbol, gsid, minPrice, maxPrice, latestPrice, buyCoinQuantity, grids, basisGridNO, records)

	return nil
}
