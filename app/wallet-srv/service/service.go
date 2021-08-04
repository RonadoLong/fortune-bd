package service

import (
	apiBinance "wq-fotune-backend/libs/binance_client"
	"wq-fotune-backend/libs/env"
	userCli "wq-fotune-backend/app/usercenter-srv/client"
	fotune_srv_user "wq-fotune-backend/app/usercenter-srv/proto"
	"wq-fotune-backend/app/wallet-srv/cache"
	"wq-fotune-backend/app/wallet-srv/dao"
)

const (
	ErrID = "wallet"
	//BinanceApiKey = "fev72IlrChwPbO8Yp3D57RkvIiuUwkIFK3dJoQi7cQaYyv00DiBwxDiXm4DH4HZq"
	//BinanceSecret = "YGvVhns0OlIxMJ1of4apa0IeYGbXsFvCrbewrTYveQz0qfxDhRalBfBJd7EUN4iP"
	BinanceApiKey = "lfNGLnHexoDNXEYeQGApIWb75ItHm7w7zOCJpxp1vvODIQFOFwChmuHxhvoleb1d"
	BinanceSecret = "8G3X3a3NxsZAyh3ZmEYRIX3d5DKK6PyqXyC6JylA0CQiQtafMZ8AUa8v8gRq43Sz"

)

type WalletService struct {
	dao          *dao.Dao
	cacheService *cache.Service
	binance      *apiBinance.BinanceClient
	UserSrv      fotune_srv_user.UserService
}

func NewWalletService() *WalletService {
	return &WalletService{
		dao:          dao.New(),
		cacheService: cache.NewService(),
		binance:      apiBinance.InitClient(BinanceApiKey, BinanceSecret),
		UserSrv:      userCli.NewUserClient(env.EtcdAddr),
	}
}
