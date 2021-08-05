package service

import (
	"wq-fotune-backend/app/usercenter-srv/internal/dao"
	walletCli "wq-fotune-backend/app/wallet-srv/client"
	walletPb "wq-fotune-backend/app/wallet-srv/proto"
	"wq-fotune-backend/libs/env"
)

type UserService struct {
	dao       *dao.Dao
	walletSrv walletPb.WalletService
}

const (
	errID = "user"
)

// NewUserService biz
func NewUserService() *UserService {
	handler := &UserService{dao.New(), walletCli.NewWalletClient(env.EtcdAddr)}
	return handler
}
