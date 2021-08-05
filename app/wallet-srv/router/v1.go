package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/micro/go-micro/v2/errors"
	"log"
	pb "wq-fotune-backend/api/wallet"
	"wq-fotune-backend/app/wallet-srv/client"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/jwt"
	"wq-fotune-backend/pkg/middleware"
	"wq-fotune-backend/pkg/response"
)

var (
	walletService pb.WalletService
)


func apiV1(group *gin.RouterGroup) {
	walletService = client.NewWalletClient(env.EtcdAddr)
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
	_, err := walletService.StrategyRunNotify(context.Background(), &req)
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
	rebate, err := walletService.GetTotalRebate(context.Background(), &pb.GetTotalRebateReq{UserId: jwtP.UserID})
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
	_, err := walletService.Withdrawal(context.Background(), &req)
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
	coin, err := walletService.ConvertCoin(context.Background(), &req)
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
	tips, err := walletService.ConvertCoinTips(context.Background(), &req)
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
	resp, err := walletService.GetUsdtDepositAddr(context.Background(), &pb.UidReq{UserId: jwtP.UserID})
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
	_, err := walletService.CreateWallet(context.Background(), &pb.UidReq{UserId: jwtP.UserID})
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
	_, err := walletService.Transfer(context.Background(), &pb.TransferReq{
		UserId:         jwtP.UserID,
		FromCoin:       req.FromCoin,
		ToCoin:         req.ToCoin,
		FromCoinAmount: req.FromCoinAmount,
	})
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
	wallet, err := walletService.GetWalletIFC(context.Background(), &pb.UidReq{UserId: jwtP.UserID})
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
	wallet, err := walletService.GetWalletUSDT(context.Background(), &pb.UidReq{UserId: jwtP.UserID})
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, wallet)
}
