package limitReq

import (
	"strconv"
	"time"
	"wq-fotune-backend/pkg/redis"
)

const (
	limitReq = "limit:requests:"
)

func SetReqCount(key string, count int) error {
	timeOut := time.Second * 15
	return redis.CacheSet(limitReq+key, count, timeOut)
}
func SetReqCountWithTimeOut(key string, count int, timeout time.Duration) error {
	timeOut := time.Second * timeout
	return redis.CacheSet(limitReq+key, count, timeOut)
}

func GetReqCount(key string) int {
	data, _ := redis.CacheGet(limitReq + key)
	count, _ := strconv.Atoi(string(data))
	return count
}
