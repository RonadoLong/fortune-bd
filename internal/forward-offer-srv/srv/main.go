package srv

import (
	"sync"
	"wq-fotune-backend/internal/forward-offer-srv/config"
	"wq-fotune-backend/internal/forward-offer-srv/global"
)

var once = sync.Once{}

// ListeningOrderService 开始监听事件
func ListeningOrderService() {
	once.Do(func() {
		global.InitRedisClient(config.Config.RedisHost, config.Config.RedisPass)
		//cache_service.InitDelayerClient()
		go listenReliableOrderQueue()
		go runOrderReceiver()
	})
}
