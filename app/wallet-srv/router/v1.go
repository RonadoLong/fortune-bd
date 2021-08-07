package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/micro/go-micro/v2/errors"
	"log"
	"wq-fotune-backend/api/response"
	pb "wq-fotune-backend/api/wallet"
	"wq-fotune-backend/app/wallet-srv/internal/service"
	"wq-fotune-backend/libs/jwt"
	"wq-fotune-backend/libs/middleware"
)

var (
	walletService = service.NewWalletService()
)


func apiV1(group *gin.RouterGroup) {
	group.POST("/strategyStartUpNotify", StrategyStarUpNotify)
	group.Use(middleware.JWTAuth())
	group.POST("/create", CreateWallet)
	group.POST("/transfer", Transfer)
	group.GET("/ifc", WalletIFC)
	group.GET("/usdt", WalletUSDT)
	group.GET("usdt/depositAddr", GetUsdtDeposit)
	group.POST("/convertCoinTips", CovertCoinTips)
	group.POST("/convertCoin", CovertCoin)
	group.POST("/withdrawal", CreateWithdrawal)
	group.GET("/totalRebate", GetTotalRebate)
}




func StrategyStarUpNotify(c *gin.Context) {
	var req pb.StrategyRunNotifyReq
	defer c.Request.Body.Close()
	if err := jsonpb.Unmarshal(c.Request.Body, &req); err != nil {
		log.Println(err)
		response.NewBindJsonErr(c, nil)
		return
	}
	err := walletService.StrategyRunNotify(context.Background(), &req, nil)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, nil)
}

func GetTotalRebate(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	var rebate pb.GetTotalRebateResp
	err := walletService.GetTotalRebate(context.Background(), &pb.GetTotalRebateReq{UserId: jwtP.UserID}, &rebate)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, rebate)
}

func CreateWithdrawal(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	var req pb.WithdrawalReq
	defer c.Request.Body.Close()
	if err := jsonpb.Unmarshal(c.Request.Body, &req); err != nil {
		log.Println(err)
		response.NewBindJsonErr(c, nil)
		return
	}
	req.UserId = jwtP.UserID
	err := walletService.Withdrawal(context.Background(), &req, nil)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, nil)
}

func CovertCoin(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	var req pb.ConvertCoinReq
	defer c.Request.Body.Close()
	if err := jsonpb.Unmarshal(c.Request.Body, &req); err != nil {
		response.NewBindJsonErr(c, nil)
		return
	}
	req.UserId = jwtP.UserID
	var coin *pb.ConvertCoinResp
	err := walletService.ConvertCoin(context.Background(), &req, nil)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, coin)
}

func CovertCoinTips(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	var req pb.ConvertCoinTipsReq
	defer c.Request.Body.Close()
	if err := jsonpb.Unmarshal(c.Request.Body, &req); err != nil {
		response.NewBindJsonErr(c, nil)
		return
	}
	req.UserId = jwtP.UserID
	var tips pb.ConvertCoinTipsResp
	err := walletService.ConvertCoinTips(context.Background(), &req, &tips)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, tips)
}

func GetUsdtDeposit(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	var resp pb.UsdtDepositAddrResp
	err := walletService.GetUsdtDepositAddr(context.Background(), &pb.UidReq{UserId: jwtP.UserID}, &resp)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, resp)
}

func CreateWallet(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	err := walletService.CreateWallet(context.Background(), &pb.UidReq{UserId: jwtP.UserID}, nil)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, nil)
}

func Transfer(c *gin.Context) {
	var req pb.TransferReq
	defer c.Request.Body.Close()
	if err := jsonpb.Unmarshal(c.Request.Body, &req); err != nil {
		response.NewBindJsonErr(c, nil)
		return
	}
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	err := walletService.Transfer(context.Background(), &pb.TransferReq{
		UserId:         jwtP.UserID,
		FromCoin:       req.FromCoin,
		ToCoin:         req.ToCoin,
		FromCoinAmount: req.FromCoinAmount,
	}, nil)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, nil)
}

func WalletIFC(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	var wallet pb.WalletBalanceResp
	err := walletService.GetWalletIFC(context.Background(), &pb.UidReq{UserId: jwtP.UserID}, &wallet)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, wallet)
}

func WalletUSDT(c *gin.Context) {
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	var wallet pb.WalletBalanceResp
	err := walletService.GetWalletUSDT(context.Background(), &pb.UidReq{UserId: jwtP.UserID}, &wallet)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, wallet)
}
