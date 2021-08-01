package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/micro/go-micro/v2/errors"
	"wq-fotune-backend/api-gateway/protocol"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/jwt"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/pkg/middleware"
	"wq-fotune-backend/pkg/response"
	"wq-fotune-backend/pkg/validate-code/phone"
	"wq-fotune-backend/service/user-srv/client"
	fotune_srv_user "wq-fotune-backend/service/user-srv/proto"
)

var (
	userService fotune_srv_user.UserService
)

func InitUserEngine(engine *gin.RouterGroup) {
	userService = client.NewUserClient(env.EtcdAddr)
	group := engine.Group("/user")
	group.POST("/login", Login)
	group.POST("/send/validate-code", SendValidateCode)
	group.POST("/register", Register)
	group.PUT("/forget/password", ForgetPassword)
	group.GET("/members", GetMembers)
	group.GET("/paymentMethods", GetPaymentMethods)

	group.Use(middleware.JWTAuth())
	group.PUT("/update", UpdateUser)
	group.PUT("/reset/password", ResetPassword)
	group.GET("/base-info", BaseInfo)
}

func Login(c *gin.Context) {
	//
	var req protocol.LoginReq
	if err := c.BindJSON(&req); err != nil {
		response.NewBindJsonErr(c, nil)
		return
	}
	loginReq := &fotune_srv_user.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	}
	user, err := userService.Login(context.Background(), loginReq)
	if err != nil {
		fromError := errors.FromError(err)
		logger.Errorf("userService.Login  调用失败 %v", err.Error())
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	token, err := middleware.NewToken(user.UserID, middleware.RoleUser)
	if err != nil {
		logger.Errorf("生成token失败 用户 %s, 角色 %d", user.UserID, middleware.RoleUser)
		response.NewInternalServerErr(c, nil)
		return
	}
	lastLogin, _ := ptypes.Timestamp(user.LastLogin)
	response.NewSuccess(c, &protocol.LoginResp{
		UserId:         user.UserID,
		Token:          token,
		InvitationCode: user.InvitationCode,
		Name:           user.Name,
		Avatar:         user.Avatar,
		Phone:          user.Phone,
		LastLogin:      lastLogin,
		LoginCount:     user.LoginCount,
	})
}

func SendValidateCode(c *gin.Context) {
	var req protocol.ValidateCodeReq
	if err := c.BindJSON(&req); err != nil {
		response.NewBindJsonErr(c, nil)
		return
	}
	if !phone.CheckPhone(req.Phone) {
		response.NewErrorParam(c, "手机号格式错误！", nil)
		return
	}

	vcodeReq := &fotune_srv_user.ValidateCodeReq{Phone: req.Phone}
	resp, err := userService.SendValidateCode(context.Background(), vcodeReq)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, resp.Code)
}

func Register(c *gin.Context) {
	var req protocol.RegisterReq
	if err := c.BindJSON(&req); err != nil {
		response.NewBindJsonErr(c, nil)
		return
	}
	//检验基本参数
	if err := req.CheckBaseParam(); err != nil {
		response.NewErrorParam(c, err.Error(), nil)
		return
	}
	registerReq := &fotune_srv_user.RegisterReq{
		Phone:           req.Phone,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
		InvitationCode:  req.InvitationCode,
		ValidateCode:    req.ValidateCode,
	}
	if _, err := userService.Register(context.Background(), registerReq); err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, nil)
}

func ResetPassword(c *gin.Context) {
	var req protocol.ChangePasswordReq
	if err := c.BindJSON(&req); err != nil {
		response.NewBindJsonErr(c, nil)
		return
	}
	if err := req.CheckPassword(); err != nil {
		response.NewErrorParam(c, err.Error(), nil)
		return
	}
	claims, _ := c.Get("claims")
	jwtP := claims.(*jwt.JWTPayload)
	resetReq := &fotune_srv_user.ChangePasswordReq{UserID: jwtP.UserID, Password: req.Password, ConfirmPassword: req.ConfirmPassword}
	if _, err := userService.ResetPassword(context.Background(), resetReq); err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, nil)
}

func UpdateUser(c *gin.Context) {
	var req protocol.UpdateUserBaseReq
	if err := c.BindJSON(&req); err != nil {
		response.NewBindJsonErr(c, nil)
		return
	}
	if err := req.CheckNotNull(); err != nil {
		response.NewErrorParam(c, err.Error(), nil)
		return
	}
	claims, _ := c.Get("claims")
	JwtP := claims.(*jwt.JWTPayload)
	updateUserReq := &fotune_srv_user.UpdateUserReq{
		Name:   req.Name,
		Avatar: req.Avatar,
		UserId: JwtP.UserID,
	}
	if _, err := userService.UpdateUser(context.Background(), updateUserReq); err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, nil)
}

func ForgetPassword(c *gin.Context) {
	var req protocol.ForgetPasswordReq
	if err := c.BindJSON(&req); err != nil {
		response.NewBindJsonErr(c, nil)
		return
	}
	if err := req.CheckPhoneVCode(); err != nil {
		response.NewErrorParam(c, err.Error(), nil)
		return
	}
	if err := req.CheckPassword(); err != nil {
		response.NewErrorParam(c, err.Error(), nil)
		return
	}

	changePassReq := &fotune_srv_user.ForgetPasswordReq{
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
		Phone:           req.Phone,
		ValidateCode:    req.ValidateCode,
	}

	if _, err := userService.ForgetPassword(context.Background(), changePassReq); err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, nil)
}

func BaseInfo(c *gin.Context) {
	claims, _ := c.Get("claims")
	JwtP := claims.(*jwt.JWTPayload)
	req := &fotune_srv_user.UserInfoReq{UserID: JwtP.UserID}
	user, err := userService.GetUserInfo(context.Background(), req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	lastLogin, _ := ptypes.Timestamp(user.LastLogin)
	response.NewSuccess(c, &protocol.LoginResp{
		UserId:         user.UserID,
		InvitationCode: user.InvitationCode,
		Name:           user.Name,
		Avatar:         user.Avatar,
		Phone:          user.Phone,
		LastLogin:      lastLogin,
		LoginCount:     user.LoginCount,
	})
}

func GetMembers(c *gin.Context) {
	resp, err := userService.GetMembers(context.Background(), &empty.Empty{})
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, resp.Members)
}

func GetPaymentMethods(c *gin.Context) {
	resp, err := userService.GetPaymentMethod(context.Background(), &empty.Empty{})
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, resp.Payments)
}
