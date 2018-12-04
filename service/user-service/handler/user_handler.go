package handler

import (
	"context"
	"shop-micro/helper"
	"shop-micro/service/user-service/proto"
)

type UserHandler struct {
	repo *userRepository
}

func NewUserHandler() (*UserHandler, error) {
	db, err := helper.CreateConnection()
	if err != nil {
		return nil, err
	}
	repository := &userRepository{DB:db}
	handler := &UserHandler{
		repo:repository,
	}
	return handler, nil
}

func (u *UserHandler) Login(c context.Context, req *shop_srv_user.LoginReq, resp *shop_srv_user.UserResp) error {
	if err := u.repo.login(req, resp); err != nil{
		return err
	}

	return nil
}

func (u *UserHandler) GetCode(c context.Context, req *shop_srv_user.PhoneCodeReq, resp *shop_srv_user.PhoneCodeResp) error {
	return nil	
}
