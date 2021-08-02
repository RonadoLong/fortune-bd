package cache

import (
	"testing"
	"wq-fotune-backend/libs/redisHelper"
)

func TestService_GetOKexQuote(t *testing.T) {
	redisSrv := &Service{redisCli: redisHelper.InitRedisClient("127.0.0.1:6379", "")}
	quote, err := redisSrv.GetOKexQuote("EOS-USDT")
	if err != nil {
		t.Log("none")
		return
	}
	t.Logf("%+v", quote)
}
