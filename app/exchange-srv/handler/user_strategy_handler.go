package handler

//func (e *ExOrderHandler) GetUserStrategyList(ctx context.Context, req *pb.GetStrategyReq, resp *pb.UserStrategyListResp) error {
//	list, err := e.exOrderSrv.GetUserStrategyList(req.UserId)
//	if err != nil {
//		return err
//	}
//	resp.UserStrategyList = list
//	//byteDate, _ := json.Marshal(strategyList)
//	//resp.StrategyList = byteDate
//	return nil
//}

//func (e *ExOrderHandler) SetUserStrategyApi(ctx context.Context, req *pb.SetUserStrategyApiReq, resp *empty.Empty) error {
//	return e.exOrderSrv.SetUserStrategyApi(req.UserId, req.StrategyId, req.ApiKey)
//}
//
//func (e *ExOrderHandler) SetUserStrategyBalance(ctx context.Context, req *pb.SetUserStrategyBalanceReq, resp *empty.Empty) error {
//	balance := fmt.Sprintf("%v", req.Balance)
//	return e.exOrderSrv.SetUserStrategyBalance(req.UserId, req.StrategyId, balance)
//}
//
//func (e *ExOrderHandler) GetUserStrategyDetail(ctx context.Context, req *pb.UserStrategyDetailReq, resp *pb.UserStrategy) error {
//	strategy, err := e.exOrderSrv.GetUserStrategyDetail(req.GetUserId(), req.GetStrategyId())
//	if err != nil {
//		return err
//	}
//	resp.UserId = strategy.UserID
//	resp.ApiKey = strategy.ApiKey
//	resp.Balance = strategy.Balance
//	resp.Platform = strategy.Platform
//	resp.StrategyId = strategy.StrategyID
//	resp.ParentStrategyId = strategy.ParentStrategyID
//	return nil
//}

//func (e *ExOrderHandler) CreateStrategy(ctx context.Context, req *pb.CreateStrategyReq, resp *empty.Empty) error {
//	return e.exOrderSrv.CreateUserStrategy(req.UserId, req.Id, req.Balance)

//secret, _ := hex.DecodeString(apiInfo.Secret)
//secretBytes, _ := encoding.AesDecrypt(secret)
//client := api.InitClient(apiInfo.ApiKey, string(secretBytes), apiInfo.Passphrase)
//accountInfo := client.GetAccount()
//if accountInfo == nil || len(accountInfo.FutureSubAccounts) == 0 {
//    return response.NewCreateStrategyErrMsg(ErrID, "获取账户资金失败")
//}
//for _, account := range accountInfo.FutureSubAccounts {
//    if account.Symbol == strategy.Symbol {
//        if account.TotalAvailBalance <= 0 {
//            return response.NewCreateStrategyErrMsg(ErrID, "账户资金为0，请充值后再运行")
//        }
//        balance := global.Float32ToString(req.Balance)
//        accountBalanceStr := global.Float64ToString(account.TotalAvailBalance)
//        accountBalance := global.StringToFloat32(accountBalanceStr)
//        logger.Infof("获取账户资金为: %v : %v", apiInfo.ApiKey, accountBalance)
//        if accountBalance < req.Balance {
//            balance = accountBalanceStr
//        }
//        startTime := global.GetCurrentTime()
//        userStrategy := &model.WqUserStrategy{
//            UserID:           req.GetUserId(),
//            StrategyID:       snowflake.SNode.Generate().String(),
//            ParentStrategyID: req.GetId(),
//            ApiKey:           apiInfo.ApiKey,
//            Platform:         strategy.ExchangeName,
//            Balance:          balance,
//            State:            1,
//            Symbol:           strategy.Symbol,
//            CreatedAt:        startTime,
//            UpdatedAt:        startTime,
//        }
//        if err := e.dao.CreateUserStrategy(userStrategy); err != nil {
//            return response.NewCreateUserStrategyErrMsg(ErrID)
//        }
//        return nil
//    }
//}
//return response.NewCreateStrategyErrMsg(ErrID, "获取账户资金失败")
//}

//func (e *ExOrderHandler) RunUserStrategy(ctx context.Context, req *pb.UserStrategyRunReq, resp *empty.Empty) error {
//	return e.exOrderSrv.RunUserStrategy(req.UserId, req.StrategyId)
//strategy, err := e.dao.GetUserStrategy(req.UserId, req.StrategyId)
//if err != nil {
//	return response.NewGetStrategyNotFoundErrMsg(ErrID)
//}

//apiInfo, err := e.dao.GetExchangeApiByUidAndApi(req.UserId, strategy.ApiKey)
//if err != nil {
//	return  response.NewExchangeApiExpireErrMsg(ErrID)
//}
//secretByte, _ := encoding.AesEncrypt([]byte(apiInfo.Secret))
//secret := hex.EncodeToString(secretByte)
//client := api.InitClient(strategy.ApiKey, secret, apiInfo.Passphrase)
//account, err := client.APIClient.OKExSwap.GetFutureUserinfo()
//if err != nil {
//	if strings.Contains(err.Error(), "30006") {
//		return errors.New("密钥失效或者不存在")
//	}
//	return response.NewInternalServerErrMsg(ErrID)
//}
//balance, _ := decimal.NewFromString(strategy.Balance)
//if balance.Equal(decimal.NewFromFloat(0)) {
//	return response.NewUserStrategyBalanceErrMsg(ErrID)
//}
////if account.FutureSubAccounts
////TODO check  account money
////var status int32 = 1
////if strategy.State == 1 {
////	status = 2
////}
//strategy = &model.WqUserStrategy{}
//strategy.State = 2
//strategy.StrategyID = req.StrategyId
//if err := e.dao.UpdateUserStrategy(req.UserId, strategy); err != nil {
//	return response.NewUserStrategyRunErrMsg(ErrID)
//}
//return nil
//}
