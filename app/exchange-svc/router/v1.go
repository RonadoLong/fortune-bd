package router

import (
	"context"
	pb "fortune-bd/api/exchange/v1"
	"fortune-bd/api/response"
	"fortune-bd/app/exchange-svc/internal/service"
	"fortune-bd/libs/jwt"
	"fortune-bd/libs/logger"
	"fortune-bd/libs/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/empty"
	jsoniter "github.com/json-iterator/go"
	"strconv"

)

var s *service.ExOrderService

func getService() *service.ExOrderService{
	if s == nil {
		s = service.NewExOrderService()
	}
	return  s
}

func apiV1(group *gin.RouterGroup) {
	group.GET("/exchange/info", GetExchangeInfoHandler)
	group.GET("/exchange/symbols/:exchange/:coin", GetTradeSymbolHandler)
	group.GET("/exchange/symbolRank", GetSymbolRankHandler)
	group.GET("/exchange/apiInfo/:user_id/:apiKey", GetApiKeyInfoHandler)
	//group.POST("/forward-offer/order", PutOrder)
	group.POST("/forward-offer/orderGrid", PutOrderGridHandler)
	group.POST("/forward-offer/profitAdd", AddProfitHandler)
	group.GET("/user/strategy/evaluationNoAuth/:user_id/:strategyId", GetUserStrategyEvaNoAuthHandler)
	group.Use(middleware.JWTAuth())
	group.POST("/api/add", AddExchangeAPIHandler)
	group.GET("/exchange/api/list", GetExApiList)
	group.PUT("/exchange/api/update", UpdateExApiHandler)
	group.DELETE("/exchange/api/:apiId", DeleteExApiHandler)
	group.GET("/user/trade/:strategyId/:pageNum/:pageSize", GetTradeListHandler)
	group.GET("/user/strategy/profit/:strategyId", GetProfitHandler)
	group.GET("/user/strategy/evaluation/:strategyId", GetUserStrategyEvaHandler)
	group.GET("/user/exchange/pos", GetExchangePosHandler)
	group.GET("/user/exchange/assert", GetUserAssertHandler)

}

func AddProfitHandler(c *gin.Context) {
	var req pb.StrategyProfitCompensateReq
	defer c.Request.Body.Close()
	if err := jsonpb.Unmarshal(c.Request.Body, &req); err != nil {
		response.NewBindJsonErr(c, nil)
		return
	}
	_, err := getService().StrategyProfitCompensate(context.Background(), &req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	response.NewSuccess(c, nil)
}

func GetSymbolRankHandler(c *gin.Context) {
	resp,err := getService().GetSymbolRankWithRateYear(context.Background(), &empty.Empty{})
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	response.NewSuccess(c, resp.Data)
}

func GetUserAssertHandler(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	resp, err := getService().GetAssetsByAllApiKey(context.Background(), &pb.GetExApiReq{UserId: jwtP.UserID})
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	response.NewSuccess(c, resp)
}

func GetTradeSymbolHandler(c *gin.Context) {
	exchange := c.Param("exchange")
	coin := c.Param("coin")
	resp, err := getService().GetTradeSymbols(context.Background(), &pb.TradeSymbolReq{Exchange: exchange, Coin: coin})
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	response.NewSuccess(c, resp.Symbols)
}

//func PutOrder(c *gin.Context) {
//	var req pb.TradeSignal
//	defer c.Request.Body.Close()
//	if err := jsonpb.Unmarshal(c.Request.Body, &req); err != nil {
//		response.NewBindJsonErr(c, nil)
//		return
//	}
//	var orderService = service.NewForwardOfferHandle()
//	err := orderService.PushSwapOrder(context.Background(), &req, nil)
//	if err != nil {
//		fromError := errors.FromError(err)
//		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
//		return
//	}
//	response.NewSuccess(c, nil)
//}

func GetExchangeInfoHandler(c *gin.Context) {
	infoList, err := getService().ExChangeInfo(context.Background(), &empty.Empty{})
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	var resp []response.ExchangeResp
	if err := jsoniter.Unmarshal(infoList.Exchanges, &resp); err != nil {
		response.NewInternalServerErr(c, nil)
		return
	}
	response.NewSuccess(c, resp)
}

func AddExchangeAPIHandler(c *gin.Context) {
	var req response.ExchangeApiReq
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
		UserId:     jwtP.UserID,
		ExchangeId: req.ExchangeID,
		ApiKey:     req.ApiKey,
		Secret:     req.Secret,
		Passphrase: req.Passphrase,
	}
	_, err := getService().AddExchangeApi(context.Background(), exAPi)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	response.NewSuccess(c, nil)
}

func GetApiKeyInfoHandler(c *gin.Context) {
	userID := c.Param("user_id")
	apiKey := c.Param("apiKey")
	req := &pb.UserApiKeyReq{
		UserId: userID,
		ApiKey: apiKey,
	}
	info, err := getService().GetApiKeyInfo(context.Background(), req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	resp := response.ExchangeApiInfoResp{
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
	apiList, err := getService().GetExchangeApiList(context.Background(), req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	var apiResp []*response.ExchangeApiResp
	if err := jsoniter.Unmarshal(apiList.ExApiList, &apiResp); err != nil {
		logger.Errorf("GetExApiList json Unmarshal err %v", err)
		response.NewInternalServerErr(c, nil)
		return
	}
	response.NewSuccess(c, apiResp)
}

func UpdateExApiHandler(c *gin.Context) {
	var req response.UpdateExchangeApiReq
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
		ApiId:      req.ApiId,
		UserId:     jwtP.UserID,
		ExchangeId: req.ExchangeID,
		ApiKey:     req.ApiKey,
		Secret:     req.Secret,
		Passphrase: req.Passphrase,
	}
	_, err := getService().UpdateExchangeApi(context.Background(), exAPi)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	response.NewSuccess(c, nil)
}

func DeleteExApiHandler(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	apiId := c.Param("apiId")
	i, err := strconv.Atoi(apiId)
	if err != nil {
		response.NewInternalServerErr(c, nil)
		return
	}
	req := &pb.UserApiReq{ApiId: int64(i), UserId: jwtP.UserID}
	_ ,err = getService().DeleteExchangeApi(context.Background(), req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	response.NewSuccess(c, nil)
}

func GetTradeListHandler(c *gin.Context) {
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
	resp, err := getService().GetTradeList(context.Background(), req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}

	var tradeResp []*response.TradeResp
	if err := jsoniter.Unmarshal(resp.TradeList, &tradeResp); err != nil {
		response.NewInternalServerErr(c, nil)
		return
	}
	response.NewSuccess(c, tradeResp)
}

func GetProfitHandler(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	strategyId := c.Param("strategyId")
	req := &pb.ProfitReq{
		UserId:     jwtP.UserID,
		StrategyId: strategyId,
	}
	resp, err := getService().GetProfitRealTime(context.Background(), req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	var profitResp []*response.ProfitResp
	if err := jsoniter.Unmarshal(resp.ProfitList, &profitResp); err != nil {
		response.NewInternalServerErr(c, nil)
		return
	}
	response.NewSuccess(c, profitResp)
}

func GetExchangePosHandler(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	req := &pb.GetExchangePosReq{UserId: jwtP.UserID}
	resp, err := getService().GetExchangePos(context.Background(), req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	response.NewSuccess(c, resp)
}

func PutOrderGridHandler(c *gin.Context) {
	var req pb.OrderReq
	defer c.Request.Body.Close()
	if err := jsonpb.Unmarshal(c.Request.Body, &req); err != nil {
		logger.Warnf("接受成交订单数据出错: %+v", err)
		response.NewBindJsonErr(c, nil)
		return
	}
	_, err := getService().EvaluationSpot(context.Background(), &req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	response.NewSuccess(c, nil)
}

func GetUserStrategyEvaHandler(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	strategyId := c.Param("strategyId")
	req := &pb.UserStrategyDetailReq{
		UserId:     jwtP.UserID,
		StrategyId: strategyId,
	}
	eva, err := getService().GetUserStrategyEva(context.Background(), req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	if len(eva.EvaDailyList) == 0 {
		response.NewSuccess(c, &response.UserStrategyEvaResp{
			RealizeProfit:  eva.RealizeProfit,
			RateReturnYear: eva.RateReturnYear,
			RateReturn:     eva.RateReturn,
			EvaDaily:       make([]interface{}, 0),
		})
		return
	}

	response.NewSuccess(c, eva)
}

func GetUserStrategyEvaNoAuthHandler(c *gin.Context) {
	userID := c.Param("user_id")
	strategyId := c.Param("strategyId")
	req := &pb.UserStrategyDetailReq{
		UserId:     userID,
		StrategyId: strategyId,
	}
	eva, err := getService().GetUserStrategyEva(context.Background(), req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	if len(eva.EvaDailyList) == 0 {
		response.NewSuccess(c, &response.UserStrategyEvaResp{
			RealizeProfit:  eva.RealizeProfit,
			RateReturnYear: eva.RateReturnYear,
			RateReturn:     eva.RateReturn,
			EvaDaily:       make([]interface{}, 0),
		})
		return
	}

	response.NewSuccess(c, eva)
}
