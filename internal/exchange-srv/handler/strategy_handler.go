package handler

//func (e *ExOrderHandler) GetStrategyList(ctx context.Context, req *pb.StrategyReq, resp *pb.StrategyList) error {
//	strategyList := e.exOrderSrv.GetStrategyList(req.PageNum, req.PageSize)
//	if len(strategyList) == 0 {
//		return response.NewStrategyNotFoundErrMsg(ErrID)
//	}
//	byteData, _ := json.Marshal(strategyList)
//	resp.StrategyList = byteData
//	return nil
//}

//func (e *ExOrderHandler) GetStrategy(ctx context.Context, req *pb.GetStrategyDetail, resp *pb.Strategy) error {
//	strategy, err := e.exOrderSrv.GetStrategy(req.GetId())
//	if err != nil {
//		return err
//	}
//	resp.Id = strategy.ID
//	resp.Tag = strategy.Tag
//	resp.Level = strategy.Level
//	resp.ExchangeName = strategy.ExchangeName
//	resp.ExchangeId = strategy.ExchangeID
//	resp.Name = strategy.Name
//	resp.Remark = strategy.Remark
//	resp.State = strategy.State
//	return nil
//}
