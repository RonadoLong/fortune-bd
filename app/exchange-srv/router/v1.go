package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/empty"
	jsoniter "github.com/json-iterator/go"
	"github.com/micro/go-micro/v2/errors"
	"strconv"
	pb "wq-fotune-backend/api/exchange"
	"wq-fotune-backend/api/protocol"
	"wq-fotune-backend/app/exchange-srv/client"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/jwt"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/pkg/middleware"
	"wq-fotune-backend/pkg/response"
)

var (
	exOrderService pb.ExOrderService
	orderService   pb.ForwardOfferService
)

func apiV1(group *gin.RouterGroup) {
	exOrderService = client.NewExOrderClient(env.EtcdAddr)
	orderService = client.NewForwardOfferClient(env.EtcdAddr)

	group.GET("/exchange/info", GetExchangeInfo)
	group.GET("/exchange/symbols/:exchange/:coin", GetTradeSymbol)
	group.GET("/exchange/symbolRank", GetSymbolRank)
	group.GET("/exchange/apiInfo/:userId/:apiKey", GetApiKeyInfo)
	group.POST("/forward-offer/order", PutOrder)
	group.POST("/forward-offer/orderGrid", PutOrderGrid)
	group.POST("/forward-offer/profitAdd", AddProfit)
	group.GET("/user/strategy/evaluationNoAuth/:userId/:strategyId", GetUserStrategyEvaNoAuth)

	group.Use(middleware.JWTAuth())
	group.POST("/api/add", AddExchangeAPI)
	group.GET("/exchange/api/list", GetExApiList)
	group.PUT("/exchange/api/update", UpdateExApi)
	group.DELETE("/exchange/api/:apiId", DeleteExApi)
	group.GET("/user/trade/:strategyId/:pageNum/:pageSize", GetTradeList)
	group.GET("/user/strategy/profit/:strategyId", GetProfit)
	group.GET("/user/strategy/evaluation/:strategyId", GetUserStrategyEva)
	group.GET("/user/exchange/pos", GetExchangePos)
	group.GET("/user/exchange/assert", GetUserAssert)

}

func AddProfit(c *gin.Context) {
	var req pb.StrategyProfitCompensateReq
	defer c.Request.Body.Close()
	if err := jsonpb.Unmarshal(c.Request.Body, &req); err != nil {
		response.NewBindJsonErr(c, nil)
		return
	}
	_, err := exOrderService.StrategyProfitCompensate(context.Background(), &req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, nil)
}

func GetSymbolRank(c *gin.Context) {
	resp, err := exOrderService.GetSymbolRankWithRateYear(context.Background(), &empty.Empty{})
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, resp.Data)
}

func GetUserAssert(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	assert, err := exOrderService.GetAssetsByAllApiKey(context.Background(), &pb.GetExApiReq{UserId: jwtP.UserID})
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, assert)
}

func GetTradeSymbol(c *gin.Context) {
	exchange := c.Param("exchange")
	coin := c.Param("coin")
	symbols, err := exOrderService.GetTradeSymbols(context.Background(), &pb.TradeSymbolReq{Exchange: exchange, Coin: coin})
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, symbols.Symbols)
}

func PutOrder(c *gin.Context) {
	var req pb.TradeSignal
	defer c.Request.Body.Close()
	if err := jsonpb.Unmarshal(c.Request.Body, &req); err != nil {
		response.NewBindJsonErr(c, nil)
		return
	}
	_, err := orderService.PushSwapOrder(context.Background(), &req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, nil)
}

func GetExchangeInfo(c *gin.Context) {
	infoList, err := exOrderService.ExChangeInfo(context.Background(), &empty.Empty{})
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	var resp []protocol.ExchangeResp
	if err := jsoniter.Unmarshal(infoList.Exchanges, &resp); err != nil {
		response.NewInternalServerErr(c, nil)
		return
	}
	response.NewSuccess(c, resp)
}

func AddExchangeAPI(c *gin.Context) {
	var req protocol.ExchangeApiReq
	if err := c.BindJSON(&req); err != nil {
		response.NewBindJsonErr(c, nil)
		return
	}
	if err := req.CheckNotNull(); err != nil {
		response.NewErrorParam(c, err.Error(), nil)
		return
	}
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	exAPi := &pb.ExchangeApi{
		UserID:     jwtP.UserID,
		ExchangeID: req.ExchangeID,
		ApiKey:     req.ApiKey,
		Secret:     req.Secret,
		Passphrase: req.Passphrase,
	}
	_, err := exOrderService.AddExchangeApi(context.Background(), exAPi)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, nil)
}

func GetApiKeyInfo(c *gin.Context) {
	userID := c.Param("userId")
	apiKey := c.Param("apiKey")
	req := &pb.UserApiKeyReq{
		UserId: userID,
		ApiKey: apiKey,
	}
	info, err := exOrderService.GetApiKeyInfo(context.Background(), req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	resp := protocol.ExchangeApiInfoResp{
		UserId:       info.UserId,
		ExchangeId:   info.ExchangeId,
		ExchangeName: info.ExchangeName,
		ApiKey:       info.ApiKey,
		Secret:       info.Secret,
		Passphrase:   info.Passphrase,
	}
	response.NewSuccess(c, resp)
}

func GetExApiList(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)

	req := &pb.GetExApiReq{UserId: jwtP.UserID}
	apiList, err := exOrderService.GetExchangeApiList(context.Background(), req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	var apiResp []*protocol.ExchangeApiResp
	if err := jsoniter.Unmarshal(apiList.ExApiList, &apiResp); err != nil {
		logger.Errorf("GetExApiList json Unmarshal err %v", err)
		response.NewInternalServerErr(c, nil)
		return
	}
	response.NewSuccess(c, apiResp)
}

func UpdateExApi(c *gin.Context) {
	var req protocol.UpdateExchangeApiReq
	if err := c.BindJSON(&req); err != nil {
		response.NewBindJsonErr(c, nil)
		return
	}
	if err := req.CheckNotNull(); err != nil {
		response.NewErrorParam(c, err.Error(), nil)
		return
	}
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	exAPi := &pb.UpdateExchangeApiReq{
		ApiID:      req.ApiId,
		UserID:     jwtP.UserID,
		ExchangeID: req.ExchangeID,
		ApiKey:     req.ApiKey,
		Secret:     req.Secret,
		Passphrase: req.Passphrase,
	}
	_, err := exOrderService.UpdateExchangeApi(context.Background(), exAPi)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, nil)
}

func DeleteExApi(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	apiId := c.Param("apiId")
	i, err := strconv.Atoi(apiId)
	if err != nil {
		response.NewInternalServerErr(c, nil)
		return
	}
	req := &pb.UserApiReq{ApiID: int64(i), UserId: jwtP.UserID}
	_, err = exOrderService.DeleteExchangeApi(context.Background(), req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, nil)
}

func GetTradeList(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	strategyId := c.Param("strategyId")
	pageNum, _ := strconv.Atoi(c.Param("pageNum"))
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	req := &pb.GetTradeListReq{
		UserId:     jwtP.UserID,
		StrategyId: strategyId,
		PageNum:    int32(pageNum),
		PageSize:   int32(pageSize),
	}
	resp, err := exOrderService.GetTradeList(context.Background(), req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}

	var tradeResp []*protocol.TradeResp
	if err := jsoniter.Unmarshal(resp.TradeList, &tradeResp); err != nil {
		response.NewInternalServerErr(c, nil)
		return
	}
	response.NewSuccess(c, tradeResp)
}

func GetProfit(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	strategyId := c.Param("strategyId")
	req := &pb.ProfitReq{
		UserId:     jwtP.UserID,
		StrategyId: strategyId,
	}
	resp, err := exOrderService.GetProfitRealTime(context.Background(), req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	var profitResp []*protocol.ProfitResp
	if err := jsoniter.Unmarshal(resp.ProfitList, &profitResp); err != nil {
		response.NewInternalServerErr(c, nil)
		return
	}
	response.NewSuccess(c, profitResp)
}

func GetExchangePos(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	req := &pb.GetExchangePosReq{UserId: jwtP.UserID}
	resp, err := exOrderService.GetExchangePos(context.Background(), req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, resp)
}

func PutOrderGrid(c *gin.Context) {
	var req pb.OrderReq
	defer c.Request.Body.Close()
	if err := jsonpb.Unmarshal(c.Request.Body, &req); err != nil {
		response.NewBindJsonErr(c, nil)
		return
	}
	_, err := exOrderService.EvaluationSpot(context.Background(), &req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, nil)
}

func GetUserStrategyEva(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	strategyId := c.Param("strategyId")
	req := &pb.UserStrategyDetailReq{
		UserId:     jwtP.UserID,
		StrategyId: strategyId,
	}
	eva, err := exOrderService.GetUserStrategyEva(context.Background(), req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	if len(eva.EvaDailyList) == 0 {
		response.NewSuccess(c, &protocol.UserStrategyEvaResp{
			RealizeProfit:  eva.RealizeProfit,
			RateReturnYear: eva.RateReturnYear,
			RateReturn:     eva.RateReturn,
			EvaDaily:       make([]interface{}, 0),
		})
		return
	}

	response.NewSuccess(c, eva)
}

func GetUserStrategyEvaNoAuth(c *gin.Context) {
	userID := c.Param("userId")
	strategyId := c.Param("strategyId")
	req := &pb.UserStrategyDetailReq{
		UserId:     userID,
		StrategyId: strategyId,
	}
	eva, err := exOrderService.GetUserStrategyEva(context.Background(), req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	if len(eva.EvaDailyList) == 0 {
		response.NewSuccess(c, &protocol.UserStrategyEvaResp{
			RealizeProfit:  eva.RealizeProfit,
			RateReturnYear: eva.RateReturnYear,
			RateReturn:     eva.RateReturn,
			EvaDaily:       make([]interface{}, 0),
		})
		return
	}

	response.NewSuccess(c, eva)
}
