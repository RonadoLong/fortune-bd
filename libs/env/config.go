package env

import (
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	runMode   = "RUN_MODE"
	redisAddr = "REDIS_ADDR"
	redisPWD  = "REDIS_PWD"
	dbDSN     = "DB_DSN"
	etcdAddr  = "ETCD_ADDR"
	mongoAddr = "MONGO_ADDR"

	mode = "MODE"
)

var (
	configMap = initial()

	RunMode   = configMap.getValue(runMode)
	RedisAddr = configMap.getValue(redisAddr)
	RedisPWD  = configMap.getValue(redisPWD)
	DBDSN     = configMap.getValue(dbDSN)
	EtcdAddr  = configMap.getValue(etcdAddr)
	MongoAddr = configMap.getValue(mongoAddr)

	GridNum = 10
	// 上报异常到消息中心URL
	MsgCenterURL = "http://localhost:20080/v1/dataCollect/systemMsg"
	// 获取交易所访问授权信息URL，url后面需要加参数 /:userid/:apikey
	ExchangeAccessURL = "http://test.ifortune.io/api/v1/exchange-order/exchange/apiInfo"
	// 获取统计信息URL url后面需要加参数 /:userId/:strategyId
	StatisticalInfoURL = "http://test.ifortune.io/api/v1/exchange-order/user/strategy/evaluationNoAuth"
	// 通知统计URL
	NotifyStatisticsURL = "http://test.ifortune.io/api/v1/exchange-order/forward-offer/orderGrid"
	// 启动策略通知接口
	NotifyStrategyStartUpURL = "http://test.ifortune.io/api/v1/wallet/strategyStartUpNotify"
	ProxyAddr                = "socks5://192.168.3.30:20170"

	EXCHANGE_ORDER_SRV_NAME = "exchange-order.srv"
	USER_SRV_NAME           = "usercenter.srv"
	WALLET_SRV_NAME         = "wallet.srv"
	QUOTE_SRV_NAME          = "quote.srv"
	COMMON_SRV_NAME         = "common.srv"
)

type envConfig map[string]string

func initial() envConfig {
	var config envConfig
	releaseMode := os.Getenv(mode)
	if releaseMode == "production" {
		config = proEnv
	} else if releaseMode == "release" {
		config = releaseEnv
	} else {
		config = developEnv
	}
	return config
}

func (env envConfig) getValue(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	if v, ok := env[key]; ok {
		return v
	}
	return ""
}

func GetProxyHttpClient() *http.Client {
	client := &http.Client{}
	client.Transport = &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return &url.URL{
				Scheme: "socks5",
				Host:   strings.Split(ProxyAddr, "//")[1],
			}, nil
		},
	}
	return client
}
