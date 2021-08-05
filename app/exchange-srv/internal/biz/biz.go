package biz

import (
	"github.com/go-redis/redis"
	"wq-fotune-backend/app/exchange-srv/cache"
	"wq-fotune-backend/app/exchange-srv/internal/dao"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/redisHelper"
)

const (
	ErrID = "exchangeOrder"
)

type ExOrderRepo struct {
	dao          *dao.Dao
	cacheService *cache.Service
}

func NewExOrderRepo() *ExOrderRepo {
	return &ExOrderRepo{
		dao:          dao.New(),
		cacheService: cache.NewService(),
	}
}

type ForwardOfferRepo struct {
	cacheService *redis.Client
	dao          *dao.Dao
}

func NewForwardOfferRepo() *ForwardOfferRepo {
	return &ForwardOfferRepo{
		dao:          dao.New(),
		cacheService: redisHelper.InitRedisClient(env.RedisAddr, env.RedisPWD),
	}
}
