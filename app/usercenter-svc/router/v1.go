package router

import (
	"context"
	"fortune-bd/api/response"
	pb "fortune-bd/api/usercenter/v1"
	"fortune-bd/app/usercenter-svc/internal/service"
	"fortune-bd/libs/jwt"
	"fortune-bd/libs/logger"
	"fortune-bd/libs/middleware"
	"fortune-bd/libs/validate-code/phone"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"

)

var (
	userService *service.UserService
)

func apiV1(group *gin.RouterGroup) {
	userService = service.NewUserService()
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
	var req response.LoginReq
	if err := c.BindJSON(&req); err != nil {
		response.NewBindJsonErr(c, nil)
		return
	}
	loginReq := &pb.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	}
	user,err := userService.Login(context.Background(), loginReq)
	if err != nil {
		fromError := errors.FromError(err)
		logger.Errorf("userService.Login  调用失败 %v", err.Error())
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	token, err := middleware.NewToken(user.UserId, middleware.RoleUser)
	if err != nil {
		logger.Errorf("生成token失败 用户 %s, 角色 %d", user.UserId, middleware.RoleUser)
		response.NewInternalServerErr(c, nil)
		return
	}
	lastLogin, _ := ptypes.Timestamp(user.LastLoginAt)
	response.NewSuccess(c, &response.LoginResp{
		UserId:         user.UserId,
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
	var req response.ValidateCodeReq
	if err := c.BindJSON(&req); err != nil {
		response.NewBindJsonErr(c, nil)
		return
	}
	if !phone.CheckPhone(req.Phone) {
		response.NewErrorParam(c, "手机号格式错误！", nil)
		return
	}

	vcodeReq := &pb.ValidateCodeReq{Phone: req.Phone}
	resp, err := userService.SendValidateCode(context.Background(), vcodeReq)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	response.NewSuccess(c, resp.Code)
}

func Register(c *gin.Context) {
	var req response.RegisterReq
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
	if _, err := userService.Register(context.Background(), registerReq); err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	response.NewSuccess(c, nil)
}

func ResetPassword(c *gin.Context) {
	var req response.ChangePasswordReq
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
	resetReq := &pb.ChangePasswordReq{UserId: jwtP.UserID, Password: req.Password, ConfirmPassword: req.ConfirmPassword}
	if _,err := userService.ResetPassword(context.Background(), resetReq); err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	response.NewSuccess(c, nil)
}

func UpdateUser(c *gin.Context) {
	var req response.UpdateUserBaseReq
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
	if _, err := userService.UpdateUser(context.Background(), updateUserReq); err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	response.NewSuccess(c, nil)
}

func ForgetPassword(c *gin.Context) {
	var req response.ForgetPasswordReq
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

	if _, err := userService.ForgetPassword(context.Background(), changePassReq); err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	response.NewSuccess(c, nil)
}

func BaseInfo(c *gin.Context) {
	claims, _ := c.Get("claims")
	JwtP := claims.(*jwt.JWTPayload)
	req := &pb.UserInfoReq{UserId: JwtP.UserID}
	user, err := userService.GetUserInfo(context.Background(), req)
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	lastLogin, _ := ptypes.Timestamp(user.LastLoginAt)
	response.NewSuccess(c, &response.LoginResp{
		UserId:         user.UserId,
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
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	response.NewSuccess(c, resp.Members)
}

func GetPaymentMethods(c *gin.Context) {
	resp, err := userService.GetPaymentMethod(context.Background(), &empty.Empty{})
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	response.NewSuccess(c, resp.Payments)
}
