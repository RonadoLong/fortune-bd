package service

import (
	"context"
	"encoding/json"
	"fortune-bd/api/response"
	"fortune-bd/app/exchange-svc/internal/biz"
	"fortune-bd/app/exchange-svc/utils"
	"fortune-bd/libs/helper"
	"fortune-bd/libs/logger"
	"github.com/shopspring/decimal"

	pb "fortune-bd/api/exchange/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	ErrID = "exchangeOrder"
)

type ExOrderService struct {
	pb.UnimplementedExOrderServer
	ExOrderSrv *biz.ExOrderRepo
}

func NewExOrderService() *ExOrderService {
	return &ExOrderService{
		ExOrderSrv: biz.NewExOrderRepo(),
	}
}

func (s *ExOrderService) ExChangeInfo(ctx context.Context, req *emptypb.Empty) (*pb.ExChangeList, error) {
	resp := &pb.ExChangeList{}
	exchangeList, err := s.ExOrderSrv.GetExchangeInfo()
	if err != nil {
		return nil, err
	}
	byteData, err := json.Marshal(exchangeList)
	if err != nil {
		logger.Errorf("ExChangeInfo json.Marshal err %v ", err)
		return nil, response.NewInternalServerErrMsg(ErrID)
	}
	resp.Exchanges = byteData
	return resp, nil
}

func (s *ExOrderService) AddExchangeApi(ctx context.Context, req *pb.ExchangeApi) (*emptypb.Empty, error) {
	err := s.ExOrderSrv.AddExchangeApi(req.UserId, req.ApiKey, req.Secret, req.Passphrase, req.ExchangeId)
	return &emptypb.Empty{}, err
}

func (s *ExOrderService) GetExchangeApiList(ctx context.Context, req *pb.GetExApiReq) (*pb.ExApiResp, error) {
	var resp = &pb.ExApiResp{}
	ret := s.ExOrderSrv.GetExchangeAccountListFromCache(req.UserId)
	if ret != nil && len(ret) > 0 {
		resp.ExApiList = ret
		return resp, nil
	}
	apiResp, err := s.ExOrderSrv.GetExchangeApiList(req.UserId)
	if err != nil {
		return nil, err
	}
	resData, err := json.Marshal(apiResp)
	if err != nil {
		logger.Errorf("GetExchangeApiList json.Marshal err %v ", err)
		return nil, response.NewInternalServerErrMsg(ErrID)
	}
	resp.ExApiList = resData
	s.ExOrderSrv.SetExchangeAccountListCache(req.UserId, resData)
	return resp, nil
}

func (s *ExOrderService) UpdateExchangeApi(ctx context.Context, req *pb.UpdateExchangeApiReq) (*emptypb.Empty, error) {
	err := s.ExOrderSrv.UpdateExchangeApi(req.UserId, req.ApiKey, req.Secret, req.Passphrase, req.ExchangeId, req.ApiId)
	return &emptypb.Empty{}, err
}

func (s *ExOrderService) DeleteExchangeApi(ctx context.Context, req *pb.UserApiReq) (*emptypb.Empty, error) {
	err := s.ExOrderSrv.DeleteExchangeApi(req.UserId, req.ApiId)
	return &emptypb.Empty{}, err
}

func (s *ExOrderService) GetTradeList(ctx context.Context, req *pb.GetTradeListReq) (*pb.TradeListResp, error) {
	var resp = &pb.TradeListResp{}
	tradeCount, err := s.ExOrderSrv.GetTradeCount(req.UserId, req.StrategyId)
	if err != nil {
		return nil, response.NewInternalServerErrMsg(ErrID)
	}
	if tradeCount == 0 {
		return nil, response.NewTradeNotFoundErrMsg(ErrID)
	}
	tradeList := s.ExOrderSrv.GetTradeList(req.UserId, req.StrategyId, req.PageNum, req.PageSize)
	resp.TradeCount = tradeCount
	byteData, _ := json.Marshal(tradeList)
	resp.TradeList = byteData
	return resp, nil
}

func (s *ExOrderService) GetProfitRealTime(ctx context.Context, req *pb.ProfitReq) (*pb.ProfitRealTimeResp, error) {
	var resp = &pb.ProfitRealTimeResp{}
	wqProfit := s.ExOrderSrv.GetProfitRealTime(req.UserId, req.StrategyId)
	if len(wqProfit) == 0 {
		return nil, response.NewProfitNotFoundErrMsg(ErrID)
	}
	byteData, _ := json.Marshal(wqProfit)
	resp.ProfitList = byteData
	return resp, nil
}

func (s *ExOrderService) GetSymbolRankWithRateYear(ctx context.Context, req *emptypb.Empty) (*pb.SymbolRankWithRateYearResp, error) {
	var resp = &pb.SymbolRankWithRateYearResp{}
	data := s.ExOrderSrv.GetSymbolRankWithRateYear()
	if len(data) == 0 {
		return nil, response.NewDataNotFound(ErrID, "暂无数据")
	}
	for _, v := range data {
		resp.Data = append(resp.Data, &pb.SymbolWithRate{
			Symbol:   v.Symbol,
			RateYear: v.RateReturnYear + "%",
			Url:      v.Url,
		})
	}
	return resp, nil
}

func (s *ExOrderService) Evaluation(ctx context.Context, req *pb.TradeReq) (*emptypb.Empty, error) {
	err := s.ExOrderSrv.EvaluationSwap(req)
	return &emptypb.Empty{}, err
}

func (s *ExOrderService) EvaluationSpot(ctx context.Context, req *pb.OrderReq) (*emptypb.Empty, error) {
	err := s.ExOrderSrv.EvaluationSpot(req)
	return &emptypb.Empty{}, err
}

func (s *ExOrderService) StrategyProfitCompensate(ctx context.Context, req *pb.StrategyProfitCompensateReq) (*emptypb.Empty, error) {
	err := s.ExOrderSrv.StrategyProfitCompensate(req.StrategyId, req.Price)
	return &emptypb.Empty{}, err
}

func (s *ExOrderService) GetUserStrategyEva(ctx context.Context, req *pb.UserStrategyDetailReq) (*pb.UserStrategyEvaResp, error) {
	var resp = &pb.UserStrategyEvaResp{}
	profit, err := s.ExOrderSrv.GetProfitByStrID(req.UserId, req.StrategyId)
	if err != nil {
		resp.RealizeProfit = "0"
		resp.RateReturnYear = "0"
		resp.RateReturn = "0"
		return nil, err
	}
	resp.RateReturn = helper.Float64ToString(profit.RateReturn)
	resp.RealizeProfit = profit.RealizeProfit
	resp.RateReturnYear = helper.Float64ToString(profit.RateReturnYear)
	profitDailyList := s.ExOrderSrv.GetProfitDailyByStrID(req.UserId, req.StrategyId, 365)

	dateMap := make(map[string]bool, 0)
	for _, daily := range profitDailyList {
		dateStr := daily.Date.Format("2006-01-02")
		if _, ok := dateMap[dateStr]; ok { //日期去重数据
			continue
		}
		resp.EvaDailyList = append(resp.EvaDailyList, &pb.EvaDaily{
			Date:        dateStr,
			ProfitDaily: daily.RealizeProfit,
		})
		dateMap[dateStr] = true
	}
	return resp, nil
}

func (s *ExOrderService) GetExchangePos(ctx context.Context, req *pb.GetExchangePosReq) (*pb.ExchangePosResp, error) {
	//获取持有币种的账户资产 先支持okex
	pos, err := s.ExOrderSrv.GetExchangePos(req.UserId, req.Exchange)
	if err != nil {
		return nil, err
	}
	if len(pos) == 0 {
		logger.Infof("no data in exchangePos user_id %s", req.UserId)
		return nil, response.NewExchangePosErrMsg(ErrID)
	}
	return &pb.ExchangePosResp{ExchangePos: pos}, nil
}

func (s *ExOrderService) GetTradeSymbols(ctx context.Context, req *pb.TradeSymbolReq) (*pb.GetSymbolsResp, error) {
	symbols, err := s.ExOrderSrv.GetTradeSymbols(req.Exchange, req.Coin)
	if err != nil {
		return nil, err
	}
	return &pb.GetSymbolsResp{Symbols: symbols}, nil
}

func (s *ExOrderService) GetApiKeyInfo(ctx context.Context, req *pb.UserApiKeyReq) (*pb.ExchangeApiResp, error) {
	info, err := s.ExOrderSrv.GetApiKeyInfo(req.UserId, req.ApiKey)
	if err != nil {
		return nil, response.NewDataNotFound(ErrID, "没有找到apiKey")
	}
	//secret, _ := hex.DecodeString(info.Secret)
	//secretKey, _ := encoding.AesDecrypt(secret)
	var resp = &pb.ExchangeApiResp{}
	resp.UserId = info.UserID
	resp.ExchangeId = info.ExchangeID
	resp.ExchangeName = info.ExchangeName
	resp.ApiKey = info.ApiKey
	resp.Passphrase = info.Passphrase
	resp.Secret = info.Secret
	return &pb.ExchangeApiResp{}, nil
}

func (s *ExOrderService) GetAssetsByAllApiKey(ctx context.Context, req *pb.GetExApiReq) (*pb.AssertsResp, error) {
	var resp = &pb.AssertsResp{}
	resp.Asserts = "0(USDT)"
	resp.Profit = "0"
	resp.ProfitPercent = "0%"
	var pos []*response.ExchangeApiResp
	var err error
	ret := s.ExOrderSrv.GetExchangeAccountListFromCache(req.UserId)
	if ret != nil {
		if err = json.Unmarshal(ret, &pos); err != nil {
			logger.Warnf("GetAssetsByAllApiKey jsonUnMarshal pos err %v ", err)
		}
	} else {
		pos, err = s.ExOrderSrv.GetExchangeApiList(req.UserId)
	}
	if len(pos) == 0 || err != nil {
		return nil, err
	}
	asserts := 0.0
	for _, p := range pos {
		asserts += helper.StringToFloat64(p.TotalUsdt)
	}
	asserts = utils.Keep2Decimal(asserts)
	resp.Asserts = helper.Float64ToString(asserts) + "(USDT)"
	strategyList := s.ExOrderSrv.GetUserStrategyByUID(req.UserId)
	if len(strategyList) == 0 {
		return nil, err
	}
	profit := 0.0
	for _, strategy := range strategyList {
		data, err := s.ExOrderSrv.GetProfitByStrID("", strategy.ID)
		if err != nil {
			continue
		}
		if data.Unit == "btc" { //如果是btc  换算成usdt
			profitDecimal, _ := decimal.NewFromString(data.RealizeProfit)
			price := s.ExOrderSrv.GetBtcTickPrice()
			newProfit := profitDecimal.Mul(decimal.NewFromFloat(price))
			logger.Infof("换算 btc数量%s 行情价格 %.2f", data.RealizeProfit, price)
			data.RealizeProfit = newProfit.RoundBank(8).String()
			logger.Infof("换算后 %s", data.RealizeProfit)
		}
		profit += helper.StringToFloat64(data.RealizeProfit)
	}
	profit = utils.Keep2Decimal(profit)
	resp.Profit = helper.Float64ToString(profit)
	if asserts != 0.0 {
		percent := profit / asserts * 100
		resp.ProfitPercent = helper.Float64ToString(utils.Keep2Decimal(percent)) + "%"
	}
	return &pb.AssertsResp{}, nil
}
