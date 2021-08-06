package biz

import (
	"github.com/go-redis/redis"
	"wq-fotune-backend/app/exchange-srv/internal/dao"
)

const (
	ErrID = "exchangeOrder"
)

type ExOrderRepo struct {
	dao          *dao.Dao
}

func NewExOrderRepo() *ExOrderRepo {
	return &ExOrderRepo{
		dao:          dao.New(),
	}
}

type ForwardOfferRepo struct {
	cacheService *redis.Client
	dao          *dao.Dao
}

func NewForwardOfferRepo() *ForwardOfferRepo {
	return &ForwardOfferRepo{
		dao:          dao.New(),
	}
}
