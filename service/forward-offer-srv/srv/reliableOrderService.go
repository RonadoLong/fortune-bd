package srv

import (
	"time"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/service/forward-offer-srv/global"
	"wq-fotune-backend/service/forward-offer-srv/srv/cache_service"
)

// listenReliableOrderQueue 处理可靠交易逻辑
func listenReliableOrderQueue() {
	logger.Infof("【 开始监听可靠交易事件 】")
	for {
		time.Sleep(time.Second)
		deq := cache_service.GetAutoAddOrderInfoFromQueue()
		if deq == nil {
			continue
		}
		logger.Infof("【 接收到可靠交易订单 】%v", global.StructToJsonStr(deq))
		client := loginExchange(&deq.ExchangeInfo)
		if client == nil {
			continue
		}
		client.ProcessReliableOrderEvent(deq)
	}
}
