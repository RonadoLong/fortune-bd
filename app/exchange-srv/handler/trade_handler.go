package handler

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/empty"
	"wq-fotune-backend/pkg/response"
	pb "wq-fotune-backend/app/exchange-srv/proto"
)

func (e *ExOrderHandler) GetTradeSymbols(ctx context.Context, req *pb.TradeSymbolReq, resp *pb.GetSymbolsResp) error {
	symbols, err := e.exOrderSrv.GetTradeSymbols(req.Exchange, req.Coin)
	if err != nil {
		return err
	}
	resp.Symbols = symbols
	return nil
}

func (e *ExOrderHandler) GetTradeList(ctx context.Context, req *pb.GetTradeListReq, resp *pb.TradeListResp) error {
	tradeCount, err := e.exOrderSrv.GetTradeCount(req.UserId, req.StrategyId)
	if err != nil {
		return response.NewInternalServerErrMsg(ErrID)
	}
	if tradeCount == 0 {
		return response.NewTradeNotFoundErrMsg(ErrID)
	}
	tradeList := e.exOrderSrv.GetTradeList(req.UserId, req.StrategyId, req.PageNum, req.PageSize)
	resp.TradeCount = tradeCount
	byteData, _ := json.Marshal(tradeList)
	resp.TradeList = byteData
	return nil
}

func (e *ExOrderHandler) GetProfitRealTime(ctx context.Context, req *pb.ProfitReq, resp *pb.ProfitRealTimeResp) error {
	wqProfit := e.exOrderSrv.GetProfitRealTime(req.UserId, req.StrategyId)
	if len(wqProfit) == 0 {
		return response.NewProfitNotFoundErrMsg(ErrID)
	}
	byteData, _ := json.Marshal(wqProfit)
	resp.ProfitList = byteData
	return nil
}

func (e *ExOrderHandler) GetSymbolRankWithRateYear(ctx context.Context, req *empty.Empty, resp *pb.SymbolRankWithRateYearResp) error {
	data := e.exOrderSrv.GetSymbolRankWithRateYear()
	if len(data) == 0 {
		return response.NewDataNotFound(ErrID, "暂无数据")
	}
	for _, v := range data {
		resp.Data = append(resp.Data, &pb.SymbolWithRate{
			Symbol:   v.Symbol,
			RateYear: v.RateReturnYear + "%",
			Url:      v.Url,
		})
	}
	return nil
}
