package dao

import (
	"fortune-bd/libs/cache"
	"github.com/go-redis/redis"
)

type Dao struct {
	RedisCli *redis.Client
}

func New() *Dao {
	return &Dao{
		RedisCli: cache.Redis(),
	}
}
