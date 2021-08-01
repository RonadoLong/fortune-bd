package redis

import (
	"github.com/go-redis/redis"
	"sync"
	"time"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/libs/redisHelper"
)

var Redis *Client
var once sync.Once

type Client struct {
	Client *redis.Client
}

func InitRedis(host, pass string) {
	logger.Infof("redis host: %s", host)
	once.Do(func() {
		Redis = &Client{
			Client: redisHelper.InitRedisClient(host, pass),
		}
	})
}

func CacheSet(key string, value interface{}, expiration time.Duration) error {
	if err := Redis.Client.Set(key, value, expiration).Err(); err != nil {
		return err
	}
	return nil
}

func CacheGet(key string) ([]byte, error) {
	return Redis.Client.Get(key).Bytes()
}

func CacheDel(key string) {
	err := Redis.Client.Del(key).Err()
	if err != nil {
		logger.Errorf("redis CacheDel has error %v", err)
	}
}