package cache_service

import (
	"strings"
	"time"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/service/forward-offer-srv/global"
)

const (
	CacheFinishOrderKey     = "reliable:success:"
	OrderBindUserKey        = "trade:bind:"
	CacheFinishOrderTimeOut = time.Minute * 10
	tradeKey                = "trade:"
)

// CacheFinishOrder 缓存成功的订单
func CacheFinishOrder(orderID, val string) {
	k := global.StringJoinString(CacheFinishOrderKey, orderID)
	err := global.RunRetry(3, func() error {
		return global.RedisClient.Set(k, val, CacheFinishOrderTimeOut).Err()
	})
	if err != nil {
		logger.Errorf("CacheFinishOrder error: %v", err.Error())
	}
}

func GetCacheFinishOrder(orderID string) string {
	k := global.StringJoinString(CacheFinishOrderKey, orderID)
	return global.RedisClient.Get(k).Val()
}

// 绑定订单ID - 用户ID 策略ID 用来做成交回调时的获取
func CacheOrderIDBindUser(orderID, uID, sID string) {
	k := global.StringJoinString(OrderBindUserKey, orderID)
	val := global.StringJoinString(uID, ",", sID)
	err := global.RunRetry(3, func() error {
		return global.RedisClient.Set(k, val, time.Minute*30).Err()
	})
	if err != nil {
		logger.Errorf("CacheFinishOrder error: %v", err.Error())
	}
}

// uid 用户ID sid 策略ID
func GetOrderIDBindUserVal(orderID string) (uid, sid string) {
	k := global.StringJoinString(OrderBindUserKey, orderID)
	val := global.RedisClient.Get(k).Val()
	if val != "" {
		vals := strings.Split(val, ",")
		if len(vals) == 2 {
			return vals[0], vals[1]
		}
	}
	return "", ""
}

func CacheAccountLeID(sID, symbol, lID string) {
	k := global.StringJoinString(tradeKey, sID, ":", symbol)
	err := global.RunRetry(3, func() error {
		return global.RedisClient.Set(k, lID, 0).Err()
	})
	if err != nil {
		logger.Errorf("CacheFinishOrder: %v", err.Error())
	}
}

func GetCacheAccountLeID(sID, symbol string) string {
	k := global.StringJoinString(tradeKey, sID, ":", symbol)
	return global.RedisClient.Get(k).Val()
}
