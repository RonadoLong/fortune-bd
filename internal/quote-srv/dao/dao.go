package dao

import (
	"github.com/go-redis/redis"
	"wq-fotune-backend/libs/global"
)

type Dao struct {
	//db *gorm.DB
	RedisCli *redis.Client
}

func New() *Dao {
	return &Dao{
		RedisCli: global.RedisCli,
	}
}
