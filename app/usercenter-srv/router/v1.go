package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/micro/go-micro/v2/errors"
	"wq-fotune-backend/api/protocol"
	"wq-fotune-backend/api/response"
	pb "wq-fotune-backend/api/usercenter"
	"wq-fotune-backend/app/usercenter-srv/internal/service"
	"wq-fotune-backend/libs/jwt"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/libs/middleware"
	"wq-fotune-backend/libs/validate-code/phone"
)

var (
	userService = service.NewUserService()
)

func apiV1(group *gin.RouterGroup) {
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
	loginReq := &pb.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	}
	var user pb.LoginResp
	err := userService.Login(context.Background(), loginReq, &user)
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

	vcodeReq := &pb.ValidateCodeReq{Phone: req.Phone}
	var resp pb.ValidateCodeResp
	err := userService.SendValidateCode(context.Background(), vcodeReq, &resp)
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
	registerReq := &pb.RegisterReq{
		Phone:           req.Phone,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
		InvitationCode:  req.InvitationCode,
		ValidateCode:    req.ValidateCode,
	}
	if err := userService.Register(context.Background(), registerReq, nil); err != nil {
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
	resetReq := &pb.ChangePasswordReq{UserID: jwtP.UserID, Password: req.Password, ConfirmPassword: req.ConfirmPassword}
	if err := userService.ResetPassword(context.Background(), resetReq, nil); err != nil {
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
	updateUserReq := &pb.UpdateUserReq{
		Name:   req.Name,
		Avatar: req.Avatar,
		UserId: JwtP.UserID,
	}
	if err := userService.UpdateUser(context.Background(), updateUserReq, nil); err != nil {
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

	changePassReq := &pb.ForgetPasswordReq{
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
		Phone:           req.Phone,
		ValidateCode:    req.ValidateCode,
	}

	if err := userService.ForgetPassword(context.Background(), changePassReq, nil); err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, nil)
}

func BaseInfo(c *gin.Context) {
	claims, _ := c.Get("claims")
	JwtP := claims.(*jwt.JWTPayload)
	req := &pb.UserInfoReq{UserID: JwtP.UserID}
	var user pb.LoginResp
	err := userService.GetUserInfo(context.Background(), req, &user)
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
	var resp pb.GetMembersResp
	err := userService.GetMembers(context.Background(), &empty.Empty{}, &resp)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, resp.Members)
}

func GetPaymentMethods(c *gin.Context) {
	var resp pb.GetPaymentMethodResp
	err := userService.GetPaymentMethod(context.Background(), &empty.Empty{}, &resp)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Detail)
		return
	}
	response.NewSuccess(c, resp.Payments)
}
