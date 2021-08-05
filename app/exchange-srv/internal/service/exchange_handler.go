package service

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/shopspring/decimal"
	"wq-fotune-backend/api-gateway/protocol"
	"wq-fotune-backend/libs/helper"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/pkg/response"
	"wq-fotune-backend/pkg/utils"
	pb "wq-fotune-backend/api/exchange"
)

const (
	BTC = "btc"
)

func (e *ExOrderService) ExChangeInfo(ctx context.Context, req *empty.Empty, resp *pb.ExChangeList) error {
	exchangeList, err := e.exOrderSrv.GetExchangeInfo()
	if err != nil {
		return err
	}
	byteData, err := json.Marshal(exchangeList)
	if err != nil {
		logger.Errorf("ExChangeInfo json.Marshal err %v ", err)
		return response.NewInternalServerErrMsg(ErrID)
	}
	resp.Exchanges = byteData
	return nil
}

func (e *ExOrderService) AddExchangeApi(ctx context.Context, req *pb.ExchangeApi, resp *empty.Empty) error {
	return e.exOrderSrv.AddExchangeApi(req.UserID, req.ApiKey, req.Secret, req.Passphrase, req.ExchangeID)
}

func (e *ExOrderService) GetExchangeApiList(ctx context.Context, req *pb.GetExApiReq, resp *pb.ExApiResp) error {
	ret := e.exOrderSrv.GetExchangeAccountListFromCache(req.UserId)
	if ret != nil {
		resp.ExApiList = ret
		return nil
	}
	apiResp, err := e.exOrderSrv.GetExchangeApiList(req.UserId)
	if err != nil {
		return err
	}
	resData, err := json.Marshal(apiResp)
	if err != nil {
		logger.Errorf("GetExchangeApiList json.Marshal err %v ", err)
		return response.NewInternalServerErrMsg(ErrID)
	}
	resp.ExApiList = resData
	e.exOrderSrv.SetExchangeAccountListCache(req.UserId, resData)
	return nil
}

func (e *ExOrderService) GetExchangePos(ctx context.Context, req *pb.GetExchangePosReq, resp *pb.ExchangePosResp) error {
	//获取持有币种的账户资产 先支持okex
	pos, err := e.exOrderSrv.GetExchangePos(req.UserId, req.Exchange)
	if err != nil {
		return err
	}
	resp.ExchangePos = pos
	if len(resp.ExchangePos) == 0 {
		logger.Infof("no data in exchangePos userId %s", req.UserId)
		return response.NewExchangePosErrMsg(ErrID)
	}
	return nil
}

func (e *ExOrderService) UpdateExchangeApi(ctx context.Context, req *pb.UpdateExchangeApiReq, resp *empty.Empty) error {
	return e.exOrderSrv.UpdateExchangeApi(req.UserID, req.ApiKey, req.Secret, req.Passphrase, req.ExchangeID, req.ApiID)
}

func (e *ExOrderService) DeleteExchangeApi(ctx context.Context, req *pb.UserApiReq, resp *empty.Empty) error {
	return e.exOrderSrv.DeleteExchangeApi(req.UserId, req.ApiID)
}

func (e *ExOrderService) GetApiKeyInfo(ctx context.Context, req *pb.UserApiKeyReq, resp *pb.ExchangeApiResp) error {
	info, err := e.exOrderSrv.GetApiKeyInfo(req.UserId, req.ApiKey)
	if err != nil {
		return response.NewDataNotFound(ErrID, "没有找到apiKey")
	}
	//secret, _ := hex.DecodeString(info.Secret)
	//secretKey, _ := encoding.AesDecrypt(secret)
	resp.UserId = info.UserID
	resp.ExchangeId = info.ExchangeID
	resp.ExchangeName = info.ExchangeName
	resp.ApiKey = info.ApiKey
	resp.Passphrase = info.Passphrase
	resp.Secret = info.Secret
	return nil
}

func (e *ExOrderService) GetAssetsByAllApiKey(ctx context.Context, req *pb.GetExApiReq, resp *pb.AssertsResp) error {
	resp.Asserts = "0(USDT)"
	resp.Profit = "0"
	resp.ProfitPercent = "0%"
	var pos []*protocol.ExchangeApiResp
	var err error
	ret := e.exOrderSrv.GetExchangeAccountListFromCache(req.UserId)
	if ret != nil {
		if err = json.Unmarshal(ret, &pos); err != nil {
			logger.Warnf("GetAssetsByAllApiKey jsonUnMarshal pos err %v ", err)
		}
	} else {
		pos, err = e.exOrderSrv.GetExchangeApiList(req.UserId)
	}
	if len(pos) == 0 || err != nil {
		return nil
	}
	asserts := 0.0
	for _, p := range pos {
		asserts += helper.StringToFloat64(p.TotalUsdt)
	}
	asserts = utils.Keep2Decimal(asserts)
	resp.Asserts = helper.Float64ToString(asserts) + "(USDT)"
	strategyList := e.exOrderSrv.GetUserStrategyByUID(req.UserId)
	if len(strategyList) == 0 {
		return nil
	}

	profit := 0.0
	for _, strategy := range strategyList {
		data, err := e.exOrderSrv.GetProfitByStrID("", strategy.ID)
		if err != nil {
			continue
		}
		if data.Unit == BTC { //如果是btc  换算成usdt
			profitDecimal, _ := decimal.NewFromString(data.RealizeProfit)
			price := e.exOrderSrv.GetBtcTickPrice()
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
	return nil
}
