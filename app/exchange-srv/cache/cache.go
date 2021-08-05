package cache

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
	"wq-fotune-backend/app/forward-offer-srv/global"
	"wq-fotune-backend/app/quote-srv/cron"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/libs/redisHelper"
)

const (
	exchangeKey = "ex:account:"
)

type Service struct {
	redisCli *redis.Client
}

func NewService() *Service {
	return &Service{
		redisCli: redisHelper.InitRedisClient(env.RedisAddr, env.RedisPWD),
	}
}

// CacheExchangeAccountList 缓存交易所数据
func (s *Service) CacheExchangeAccountList(userId string, data []byte) {
	var key = global.StringJoinString(exchangeKey, userId)
	s.redisCli.Set(key, data, time.Second*10)
}

// GetExchangeAccountList 获取用户在缓存中的数据
func (s *Service) GetExchangeAccountList(userId string) []byte {
	var key = global.StringJoinString(exchangeKey, userId)
	result, err := s.redisCli.Get(key).Bytes()
	if err != nil {
		return nil
	}
	return result
}

func (s *Service) GetOKexQuote(symbol string) (*cron.Ticker, error) {
	data := &cron.Ticker{}
	bytes, err := s.redisCli.HGet(cron.TickOkexAll, symbol).Bytes()
	if err != nil {
		logger.Warnf("reids获取tick:okex此品种获取失败 %s %v", symbol, err)
		return nil, err
	}
	if err := json.Unmarshal(bytes, data); err != nil {
		logger.Warnf("reids获取tick:okex此品种解析失败 %s %v", symbol, err)
		return nil, err
	}
	return data, nil
}

func (s *Service) GetHuobiQuote(symbol string) (*cron.Ticker, error) {
	data := &cron.Ticker{}
	bytes, err := s.redisCli.HGet(cron.TickHuobiAll, symbol).Bytes()
	if err != nil {
		logger.Warnf("reids获取tick:huobi此品种获取失败 %s %v", symbol, err)
		return nil, err
	}
	if err := json.Unmarshal(bytes, data); err != nil {
		logger.Warnf("reids获取tick:huobi此品种解析失败 %s %v", symbol, err)
		return nil, err
	}
	return data, nil
}

func (s *Service) GetBinanceQuote(symbol string) (*cron.Ticker, error) {
	data := &cron.Ticker{}
	bytes, err := s.redisCli.HGet(cron.TickBinanceAll, symbol).Bytes()
	if err != nil {
		logger.Warnf("reids获取tick:binance此品种获取失败 %s %v", symbol, err)
		return nil, err
	}
	if err := json.Unmarshal(bytes, data); err != nil {
		logger.Warnf("reids获取tick:binance此品种解析失败 %s %v", symbol, err)
		return nil, err
	}
	return data, nil
}

func (s *Service) CacheData(key string, data interface{}, duration time.Duration) error {
	return s.redisCli.Set(key, data, duration).Err()
}
