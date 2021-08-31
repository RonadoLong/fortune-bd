package model

import (
	"context"
	"errors"
	"fmt"
	"fortune-bd/app/grid-strategy-svc/util/goex/binance"
	"fortune-bd/app/grid-strategy-svc/util/grid"
	"fortune-bd/app/grid-strategy-svc/util/huobi"
	"strconv"
	"strings"
	"time"


	"github.com/globalsign/mgo/bson"
	"github.com/zhufuyi/pkg/krand"
	"github.com/zhufuyi/pkg/logger"
)

// 同步本地和交易所的实际委托订单
func (g *GridProcess) syncGridStrategy() error {
	defer func() {
		if e := recover(); e != nil {
			logger.Error("syncGridStrategy panic", logger.Any("e", e))
		}
	}()

	if g.ExchangeAccount == nil {
		return errors.New("not found account")
	}

	latestPrice, err := GetLatestPrice(g.Exchange, g.Symbol)
	if err != nil {
		logger.Warn("GetLatestPrice error", logger.Err(err), logger.String("exchange", g.Exchange), logger.String("symbol", g.Symbol))
		return fmt.Errorf("获取%s交易所%s的最新价格失败", g.Exchange, g.Symbol)
	}

	// 判断策略是否在运行状态
	query := bson.M{"_id": bson.ObjectIdHex(g.Gsid)}
	gs, err := FindGridStrategy(query, bson.M{})
	if err != nil {
		return err
	}
	if !gs.IsRun {
		return fmt.Errorf("grid strategy %s had stoped", g.Gsid)
	}

	switch gs.Type {
	case GridTypeTrend: // 如果是趋势网格，达到阈值后自动生成新网格
		if isNeedGenNewGrid(gs, latestPrice) {
			gf, err := CalculateBestGrid(gs.Exchange, gs.Symbol, gs.TotalSum, latestPrice, 0.001)
			if err != nil {
				logger.Warn("CalculateBestGrid error", logger.Err(err), logger.String("params", fmt.Sprintf("%v,%v,%v,%v", gs.Exchange, gs.Symbol, gs.TotalSum, latestPrice)))
				return err
			}

			grids, err := GenNewGrid(gs.Exchange, gs.Symbol, gs.TotalSum, gf.GridNum, gf.MinPrice, gf.MaxPrice)
			if err != nil {
				logger.Warn("GenNewGrid error", logger.Err(err), logger.String("params", fmt.Sprintf("%v,%v,%v,%v,%v,%v", gs.Exchange, gs.Symbol, gs.TotalSum, gs.GridNum, gf.MinPrice, gf.MaxPrice)))
				return err
			}
			gs.GridNum = gf.GridNum

			return UpdateRunningGridStrategy(gs, g.ExchangeAccount, grids, gf.MinPrice, gf.MaxPrice, latestPrice, false)
		}

	case GridTypeInfinite:
		moveStep := checkNeedUpdateGrid(gs, latestPrice)

		if moveStep != 0 {
			el := GetExchangeLimitCache(GetKey(g.Exchange, g.Symbol))
			moveGridSize := int(gs.GridNum / 3)
			minPrice := gs.MinPrice
			switch moveStep {
			case -1:
				if gs.IntervalSize <= 0 {
					return errors.New("intervalSize is valid")
				}
				for i := 0; i < moveGridSize; i++ {
					minPriceTmp := minPrice / gs.IntervalSize
					if minPriceTmp <= gs.StartupMinPrice {
						break
					}
					minPrice = minPriceTmp
				}

			case 1:
				minPrice = g.Grids[gs.GridNum-moveGridSize].Price
			}

			grids := grid.GenerateGS(FloatRound(minPrice, el.PricePrecision), gs.IntervalSize, gs.TotalSum, gs.GridNum, el.PricePrecision, el.QuantityPrecision)
			return UpdateRunningGridStrategy(gs, g.ExchangeAccount, grids, grids[gs.GridNum].Price, grids[0].Price, latestPrice, false)
		}
	}

	// 获取交易所历史订单记录
	historyOrders, err := g.ExchangeAccount.GetHistoryOrdersInfo(g.Symbol, huobi.OrderStateSubmitted, "")
	if err != nil {
		return err
	}

	// 交易所历史订单记录中筛选出属于当前策略的记录
	submittedOrderMap := map[string]bool{} // key为订单id
	// 区分不同交易所
	switch g.Exchange {
	case ExchangeHuobi:
		hos := historyOrders.([]*huobi.HistoryOrders)
		for _, ho := range hos {
			if strings.Contains(ho.ClientOrderID, cutGsid(g.Gsid)) { // 根据策略id筛选
				key := fmt.Sprintf("%v", ho.ID)
				submittedOrderMap[key] = true
			}
		}

	case ExchangeBinance:
		hos := historyOrders.([]*binance.OrderInfo)
		for _, ho := range hos {
			if strings.Contains(ho.ClientOrderID, cutGsid(g.Gsid)) { // 根据策略id筛选
				key := fmt.Sprintf("%v", ho.ID)
				submittedOrderMap[key] = true
			}
		}
	}

	// 判断交易所实际挂单数和网格数是否一致
	if len(submittedOrderMap) >= len(g.Grids)-1 {
		logger.Info("no limit order need to sync", logger.String("gsid", g.Gsid), logger.Int("submittedOrderSize", len(submittedOrderMap)))
		return nil
	}
	logger.Infof("strategy %s match %d submitted Order records, need %d", g.Gsid, len(submittedOrderMap), len(g.Grids)-1)

	// 读取本地网格挂单记录
	query = bson.M{"gsid": bson.ObjectIdHex(g.Gsid), "orderState": "submitted"}
	gtrs, err := FindGridTradeRecords(query, bson.M{}, 0, 200, "-gid")
	if err != nil {
		return err
	}

	var (
		tradeTime     int64
		tradeAmount   float64
		orderID       string
		clientOrderID string
		orderStatus   string
	)

	// 比较本地和交易所的委托订单，如果状态不存在，修改本地委托记录，并添加新委托订单
	for _, gtr := range gtrs {
		if _, ok := submittedOrderMap[gtr.OrderID]; ok {
			//logger.Infof("submitted order %s already exist, gridNO=%d", gtr.OrderID, gtr.GID)
			continue
		}
		logger.Info("record is inconsistent", logger.Any("record", gtr))

		// 获取当前订单在交易所的信息
		orderInfo, err := g.ExchangeAccount.GetOrderInfo(gtr.OrderID, gtr.Symbol)
		if err != nil {
			logger.Error("GetOrderInfo error", logger.Err(err), logger.String("orderID", gtr.OrderID))
			switch gtr.Exchange {
			case ExchangeHuobi:
				if strings.Contains(err.Error(), "invalid") {
					updateOrderState(gtr.OrderID, "invalid")
				}
			case ExchangeBinance:
				if strings.Contains(err.Error(), "Unknown") {
					updateOrderState(gtr.OrderID, "unknown")
				}
			}

			continue
		}

		// 区分不同交易所
		switch g.Exchange {
		case ExchangeHuobi:
			oi := orderInfo.(*huobi.OrderInfo)
			if gtr.OrderState == oi.State {
				continue
			}

			if oi.State != huobi.OrderStateFilled {
				if oi.State == huobi.OrderStateCanceled { // 在官网手动取消了订单
					updateOrderState(fmt.Sprintf("%v", oi.ID), huobi.OrderStateCanceled)
				}
				continue
			}
			tradeTime = oi.FinishedAt
			tradeAmount = str2Float64(oi.Amount)
			orderID = fmt.Sprintf("%v", oi.ID)
			clientOrderID = oi.ClientOrderId
			orderStatus = oi.State

		case ExchangeBinance:
			oi := orderInfo.(*binance.OrderInfo)
			if gtr.OrderState == oi.State {
				continue
			}

			if oi.State != binance.OrderStateFilled {
				if oi.State == binance.OrderStateCanceled { // 在官网手动取消的订单
					updateOrderState(oi.ID, binance.OrderStateCanceled)
				}
				continue
			}
			tradeTime = oi.UpdateAt
			tradeAmount = str2Float64(oi.Amount)
			orderID = oi.ID
			clientOrderID = oi.ClientOrderID
			orderStatus = oi.State

		default:
			logger.Errorf("unsupported exchange %s", g.Exchange)
			continue
		}

		logger.Info("exchange respond data", logger.String("data", fmt.Sprintf("%s, %v, %v, %v, %v, %v", g.Exchange, orderID, clientOrderID, tradeAmount, tradeTime, orderStatus)))
		time.Sleep(time.Microsecond * 50) // 延时

		// 同步订单记录
		err = g.syncLimitOrder(latestPrice, submittedOrderMap, tradeTime, tradeAmount, orderID, clientOrderID, orderStatus)
		if err != nil {
			logger.Warn("syncLimitOrder error",
				logger.Err(err),
				logger.String("param", fmt.Sprintf("%v, %v, %v, %v, %v, %v", submittedOrderMap, tradeTime, tradeAmount, orderID, clientOrderID, orderStatus)))
			continue
		}

		logger.Info("syncLimitOrder success", logger.String("param", fmt.Sprintf("%v,%v,%v,%v,%v", tradeTime, tradeAmount, orderID, clientOrderID, orderStatus)))
	}

	err = g.checkMissGrid(latestPrice)
	if err != nil {
		logger.Error("checkMissGrid error", logger.Float64("latestPrice", latestPrice))
	}

	return nil
}

// 判断是否需要新的网格
func isNeedGenNewGrid(gs *GridStrategy, latestPrice float64) bool {
	gsid := gs.ID.Hex()

	minLimitPrice := gs.EntryPrice * 0.99
	//currentMinLimitPrice := gs.MinPrice * 0.99
	maxLimitPrice := gs.MaxPrice * 1.01

	// 禁止网格往下移动
	if minLimitPrice <= latestPrice && latestPrice <= maxLimitPrice {
		logger.Infof("market price is in the range, ignore handle, gsid=%s,  latestPrice = %v, entryLimitPrice=%v,maxLimitPrice =%v", gsid, latestPrice, minLimitPrice, maxLimitPrice)
		return false
	}

	// 网格可以往下移动
	//if currentMinLimitPrice < latestPrice && latestPrice < maxLimitPrice {
	//	logger.Infof("price is in the range, ignore handle, gsid=%s, currentMinLimitPrice = %v, latestPrice = %v, maxLimitPrice =%v", gsid, currentMinLimitPrice, latestPrice, maxLimitPrice)
	//	return false
	//}

	if latestPrice < minLimitPrice {
		logger.Infof("latest price(%v) had fallen to the entry price(%v), ignore handle, gsid=%s, symbol=%s", latestPrice, minLimitPrice, gsid, gs.Symbol)
		return false
	}

	nowSecond, latestRecordSecond := time.Now().Local().Unix(), int64(0)
	// 最近一次记录
	gtrs, err := FindGridTradeRecords(bson.M{"gsid": gs.ID}, bson.M{}, 0, 1)
	if err != nil {
		logger.Warn("FindGridTradeRecords error", logger.Err(err))
		return false
	}
	if len(gtrs) > 0 {
		latestRecordSecond = gtrs[0].CreatedAt.Local().Unix()
		if nowSecond-latestRecordSecond < 3600 {
			logger.Infof("gsid=%s, symbol=%s, latest price %v is out of limit price(%v ~ %v), but not more than 1 hours(%d, %d), ignore processing.", gsid, gs.Symbol, latestPrice, minLimitPrice, maxLimitPrice, nowSecond, latestRecordSecond)
			return false
		}
	}
	logger.Info("trend grid will generate new grid", logger.String("gsid", gsid), logger.String("symbol", gs.Symbol), logger.String("oldParam", fmt.Sprintf("%v,%v,%v,%v,%v", latestPrice, minLimitPrice, maxLimitPrice, nowSecond, latestRecordSecond)))

	return true
}

// 判断是否需要更新的网格，返回值：-1表示往下移动网格，0表示不需要更新网格，1表示网上移动网格
func checkNeedUpdateGrid(gs *GridStrategy, latestPrice float64) int {
	gsid := gs.ID.Hex()
	moveStep := 0

	minLimitPrice := gs.MinPrice
	maxLimitPrice := gs.MaxPrice

	if gs.MinPrice <= latestPrice && latestPrice <= gs.MaxPrice {
		logger.Infof("market price is in the range, ignore handle, gsid=%s,  latestPrice = %v, entryLimitPrice=%v,maxLimitPrice =%v", gsid, latestPrice, minLimitPrice, maxLimitPrice)
		return moveStep
	}

	// 当前行情如果跌破启动的最小价格，禁止网格继续往下移动
	if latestPrice < gs.StartupMinPrice {
		logger.Infof("latest price(%v) had fallen to the min price(%v), ignore handle, gsid=%s, symbol=%s", latestPrice, minLimitPrice, gsid, gs.Symbol)
		return moveStep
	}

	// 需要网下移动网格
	if gs.StartupMinPrice <= latestPrice && latestPrice <= gs.MinPrice {
		moveStep = -1
	}

	// 需要网上移动网格
	if gs.MaxPrice < latestPrice {
		moveStep = 1
	}

	nowSecond, latestRecordSecond := time.Now().Local().Unix(), int64(0)
	// 最近一次记录
	gtrs, err := FindGridTradeRecords(bson.M{"gsid": gs.ID}, bson.M{}, 0, 1)
	if err != nil {
		logger.Warn("FindGridTradeRecords error", logger.Err(err))
		return 0
	}
	if len(gtrs) > 0 {
		latestRecordSecond = gtrs[0].CreatedAt.Local().Unix()
		if nowSecond-latestRecordSecond < 3600 {
			logger.Infof("infinite grid, time is not up yet. gsid=%s, symbol=%s, latest price %v is out of limit price(%v ~ %v), but not more than 1 hours(%d, %d), ignore processing.", gsid, gs.Symbol, latestPrice, minLimitPrice, maxLimitPrice, nowSecond, latestRecordSecond)
			return 0
		}
	}
	logger.Info("infinite grid will update grid", logger.String("gsid", gsid), logger.String("symbol", gs.Symbol), logger.String("oldParam", fmt.Sprintf("%v,%v,%v,%v,%v", latestPrice, minLimitPrice, maxLimitPrice, nowSecond, latestRecordSecond)))

	return moveStep
}

// 同步委托订单
func (g *GridProcess) syncLimitOrder(latestPrice float64, submittedOrderMap map[string]bool, tradeTime int64, tradeAmount float64, orderID, clientOrderID, orderStatus string) error {
	var newRecord *TradeRecord
	var err error

	_, gridNO := parseClientOrderID(g.Exchange, clientOrderID)
	grid := g.Grids[gridNO]

	if gridNO == 0 && grid.Price < latestPrice {
		logger.Info("latestPrice has exceeded the maxPrice")
		g.Grids[gridNO].Side = ""
		g.Grids[gridNO].OrderID = ""
		g.BasisGridNO = 0
	} else if gridNO == len(g.Grids)-1 && grid.Price > latestPrice {
		logger.Info("latestPrice less than the minPrice")
		g.Grids[gridNO].Side = ""
		g.Grids[gridNO].OrderID = ""
		g.BasisGridNO = len(g.Grids) - 1
	} else {
		// 获取当前基准价格线
		for _, v := range g.Grids {
			if v.Price > latestPrice {
				g.BasisGridNO = v.GID
			}
		}

		// 如果刚好在基准网格线，不需要添加新的委托订单
		if g.BasisGridNO == grid.GID {
			logger.Debug("do not need add new limit order", logger.Int("BasisGridNO", g.BasisGridNO))
			g.Grids[gridNO].Side = ""
			g.Grids[gridNO].OrderID = ""
		} else {
			newRecord, err = g.placeLimitOrder(gridNO, latestPrice)
			if err != nil {
				logger.Error("g.placeLimitOrder error", logger.String("gsid", g.Gsid), logger.String("symbol", g.Symbol), logger.Any("grid", grid))
				newRecord = nil
			}

			if newRecord != nil && newRecord.OrderID != "" {
				submittedOrderMap[newRecord.OrderID] = true
			}
		}
	}

	g.updateOldOrderRecord(gridNO, orderID, orderStatus, tradeTime, tradeAmount) // 更新旧的委托订单状态

	return g.updateGridAndAddRecord(gridNO, newRecord) // 更新网格参数和添加新网格
}

// 同步数据库和缓存网格数据
func (g *GridProcess) syncGridCache() {
	// 读取本地网格挂单记录
	query := bson.M{"gsid": bson.ObjectIdHex(g.Gsid), "orderState": "submitted"}
	gtrs, err := FindGridTradeRecords(query, bson.M{}, 0, 100, "-gid")
	if err != nil {
		logger.Error("FindGridTradeRecords error", logger.Err(err), logger.Any("query", query))
		return
	}

	currentOrderMap := map[int]*GridTradeRecord{}
	for _, gtr := range gtrs {
		currentOrderMap[gtr.GID] = gtr
	}
	isNeedSync := false
	for _, grid := range g.Grids {
		gridNO := grid.GID
		if gtr, ok := currentOrderMap[gridNO]; ok {
			if grid.OrderID != gtr.OrderID {
				g.Grids[gridNO].OrderID = gtr.OrderID
				g.Grids[gridNO].Side = gtr.Side
				isNeedSync = true
			}
		} else { // 不一致更新到数据库
			g.Grids[gridNO].OrderID = ""
			g.Grids[gridNO].Side = ""
			isNeedSync = true
		}

		if isNeedSync {
			isNeedSync = false
			query := bson.M{"gsid": bson.ObjectIdHex(g.Gsid)}
			update := bson.M{
				"$set": bson.M{
					fmt.Sprintf("grids.%d.side", gridNO):    g.Grids[gridNO].Side,
					fmt.Sprintf("grids.%d.orderId", gridNO): g.Grids[gridNO].OrderID,
				},
			}
			err := UpdateGridPendingOrder(query, update)
			if err != nil {
				logger.Error("UpdateGridPendingOrder error", logger.Err(err), logger.Any("query", query), logger.Any("update", update))
			}
		}
	}
}

// 检查网格是否有遗漏委托订单
func (g *GridProcess) checkMissGrid(latestPrice float64) error {
	g.syncGridCache()

	missGridNOs := []int{}
	existGridNOs := []int{}

	for _, grid := range g.Grids {
		if grid.Price > latestPrice {
			g.BasisGridNO = grid.GID
		}

		if grid.OrderID != "" { // 如果不为空，说明已经有委托单了
			existGridNOs = append(existGridNOs, grid.GID)
			continue
		}

		if grid.GID == 0 && grid.Price < latestPrice { // 网格最高价格小于最新价格，除了第0个网格为空，其他网格必须存在
			g.BasisGridNO = 0
			continue
		} else if grid.GID == len(g.Grids)-1 && grid.Price > latestPrice { // 网格最低价格大于最新价格，除了最后一个网格为空，其他网格必须存在
			g.BasisGridNO = len(g.Grids) - 1
			continue
		}
		missGridNOs = append(missGridNOs, grid.GID)
	}

	isExistBasisGridNO := false
	for _, v := range existGridNOs {
		if v == g.BasisGridNO {
			isExistBasisGridNO = true
		}
	}

	// 获取实际的缺失的网格线
	realMissGridNOMap := map[int]bool{}
	for _, gridNO := range missGridNOs {
		if gridNO == g.BasisGridNO {
			continue
		}

		// 如果basisGridNO存在了，忽略basisGridNO+1
		if isExistBasisGridNO && gridNO == g.BasisGridNO+1 {
			continue
		}

		realMissGridNOMap[gridNO] = true
	}

	logger.Info("missGrid", logger.Int("BasisGridNO", g.BasisGridNO), logger.Any("missGridNOs", realMissGridNOMap))
	if len(existGridNOs) >= len(g.Grids)-1 {
		logger.Info("no need add new grid record")
		return nil
	}

	for gridNO := range realMissGridNOMap {
		grid := g.Grids[gridNO]
		newRecord, err := g.placeLimitOrder(grid.GID, latestPrice)
		if err != nil {
			logger.Error("g.placeLimitOrder error", logger.String("gsid", g.Gsid), logger.String("symbol", g.Symbol), logger.Any("grid", grid))
			continue
		}

		err = g.updateGridAndAddRecord(grid.GID, newRecord)
		if err != nil {
			logger.Error("updateGridAndAddRecord error", logger.Err(err), logger.Any("newRecord", newRecord))
			continue
		}
	}

	return nil
}

// 下委托订单
func (g *GridProcess) placeLimitOrder(gridNO int, latestPrice float64) (*TradeRecord, error) {
	grid := g.Grids[gridNO]
	// 根据当前价格判断网格为委托买单还是卖单
	side, price, amount, clientOrderID := "", "", "", ""
	quantity := 0.0
	if grid.Price < latestPrice { // 小于最新价格，为买单
		// 委托买单
		side = "buy"
		//price = fmt.Sprintf("%v", grid.Price)
		price = strconv.FormatFloat(grid.Price, 'f', -1, 64)
		amount = fmt.Sprintf("%v", grid.BuyQuantity)
		quantity = grid.BuyQuantity
		clientOrderID = NewGridClientOrderID(g.Exchange, side, g.Gsid, gridNO)
	} else { // 大于等于最新价格，为卖单
		// 委托卖单
		side = "sell"
		//price = fmt.Sprintf("%v", grid.Price)
		price = strconv.FormatFloat(grid.Price, 'f', -1, 64)
		amount = fmt.Sprintf("%v", grid.SellQuantity)
		quantity = grid.SellQuantity
		clientOrderID = NewGridClientOrderID(g.Exchange, side, g.Gsid, gridNO)
	}
	orderID, err := g.ExchangeAccount.PlaceLimitOrder(side, g.Symbol, price, amount, clientOrderID)
	if err != nil {
		logger.Error("placeLimitOrder error", logger.Err(err), logger.String("params", fmt.Sprintf("%s, %s, %s, %s, %s", side, g.Symbol, price, amount, clientOrderID)))
		return nil, err
	}

	grid.Side = side
	grid.OrderID = orderID
	g.Grids[gridNO] = grid
	g.BasisGridNO = gridNO

	newRecord := &TradeRecord{
		GID:           grid.GID,
		OrderID:       orderID,
		ClientOrderID: clientOrderID,
		OrderType:     "limit",
		Side:          side,
		Price:         grid.Price,
		Quantity:      quantity,
		Volume:        grid.Price * quantity,
		Unit:          g.AnchorSymbol,
		OrderState:    huobi.OrderStateSubmitted,
	}

	return newRecord, nil
}

// 更新旧的委托记录
func (g *GridProcess) updateOldOrderRecord(gridNO int, orderID string, orderStatus string, tradeTime int64, tradeAmount float64) error {
	buyPrice := 0.0
	fee := GetExchangeFees(g.Exchange)

	// 更新订单状态和时间
	query := bson.M{"gsid": bson.ObjectIdHex(g.Gsid), "orderID": orderID}
	gtr, err := FindGridTradeRecord(query, bson.M{})
	if err != nil {
		logger.Error("FindGridTradeRecord error", logger.Err(err), logger.Any("query", query))
		return err
	}

	if gtr.Volume != tradeAmount {
		tradeAmount = gtr.Volume
	}

	stateTime := time.Unix(tradeTime/1000, tradeTime%1000)
	fees := tradeAmount * fee
	update := bson.M{
		"$set": bson.M{
			"orderState": orderStatus,
			"stateTime":  stateTime,
			"fees":       fees,
		},
	}
	err = UpdateGridTradeRecord(query, update)
	if err != nil {
		logger.Error("UpdateGridTradeRecord error", logger.Err(err), logger.Any("query", query), logger.Any("update", update))
		return err
	}
	logger.Info("update gridTradeRecord success", logger.String("orderID", orderID), logger.Any("update", update), logger.Float64("fee", fee))

	gtr.OrderState = orderStatus
	gtr.StateTime = stateTime
	gtr.Fees = fees

	// 如果成交的是卖单，需要通知统计
	if gtr.Side == "sell" {
		buyOrderID := ""
		// 判断是否启动网格时委托的卖单
		if gtr.IsStartUpOrder {
			// 查找市价单价格
			query = bson.M{"gsid": gtr.GSID, "side": "buy", "orderType": "market"}
			mgtr, err := FindGridTradeRecord(query, bson.M{})
			if err != nil {
				logger.Error("FindGridTradeRecord error", logger.Err(err), logger.Any("query", query))
			} else {
				buyPrice = mgtr.Price
				buyOrderID = mgtr.OrderID
			}
		} else {
			// 下一格的买入价格
			if gridNO < len(g.Grids)-1 {
				buyPrice = g.Grids[gridNO+1].Price
			}
			// todo 获取最近下一格买入委托订单id
		}

		g.notifyStatistics(gtr, buyPrice)

		query = bson.M{"gsid": bson.ObjectIdHex(g.Gsid), "orderID": orderID}
		update = bson.M{
			"$set": bson.M{
				"buyPrice":   buyPrice,
				"buyOrderID": buyOrderID,
			},
		}
		UpdateGridTradeRecord(query, update)
	}

	return nil
}

// 更新本地的网格和委托记录
func (g *GridProcess) updateGridAndAddRecord(gridNO int, newRecord *TradeRecord) error {
	// 更新网格参数
	query := bson.M{"gsid": bson.ObjectIdHex(g.Gsid)}
	update := bson.M{
		"$set": bson.M{
			"basisGridNO":                           g.BasisGridNO,
			fmt.Sprintf("grids.%d.side", gridNO):    g.Grids[gridNO].Side,
			fmt.Sprintf("grids.%d.orderId", gridNO): g.Grids[gridNO].OrderID,
		},
	}
	err := UpdateGridPendingOrder(query, update)
	if err != nil {
		logger.Error("UpdateGridPendingOrder error", logger.Err(err), logger.Any("query", query), logger.Any("update", update))
		return err
	}
	logger.Info("update gridPendingOrder success", logger.String("gsid", g.Gsid), logger.Int("gridNO", gridNO))

	// 添加新订单记录
	if newRecord != nil {
		gridTradeRecord := &GridTradeRecord{
			GSID: bson.ObjectIdHex(g.Gsid),
			GID:  newRecord.GID,

			OrderID:       newRecord.OrderID,
			ClientOrderID: newRecord.ClientOrderID,
			OrderType:     newRecord.OrderType,
			Side:          newRecord.Side,
			Price:         newRecord.Price,
			Quantity:      newRecord.Quantity,
			Volume:        newRecord.Volume,
			Unit:          newRecord.Unit,

			OrderState: newRecord.OrderState,
			StateTime:  time.Now().Local(),

			Exchange: g.Exchange,
			Symbol:   g.Symbol,
		}
		err = gridTradeRecord.Insert()
		if err != nil {
			logger.Error("UpdateGridPendingOrder error", logger.Err(err), logger.Any("query", query), logger.Any("update", update))
			return err
		}

		logger.Info("add new GridTradeRecord success", logger.String("orderID", gridTradeRecord.OrderID), logger.String("clientOrderID", newRecord.ClientOrderID))
	}

	return nil
}

// ListenOrder 检查订单
func (g *GridProcess) ListenOrder() {
	if g.Ctx == nil {
		g.Ctx, g.Cancel = context.WithCancel(context.Background())
	}

	logger.Info("start to listen grid order", logger.String("gsid", g.Gsid))
	ticker := time.NewTicker(2*time.Hour + time.Duration(krand.Int(1000))*time.Second)
	//ticker := time.NewTicker(time.Duration(100+krand.Int(100)) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-g.Ctx.Done():
			logger.Info(" stopped the listen order goroutine success", logger.String("gsid", g.Gsid))
			return

		case <-ticker.C:
			err := g.syncGridStrategy()
			if err != nil {
				logger.Error("syncGridStrategy error", logger.Err(err), logger.String("gsid", g.Gsid))
			}
		}
	}
}

func str2Float64(str string) float64 {
	if str == "" {
		return 0
	}

	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}

	return f
}

func updateOrderState(orderID string, state string) {
	query := bson.M{"orderID": orderID}
	update := bson.M{"$set": bson.M{"orderState": state}}
	err := UpdateGridTradeRecord(query, update)
	if err != nil {
		logger.Warn("UpdateGridTradeRecord error", logger.Err(err), logger.String("orderID", orderID))
	}
}
