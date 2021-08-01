package service

import (
	"github.com/go-redis/redis"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/redisHelper"
	"wq-fotune-backend/service/exchange-srv/cache"
	"wq-fotune-backend/service/exchange-srv/dao"
)

const (
	ErrID = "exchangeOrder"
)

type ExOrderService struct {
	dao          *dao.Dao
	cacheService *cache.Service
}

func NewExOrderService() *ExOrderService {
	return &ExOrderService{
		dao:          dao.New(),
		cacheService: cache.NewService(),
	}
}

type ForwardOfferService struct {
	cacheService *redis.Client
	dao          *dao.Dao
}

func NewForwardOfferService() *ForwardOfferService {
	return &ForwardOfferService{
		dao:          dao.New(),
		cacheService: redisHelper.InitRedisClient(env.RedisAddr, env.RedisPWD),
	}
}
