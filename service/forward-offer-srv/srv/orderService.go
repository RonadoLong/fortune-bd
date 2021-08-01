package srv

import (
	jsoniter "github.com/json-iterator/go"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/service/forward-offer-srv/global"
	"wq-fotune-backend/service/forward-offer-srv/srv/model"
	"wq-fotune-backend/service/forward-offer-srv/srv/session"
)

// runOrderReceiver 消费交易相关的数据
func runOrderReceiver() {
	logger.Info("【 开始消费队列中的请求数据 】")
	for {
		msg := global.GetReqOrderMessageFromQueue()
		if msg == "" {
			continue
		}
		logger.Infof("【 消费者接收到订单请求 】：%v", msg)
		var req *model.ExchangeReq
		err := jsoniter.UnmarshalFromString(msg, &req)
		if err != nil {
			logger.Infof("序列化结构体失败: %s", err.Error())
			logger.Info(msg)
			continue
		}
		client := loginExchange(req.ExchangeInfo)
		err = jsoniter.UnmarshalFromString(req.Data, &req.TradeReq)
		if err != nil {
			logger.Infof("序列化结构体失败: %s", err.Error())
			logger.Info(msg)
			continue
		}
		if client != nil {
			go client.BroadCastOrder(req)
		}
	}
}

func loginExchange(exchangeInfo *model.ExchangeInfo) *session.Client {
	if exchangeInfo == nil {
		logger.Warn("交易所账户为空")
		return nil
	}
	if exchangeInfo.EcPass == "" {
		logger.Errorf("login exchange: %s %s", exchangeInfo.APIKey, "PASSPHRASE 不能为空")
		return nil
	}
	client := session.InitClientNormal(exchangeInfo)
	//client := session.RegisterClientOrGetClient()
	if client == nil {
		logger.Errorf("login exchange :%s", exchangeInfo.APIKey, "交易账号不对或者不存在")
	}
	return client
}
