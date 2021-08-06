package cache

import (
	"errors"
	"github.com/go-redis/redis"
	"time"
	"wq-fotune-backend/libs/cache"
	"wq-fotune-backend/libs/logger"
)

type Service struct {
	redisCli *redis.Client
}

var userRunKey = "run-strategy:"

func NewService() *Service {
	return &Service{
		redisCli: cache.Redis(),
	}
}

// CacheUserStrategyRunInfo 用户启动策略后存一条redis记录
func (s *Service) CacheUserStrategyRunInfo(userId string) {
	timeOut := time.Hour * 24 * 30
	if err := s.redisCli.Set(userRunKey+userId, "", timeOut).Err(); err != nil {
		logger.Errorf("用户启动数据redis保存失败 %s %v", userId, "")
	}
}

var KeyNotFound = errors.New("None ")

func (s *Service) GetUserStrategyRunInfo(userId string) (times string, err error) {
	times, err = s.redisCli.Get(userRunKey + userId).Result()
	if err != nil {
		if err == redis.Nil {
			logger.Warnf("没有找到启动记录 %s", userId)
			return "", KeyNotFound
		}
		return "", err
	}
	return times, nil
}
