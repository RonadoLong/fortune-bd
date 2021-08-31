package biz

import (
	pb "fortune-bd/api/usercenter/v1"
	"fortune-bd/app/wallet-svc/cache"
	"fortune-bd/app/wallet-svc/internal/dao"
	"fortune-bd/libs/exchangeclient"
)

const (
	ErrID = "wallet"
	//BinanceApiKey = "fev72IlrChwPbO8Yp3D57RkvIiuUwkIFK3dJoQi7cQaYyv00DiBwxDiXm4DH4HZq"
	//BinanceSecret = "YGvVhns0OlIxMJ1of4apa0IeYGbXsFvCrbewrTYveQz0qfxDhRalBfBJd7EUN4iP"
	BinanceApiKey = "lfNGLnHexoDNXEYeQGApIWb75ItHm7w7zOCJpxp1vvODIQFOFwChmuHxhvoleb1d"
	BinanceSecret = "8G3X3a3NxsZAyh3ZmEYRIX3d5DKK6PyqXyC6JylA0CQiQtafMZ8AUa8v8gRq43Sz"

)

type WalletRepo struct {
	dao          *dao.Dao
	cacheService *cache.Service
	binance      *exchangeclient.Binance
	UserSrv      pb.UserClient
}

func NewWalletRepo() *WalletRepo {
	return &WalletRepo{
		dao:          dao.New(),
		cacheService: cache.NewService(),
		binance:      exchangeclient.InitBinance(BinanceApiKey, BinanceSecret),
		//UserSrv:      userCli.NewUserClient(env.EtcdAddr),
	}
}
