package redisHelper

import (
	"github.com/go-redis/redis"
	"time"
	"wq-fotune-backend/libs/logger"
)

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
	RedisClient := redis.NewClient(options)
	go func() {
		var ticker = time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_, err := RedisClient.Ping().Result()
				if err != nil {
					logger.Errorf("conned redisUtils err %s", err.Error())
				}
			}
		}
	}()
	return RedisClient
}
