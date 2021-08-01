package cache

import (
	"github.com/go-redis/redis"
	"time"
	"wq-fotune-backend/libs/global"
)

type Service struct {
	redisCli *redis.Client
}

func NewService() *Service {
	return &Service{
		redisCli: global.RedisCli,
	}
}

func (s *Service) CacheData(key string, data interface{}, duration time.Duration) error {
	return s.redisCli.Set(key, data, duration).Err()
}

func (s *Service) GetData(key string) ([]byte, error) {
	return s.redisCli.Get(key).Bytes()
}
