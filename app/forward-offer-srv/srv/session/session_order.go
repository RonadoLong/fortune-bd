package session

import (
	"time"
	"wq-fotune-backend/app/forward-offer-srv/global"
	"wq-fotune-backend/app/forward-offer-srv/srv/cache_service"
	"wq-fotune-backend/app/forward-offer-srv/srv/model"
)

const maxCount = 200

var errMap = map[string]int{
	"30012": global.CollectAccountErrorCode,
	"30010": global.CollectAccountErrorCode,
	"30026": global.CollectOrderLimiterErrorCode,
	"30014": global.CollectOrderLimiterErrorCode,
	"35047": global.CollectPriceErrorCode,
	"32080": global.CollectPriceErrorCode,
	"35002": global.CollectOrderLimiterErrorCode,
	"35004": global.CollectOrderLimiterErrorCode, // 合约正在结算时
}

// 提交订单
func (c *Client) ProcessCreateOrder(req model.OrderReq, tyrCount int) *model.CheckError {
	_, err := c.ApiClient.PostOrder(req)
	if err != nil {
		if code, ok := errMap[err.Desc.ErrorCode]; ok && err.Code != 1500 {
			if code == global.CollectOrderLimiterErrorCode { // 判断是否可以重试的订单
				if tyrCount >= maxCount {
					return model.CreateCheckError(code, global.StructToJsonStr(err.Desc)) // 请求频繁
				}
				time.Sleep(time.Second)
				tyrCount++
				c.ProcessCreateOrder(req, tyrCount)
			}
			return model.CreateCheckError(code, global.StructToJsonStr(err.Desc))
		}
		return model.CreateCheckError(global.CollectOrderParamsErrorCode, global.StructToJsonStr(err.Desc))
	}
	return nil
}

// 发送失败订单到延迟队列
func (c *Client) pushReliableOrderToWaitQueue(req model.OrderReq, deq *model.DelayerReq) {
	delayerReq := model.CreateDelayerReq(req, deq.UserID, deq.StrategyID)
	delayerReq.APIKey = c.ApiClient.GetApiKey()
	delayerReq.SecretKey = c.ApiClient.GetApiSecretKey()
	delayerReq.EcPass = c.ApiClient.GetApiPassphrase()
	cache_service.PushDelayerInfoToQueue(delayerReq)
}
