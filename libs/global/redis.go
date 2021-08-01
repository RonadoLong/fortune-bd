package global

import (
	"github.com/go-redis/redis"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/redisHelper"
)

var RedisCli *redis.Client

func InitRedis() {
	RedisCli = redisHelper.InitRedisClient(env.RedisAddr, env.RedisPWD)
}
