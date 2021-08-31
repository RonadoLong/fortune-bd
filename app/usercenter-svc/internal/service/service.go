package service

import (
	"context"
	"errors"
	walletpb "fortune-bd/api/wallet/v1"
	"fortune-bd/app/usercenter-svc/internal/model"
)

const (
	errID = "user"
)


func (u *UserService) GetUserInfoByPhone(phone string) (*model.WqUserBase, error) {
	user := u.dao.GetWqUserBaseByPhone(phone)
	if user == nil {
		return nil, errors.New("该手机没有绑定用户")
	}
	return user, nil
}

func (u *UserService) AddIfcBalance(userMasterId, inUserID, exchange, _type string, volume float64) error {
	_, err := u.walletSrv.AddIfcBalance(context.Background(), &walletpb.AddIfcBalanceReq{
		UserMasterId: userMasterId,
		InUserId:     inUserID,
		Volume:       volume,   //手数
		Type:         _type,    //register api  strategy
		Exchange:     exchange, //注册这里留空
	})
	return err
}