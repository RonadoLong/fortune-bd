package cache

import (
	"encoding/json"
	"time"
	"wq-fotune-backend/app/forward-offer-srv/global"
	"wq-fotune-backend/app/quote-srv/cron"
	"wq-fotune-backend/libs/cache"
	"wq-fotune-backend/libs/logger"
)

const (
	exchangeKey = "ex:account:"
)



// CacheExchangeAccountList 缓存交易所数据
func CacheExchangeAccountList(userId string, data []byte) {
	var key = global.StringJoinString(exchangeKey, userId)
	cache.Redis().Set(key, data, time.Second*10)
}

// GetExchangeAccountList 获取用户在缓存中的数据
func GetExchangeAccountList(userId string) []byte {
	var key = global.StringJoinString(exchangeKey, userId)
	result, err := cache.Redis().Get(key).Bytes()
	if err != nil {
		return nil
	}
	return result
}

func GetOKexQuote(symbol string) (*cron.Ticker, error) {
	data := &cron.Ticker{}
	bytes, err := cache.Redis().HGet(cron.TickOkexAll, symbol).Bytes()
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

func  GetHuobiQuote(symbol string) (*cron.Ticker, error) {
	data := &cron.Ticker{}
	bytes, err := cache.Redis().HGet(cron.TickHuobiAll, symbol).Bytes()
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

func GetBinanceQuote(symbol string) (*cron.Ticker, error) {
	data := &cron.Ticker{}
	bytes, err := cache.Redis().HGet(cron.TickBinanceAll, symbol).Bytes()
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

func CacheData(key string, data interface{}, duration time.Duration) error {
	return cache.Redis().Set(key, data, duration).Err()
}
