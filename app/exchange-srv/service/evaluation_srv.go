package service

import (
	"fmt"
	"github.com/shopspring/decimal"
	"log"
	"strconv"
	"strings"
	"time"
	"wq-fotune-backend/libs/helper"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/pkg/evaluation"
	exchange_info "wq-fotune-backend/pkg/exchange-info"
	"wq-fotune-backend/pkg/goex"
	"wq-fotune-backend/pkg/response"
	"wq-fotune-backend/pkg/utils"
	"wq-fotune-backend/app/exchange-srv/dao"
	"wq-fotune-backend/app/exchange-srv/model"
	pb "wq-fotune-backend/app/exchange-srv/proto"
	"wq-fotune-backend/app/forward-offer-srv/global"
)

//不再使用
func (e *ExOrderService) EvaluationSwap(req *pb.TradeReq) error {
	log.Println("开始计算合约统计")
	//参数处理
	direction := strings.ToLower(req.Direction)
	price, _ := strconv.ParseFloat(req.GetPrice(), 64)
	pos, _ := strconv.ParseFloat(req.Volume, 64)

	log.Println(price, "trade price")
	//合约交易的统计
	strategy, err := e.dao.GetUserStrategy(req.GetUserId(), req.GetStrategyId())
	if err != nil {
		return response.NewEvaluationStrategyErrMsg(ErrID)
	}
	trade := model.NewWqTrade(req.UserId, req.TradeId, req.ApiKey, req.StrategyId, req.Symbol, req.Volume, req.Commission, req.Direction)
	// price

	oldTrade, err := e.dao.GetTradeByLastID(req.UserId, req.StrategyId)
	if err != nil { //if no trade create a new
		if err != dao.RowNotFoundErr {
			return response.NewInternalServerErrMsg(ErrID)
		}
		log.Println("first one ", trade.TradeID)
		trade.OpenPrice = price
		trade.Pos = req.Volume
		trade.AvgPrice = price
		trade.PosDirection = direction
		if err := e.dao.CreateTrade(trade); err != nil {
			return response.NewInternalServerErrMsg(ErrID)
		}
		return nil
	}

	oldPos, _ := strconv.ParseFloat(oldTrade.Pos, 64)
	log.Println(oldPos, "oldPos")
	if oldPos == 0 { // if oldPos == 0 create a new 上条记录持仓为0 那么本次开仓
		trade.OpenPrice = price
		trade.Pos = req.Volume
		trade.AvgPrice = price
		trade.PosDirection = direction
		if err := e.dao.CreateTrade(trade); err != nil {
			return response.NewInternalServerErrMsg(ErrID)
		}
		return nil
	}
	//

	// add buy or add sell
	if (direction == "buy" && direction == oldTrade.PosDirection) || (direction == "sell" && direction == oldTrade.PosDirection) {
		posFloat := pos + oldPos                          //如果两次方向相同 那么仓位相加
		posNew := decimal.NewFromFloat(posFloat).String() //如果两次方向相同 那么仓位相加
		trade.OpenPrice = price                           // 本次交易均价
		trade.Pos = posNew
		trade.PosDirection = direction               // 目前持仓方向
		lastTotalPrice := oldTrade.AvgPrice * oldPos // 9500 * 0.1
		thisTotalPrice := price * pos                //9500 * 0.1

		avgPrice := (lastTotalPrice + thisTotalPrice) / posFloat // 19000 / 0.2
		avgPriceStr := fmt.Sprintf("%.2f", avgPrice)
		trade.AvgPrice, _ = strconv.ParseFloat(avgPriceStr, 64)
		if err := e.dao.CreateTrade(trade); err != nil {
			return response.NewInternalServerErrMsg(ErrID)
		}
		return nil
	}

	commission, _ := strconv.ParseFloat(trade.Commission, 64)
	balance := strategy.TotalSum
	//close  平仓
	var eva evaluation.Evaluation
	okex := evaluation.NewOKex(trade.Symbol, pos, oldTrade.AvgPrice, commission, price, balance)
	eva = okex
	// close sell or buy
	if (direction == "buy" && direction != oldTrade.PosDirection) || (direction == "sell" && direction != oldTrade.PosDirection) {
		okex.Direction = getEvaDirection(direction)
		okex.Type = evaluation.ContractTradingUsdt
		diff := pos - oldPos

		trade.Profit = decimal.NewFromFloat(eva.CalProfit()).String()
		trade.OpenPrice = oldTrade.AvgPrice
		trade.ClosePrice = price
		if diff == 0 {
			trade.Pos = "0"
			trade.PosDirection = ""
		}
		//close sell and continue sell
		if diff < 0 {
			trade.AvgPrice = oldTrade.AvgPrice
			trade.Pos = decimal.NewFromFloat(oldPos - pos).String()
			trade.PosDirection = getPosDirection(direction)
		}
		//close sell and open buy
		if diff > 0 {
			trade.AvgPrice = price
			trade.Pos = decimal.NewFromFloat(diff).String()
			trade.PosDirection = getPosDirection(direction)
		}

	}
	if err := e.dao.CreateTrade(trade); err != nil {
		return response.NewInternalServerErrMsg(ErrID)
	}
	//update wqProfit
	wqProfit, err := e.dao.GetProfitByID(req.UserId, req.StrategyId)
	runDay := int(time.Now().Sub(strategy.CreatedAt).Hours() / 24)
	if runDay == 0 {
		runDay += 1
	}
	if err != nil {
		profit := &model.WqProfit{
			UserID:          req.UserId,
			ApiKey:          strategy.ApiKey,
			StrategyID:      req.StrategyId,
			Symbol:          req.Symbol,
			RealizeProfit:   trade.Profit,
			UnRealizeProfit: "",
			Position:        0,
			RateReturn:      eva.RateReturn(),
			RateReturnYear:  0,
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}
		profit.RateReturnYear = profit.RateReturn / (float64(runDay) / 360)
		profit.RateReturnYear = utils.Keep2Decimal(profit.RateReturnYear)
		if err := e.dao.CreateProfit(profit); err != nil {
			logger.Errorf("CreateProfit has err %v", err)
		}
		return nil
	}
	newProfit, _ := decimal.NewFromString(trade.Profit)
	oldProfit, _ := decimal.NewFromString(wqProfit.RealizeProfit)
	totalProfit := newProfit.Add(oldProfit)
	updateWqProfit := &model.WqProfit{
		RealizeProfit: totalProfit.String(),
	}
	totalProfitFloat, _ := totalProfit.Float64()
	okex = &evaluation.Okex{
		Principal: balance,
		Profit:    totalProfitFloat,
	}
	updateWqProfit.RateReturn = okex.RateReturn()
	updateWqProfit.RateReturnYear = updateWqProfit.RateReturn / (float64(runDay) / 360)
	updateWqProfit.RateReturnYear = utils.Keep2Decimal(updateWqProfit.RateReturnYear)
	if err := e.dao.UpdateProfit(req.UserId, req.StrategyId, updateWqProfit); err != nil {
		logger.Errorf("UpdateProfit has err %v", err)
	}
	return nil
}

func getCommission(buyPrice, volume float64, exchane string) float64 {
	feeRate := 0.0
	switch exchane {
	case exchange_info.HUOBI:
		feeRate = 0.002
	case exchange_info.BINANCE:
		feeRate = 0.001
	case exchange_info.OKEX:
		feeRate = 0.001
	}
	feeDecimail := decimal.NewFromFloat(buyPrice).Mul(decimal.NewFromFloat(volume)).Mul(decimal.NewFromFloat(feeRate)).Round(8)
	f, _ := feeDecimail.Float64()
	return f
}

func (e *ExOrderService) EvaluationSpot(req *pb.OrderReq) error {
	//币币交易实时统计 规则很简单 只有买币卖币
	tradeAt := utils.FormatTimeFromUnix(goex.ToInt64(req.TradeAt))
	logger.Infof("订单进入orderID%v strategy %v 订单成交时间 %v", req.OrderId, req.StrategyId, tradeAt)
	strategy, err := e.dao.GetUserStrategy(req.GetUserId(), req.GetStrategyId())
	if err != nil {
		return response.NewEvaluationStrategyErrMsg(ErrID)
	}
	runDay := int(time.Now().Sub(strategy.CreatedAt).Hours() / 24)
	balance := strategy.TotalSum

	runDay += 1

	volume, _ := strconv.ParseFloat(req.Volume, 64)
	buyPrice, _ := strconv.ParseFloat(req.BuyPrice, 64)
	sellPrice, _ := strconv.ParseFloat(req.SellPrice, 64)
	profit := (sellPrice * volume) - (buyPrice * volume)
	profit = utils.Keep8Decimal(profit)

	commission := 0.0
	order, err := e.dao.GetOrderRecord(req.OrderId)
	if err != nil {
		logger.Warnf("现货统计没有找到订单 %v", err)
	} else {
		commission = utils.Keep8Decimal(order.Fees)
	}
	buyFee := getCommission(buyPrice, volume, req.Exchange)
	commission = commission + buyFee
	commission = utils.Keep8Decimal(commission)
	//利润等于减去手续费
	profit, _ = decimal.NewFromFloat(profit).Sub(decimal.NewFromFloat(commission)).Float64()
	trade := &model.WqTradeRecord{
		ID:            0,
		UserID:        req.UserId,
		ApiKey:        req.ApiKey,
		StrategyID:    req.StrategyId,
		OrderID:       req.OrderId,
		Symbol:        req.Symbol,
		RealizeProfit: helper.Float64ToString(profit),
		BuyPrice:      req.BuyPrice,
		SellPrice:     req.SellPrice,
		Unit:          req.Unit,
		Commission:    global.Float64ToString(commission),
		Volume:        req.Volume,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := e.dao.CreateWqTradeRecord(trade); err != nil {
		logger.Errorf("CreateWqTradeRecord has err %v %+v", err, trade)
	}
	var eva evaluation.Evaluation
	okex := &evaluation.Okex{
		Principal: balance,
		Profit:    utils.Keep2Decimal(profit),
	}
	eva = okex

	wqProfit, err := e.dao.GetProfitByID(req.UserId, req.StrategyId)
	if err != nil { //新建profit
		profit := &model.WqProfit{
			UserID:          req.UserId,
			ApiKey:          strategy.ApiKey,
			StrategyID:      req.StrategyId,
			Symbol:          req.Symbol,
			RealizeProfit:   global.Float64ToString(profit),
			UnRealizeProfit: "",
			Position:        0,
			RateReturn:      eva.RateReturn(),
			RateReturnYear:  0,
			Unit:            req.Unit,
			Commission:      global.Float64ToString(commission),
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}
		profit.RateReturnYear = profit.RateReturn / (float64(runDay) / 365)
		profit.RateReturnYear = utils.Keep2Decimal(profit.RateReturnYear)
		logger.Infof("插入数据 %+v", profit)
		if err := e.dao.CreateProfit(profit); err != nil {
			logger.Errorf("CreateProfit has err %v", err)
		}
		return nil
	}
	//更新profit
	newProfit := decimal.NewFromFloat(profit)
	oldProfit, _ := decimal.NewFromString(wqProfit.RealizeProfit)
	totalProfit := newProfit.Add(oldProfit)
	updateWqProfit := &model.WqProfit{
		RealizeProfit: totalProfit.String(),
	}
	totalProfitFloat, _ := totalProfit.Float64()
	okex = &evaluation.Okex{
		Principal: balance,
		Profit:    totalProfitFloat,
	}

	oldCommission := helper.StringToFloat64(wqProfit.Commission)
	newCommission := decimal.NewFromFloat(oldCommission).Add(decimal.NewFromFloat(commission))
	updateWqProfit.Commission = newCommission.String()
	updateWqProfit.RateReturn = okex.RateReturn()
	updateWqProfit.RateReturnYear = updateWqProfit.RateReturn / (float64(runDay) / 365)
	updateWqProfit.RateReturnYear = utils.Keep2Decimal(updateWqProfit.RateReturnYear)
	updateWqProfit.UpdatedAt = time.Now()
	if err := e.dao.UpdateProfit(req.UserId, req.StrategyId, updateWqProfit); err != nil {
		logger.Errorf("UpdateProfit has err %v", err)
	}
	return nil
}

func (e *ExOrderService) CreateWqProfitDaily(profit *model.WqProfit) error {
	profit.UpdatedAt = time.Now()
	profit.CreatedAt = time.Now()
	profit.ID = 0
	date := time.Now().Add(-1 * time.Hour)
	profitDaily := &model.WqProfitDaily{
		WqProfit: *profit,
		Date:     date,
	}
	return e.dao.CreateProfitDaily(profitDaily)
}

func (e *ExOrderService) StrategyProfitCompensate(strategyId string, price float64) error {
	price = utils.Keep8Decimal(price)
	profitData, err := e.GetProfitByStrID("", strategyId)
	if err != nil {
		return response.NewDataNotFound(ErrID, "没有找到策略统计")
	}
	strategy, err := e.dao.GetUserStrategy("", strategyId)
	if err != nil {
		return response.NewDataNotFound(ErrID, "没有找到策略")
	}
	runDay := int(time.Now().Sub(strategy.CreatedAt).Hours() / 24)
	runDay += 1
	profit, _ := decimal.NewFromString(profitData.RealizeProfit)
	priceDecimal := decimal.NewFromFloat(price)
	newProfit := profit.Add(priceDecimal)
	profitData.RealizeProfit = newProfit.String()

	newProfitFloat, _ := newProfit.Float64()
	okex := &evaluation.Okex{
		Principal: strategy.TotalSum,
		Profit:    newProfitFloat,
	}
	profitData.RateReturn = okex.RateReturn()
	profitData.RateReturnYear = profitData.RateReturn / (float64(runDay) / 365)
	profitData.RateReturnYear = utils.Keep2Decimal(profitData.RateReturnYear)
	profitData.UpdatedAt = time.Now()

	if err := e.dao.UpdateProfit("", strategyId, profitData); err != nil {
		logger.Errorf("UpdateProfit has err %v", err)
	}

	profitDaily, err := e.dao.GetLastProfitByStrategyId(strategyId)
	if err != nil {
		logger.Warnf("GetLastProfitByStrategyId has error %v strategyId %s", err, strategyId)
	}
	profitDaily.RealizeProfit = profitData.RealizeProfit
	profitDaily.RateReturn = profitData.RateReturn
	profitDaily.RateReturnYear = profitData.RateReturnYear
	profitDaily.UpdatedAt = profitData.UpdatedAt
	if err := e.dao.UpdateProfitDaily(profitDaily); err != nil {
		logger.Errorf("UpdateProfitDaily has err %v", err)
	}
	return nil
}
