package cache

import (
	"fortune-bd/libs/env"
	"fortune-bd/libs/logger"
	"github.com/go-redis/redis"
	"time"

)

var rdb *redis.Client

func Redis() *redis.Client{
	if rdb == nil {
		rdb = InitRedisClient(env.RedisAddr, env.RedisPWD)
	}
	return rdb
}

func InitRedisClient(host string, password string) *redis.Client {
	var options = &redis.Options{
		Addr:        host,
		DB:          0,
		IdleTimeout: 240 * time.Second,
		DialTimeout: time.Second * 10,
	}
	if password != "" {
		options.Password = password
	}
	return redis.NewClient(options)
}


func CacheSet(key string, value interface{}, expiration time.Duration) error {
	if err := rdb.Set(key, value, expiration).Err(); err != nil {
		return err
	}
	return nil
}

func CacheGet(key string) ([]byte, error) {
	return rdb.Get(key).Bytes()
}

func CacheDel(key string) {
	err := rdb.Del(key).Err()
	if err != nil {
		logger.Errorf("redis CacheDel has error %v", err)
	}
}