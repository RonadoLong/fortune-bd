package service

import (
	"go.mongodb.org/mongo-driver/bson"
	"wq-fotune-backend/app/exchange-srv/model"
)

//弃用
//func (e *ExOrderService) GetUserStrategyList(userID string) ([]*pb.UserStrategyWithRate, error) {
//	strategyListResp := make([]*pb.UserStrategyWithRate, 0)
//	//获取app端 用户策略列表 顺便获取合约的盈亏
//	strategyList := e.dao.GetUserStrategyList(userID)
//	if len(strategyList) == 0 {
//		return nil, response.NewUserStrategyNotFoundErrMsg(ErrID)
//	}
//
//	for _, strategy := range strategyList {
//		strategyResp := &pb.UserStrategyWithRate{
//			UserId:           strategy.UserID,
//			StrategyId:       strategy.StrategyID,
//			ParentStrategyId: strategy.ParentStrategyID,
//			Platform:         strategy.Platform,
//			ApiKey:           strategy.ApiKey,
//			Balance:          strategy.Balance,
//			State:            strategy.State,
//			RunAt:            strategy.CreatedAt.String(),
//			Symbol:           strategy.Symbol,
//			TotalProfit:      "0",
//			RealizeProfit:    "0",
//			UnRealizeProfit:  "0",
//			RateReturnYear:   "0",
//		}
//		apiInfo, err := e.dao.GetExchangeApiByUidAndApi(strategy.UserID, strategy.ApiKey)
//		if err != nil {
//			logger.Infof("GetUserStrategyList GetExchangeApiByUidAndApi has err %v userID %s apiKey %s", err.Error(), strategy.UserID, strategy.ApiKey)
//			strategyListResp = append(strategyListResp, strategyResp)
//			continue
//		}
//		secret, _ := hex.DecodeString(apiInfo.Secret)
//		secretBytes, _ := encoding.AesDecrypt(secret)
//		okClient := api.InitClient(apiInfo.ApiKey, string(secretBytes), apiInfo.Passphrase)
//
//		currencyInfo := goex.CurrencyPair{}
//		if strings.Contains(strategy.Symbol, "ETH-USD") {
//			currencyInfo = goex.ETH_USD
//		}
//		if strings.Contains(strategy.Symbol, "BTC-USD") {
//			currencyInfo = goex.BTC_USD
//		}
//		if strings.Contains(strategy.Symbol, "ETH-USDT") {
//			currencyInfo = goex.ETH_USDT
//		}
//		if strings.Contains(strategy.Symbol, "BTC-USDT") {
//			currencyInfo = goex.BTC_USDT
//		}
//		accInfo, err := okClient.APIClient.OKExSwap.GetFutureAccountInfo(currencyInfo)
//		if err != nil {
//			logger.Warnf("GetUserStrategyList GetFutureAccountInfo has err %v", err)
//			strategyListResp = append(strategyListResp, strategyResp)
//			continue
//		}
//		logger.Infof("%v", accInfo.Info)
//		realizePnl := decimal.NewFromFloat(accInfo.Info.RealizedPnl)
//		unRealizePnl := decimal.NewFromFloat(accInfo.Info.UnrealizedPnl)
//		strategyResp.RealizeProfit = realizePnl.String()     //已实现收益
//		strategyResp.UnRealizeProfit = unRealizePnl.String() //未实现收益
//
//		totalProfit := realizePnl.Add(unRealizePnl)
//		strategyResp.TotalProfit = totalProfit.String()
//		balance := globalF.StringToFloat64(strategy.Balance)
//		totalProfitFloat, _ := totalProfit.Float64()
//
//		if totalProfitFloat == 0.0 || balance == 0.0 {
//			logger.Infof("totalProfitFloat == 0 apikey %s symbol %s", strategy.ApiKey, strategy.Symbol)
//			strategyListResp = append(strategyListResp, strategyResp)
//			continue
//		}
//		log.Println("totalProfitFloat", totalProfitFloat, balance)
//		rateReturn := utils.Keep2Decimal(totalProfitFloat / balance)
//		runDay := int(time.Now().Sub(strategy.CreatedAt).Hours() / 24)
//		if runDay == 0 {
//			runDay += 1
//		}
//		rateReturnYear := rateReturn / (float64(runDay) / 360) * 100 //年化率
//		log.Println(rateReturnYear, rateReturn, "yy")
//		strategyResp.RateReturnYear = decimal.NewFromFloat(utils.Keep2Decimal(rateReturnYear)).String()
//		strategyListResp = append(strategyListResp, strategyResp)
//	}
//	return strategyListResp, nil
//}

//弃用
//func (e *ExOrderService) SetUserStrategyApi(userId, strategyID, apiKey string) error {
//	strategyOld, _ := e.dao.GetUserStrategy(userId, strategyID)
//	if strategyOld != nil {
//		if strategyOld.ApiKey == apiKey {
//			return response.NewSetApiKeySameErrMsg(ErrID)
//		}
//	}
//	if err := e.dao.SetUserStrategyApi(userId, strategyID, apiKey); err != nil {
//		return response.NewSetUserStrategyApiKeyErrMsg(ErrID)
//	}
//	return nil
//}

//弃用
//func (e *ExOrderService) SetUserStrategyBalance(userId, strategyID, balance string) error {
//	strategyOld, _ := e.dao.GetUserStrategy(userId, strategyID)
//	if strategyOld != nil {
//		if strategyOld.Balance == balance {
//			return response.NewSetBalanceSameErrMsg(ErrID)
//		}
//	}
//	if err := e.dao.SetUserStrategyBalance(userId, strategyID, balance); err != nil {
//		return response.NewSetBalanceErrMsg(ErrID)
//	}
//	return nil
//}

//弃用
//func (e *ExOrderService) GetUserStrategyDetail(userId, strategyID string) (*model.WqUserStrategy, error) {
//	strategy, err := e.dao.GetUserStrategy(userId, strategyID)
//	if err != nil {
//		return nil, response.NewUserStrategyDetailErrMsg(ErrID)
//	}
//	return strategy, nil
//}

//// 弃用
//func (e *ExOrderService) CreateUserStrategy(userId string, parentStrategyId int64, balance float32) error {
//	strategy, err := e.dao.GetStrategy(parentStrategyId)
//	if err != nil {
//		return response.NewGetStrategyNotFoundErrMsg(ErrID)
//	}
//	OldUserStrategy, _ := e.dao.GetUserStrategyByParentID(userId, strategy.ID)
//	if OldUserStrategy != nil {
//		if OldUserStrategy.State == 1 {
//			return response.NewCreateStrategyErrMsg(ErrID, "不得重复创建同个机器人")
//		}
//	}
//	apiInfo, _ := e.dao.GetExchangeApiByUidAndExID(userId, strategy.ExchangeID)
//	if apiInfo == nil {
//		return response.NewCreateStrategyErrMsg(ErrID, "请根据要求绑定交易所")
//	}
//
//	if balance <= 0 {
//		return response.NewCreateStrategyErrMsg(ErrID, "投资金额不能为0")
//	}
//	startTime := globalF.GetCurrentTime()
//	userStrategy := &model.WqUserStrategy{
//		UserID: userId,
//		//GroupID:          strategy.GroupID,
//		StrategyID:       snowflake.SNode.Generate().String(),
//		ParentStrategyID: parentStrategyId,
//		ApiKey:           apiInfo.ApiKey,
//		Platform:         strategy.ExchangeName,
//		Balance:          globalF.Float32ToString(balance),
//		State:            1,
//		//Symbol:           strategy.Symbol,
//		CreatedAt: startTime,
//		UpdatedAt: startTime,
//	}
//	if err := e.dao.CreateUserStrategy(userStrategy); err != nil {
//		return response.NewCreateUserStrategyErrMsg(ErrID)
//	}
//	return nil
//}

//弃用
//func (e *ExOrderService) RunUserStrategy(userID, strategyID string) error {
//	// 这个接口直接暂停了策略 没有目前启动功能 创建策略时就已经启动
//	strategy, err := e.dao.GetUserStrategy(userID, strategyID)
//	if err != nil {
//		return response.NewGetStrategyNotFoundErrMsg(ErrID)
//	}
//	balance, _ := decimal.NewFromString(strategy.Balance)
//	if balance.Equal(decimal.NewFromFloat(0)) {
//		return response.NewUserStrategyBalanceErrMsg(ErrID)
//	}
//
//	strategy = &model.WqUserStrategy{}
//	strategy.State = 2
//	strategy.StrategyID = strategyID
//	if err := e.dao.UpdateUserStrategy(userID, strategy); err != nil {
//		return response.NewUserStrategyRunErrMsg(ErrID)
//	}
//	return nil
//}

func (e *ExOrderService) GetUserStrategyOfRun() []*model.GridStrategy {
	return e.dao.GetUserStrategyOfRun(nil)
}

func (e *ExOrderService) GetUserStrategyByUID(userID string) []*model.GridStrategy {
	sql := bson.M{"uid": userID}
	return e.dao.GetUserStrategyOfRun(sql)
}

func (e *ExOrderService) UpdateProfit(strategyID string, profit *model.WqProfit) error {
	return e.dao.UpdateProfit("", strategyID, profit)
}
