package limitReq

import (
	"strconv"
	"time"
	"wq-fotune-backend/libs/cache"
)

const (
	limitReq = "limit:requests:"
)

func SetReqCount(key string, count int) error {
	timeOut := time.Second * 15
	return cache.CacheSet(limitReq+key, count, timeOut)
}
func SetReqCountWithTimeOut(key string, count int, timeout time.Duration) error {
	timeOut := time.Second * timeout
	return cache.CacheSet(limitReq+key, count, timeOut)
}

func GetReqCount(key string) int {
	data, _ := cache.CacheGet(limitReq + key)
	count, _ := strconv.Atoi(string(data))
	return count
}
