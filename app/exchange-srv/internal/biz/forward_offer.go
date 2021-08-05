package biz

import (
	"encoding/hex"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	pb "wq-fotune-backend/api/exchange"
	"wq-fotune-backend/app/forward-offer-srv/global"
	"wq-fotune-backend/app/forward-offer-srv/srv/model"
	"wq-fotune-backend/libs/helper"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/pkg/encoding"
	"wq-fotune-backend/pkg/utils"
)

const TradeRequestQueue = "trade:%s:request"

// PushSwapOrder direction 不再使用
func (f *ForwardOfferRepo) PushSwapOrder(req *pb.TradeSignal) error {
	if req == nil || req.Orders == nil {
		return nil
	}
	//todo 检查是否有策略组ID的策略在运行-> 后期需要过滤组ID，品种(打开注释部分)
	strategyOfRun := f.dao.GetUserStrategyOfRun(nil)
	if len(strategyOfRun) == 0 {
		logger.Warnf("没有在运行的策略")
		return nil
	}
	for _, strategy := range strategyOfRun {
		logger.Infof("正在运行的策略: %v uid: %s", strategy.ID, strategy.UID)
		//if req.SharedID != strategy.GroupID || req.Orders.Symbol != strategy.Symbol {
		//	continue
		//}
		platform := f.dao.GetExchangeApiListByUidAndPlatform(strategy.UID, "okex")
		if len(platform) == 0 {
			continue
		}
		//todo 需要封装和优化
		// 组装交易服务需要的数据 如这个样例service/forward-offer-srv/srv/session/session_test.go
		val := &model.OrderReq{}
		val.OrderID = global.GetOkexOrderID()
		val.Symbol = req.Orders.Symbol
		//val.OrderQty = global.StringToFloat64(string(req.Orders.OrderQty))
		val.OrderQty = 1
		val.Price = req.Orders.Price
		val.TryCount = 10
		val.Direction = req.Orders.Side
		marshalToString, _ := jsoniter.MarshalToString(val)
		d := map[string]string{
			"type":  "autoAdd",
			"value": marshalToString,
		}
		data, _ := jsoniter.MarshalToString(d)
		var tradeReq = model.ExchangeReq{}
		tradeReq.Exchange = "okex"
		tradeReq.Data = data
		tradeReq.UserID = strategy.UID
		tradeReq.StrategyID = strategy.ID
		secret, _ := hex.DecodeString(platform[0].Secret)
		secretKey, _ := encoding.AesDecrypt(secret)
		tradeReq.ExchangeInfo = &model.ExchangeInfo{
			APIKey:    platform[0].ApiKey,
			SecretKey: string(secretKey),
			EcPass:    platform[0].Passphrase,
		}
		key := fmt.Sprintf(TradeRequestQueue, "okex")
		bytes, _ := jsoniter.Marshal(tradeReq)
		// 发送订单到队列
		f.PushReqOrderMessage(key, bytes)
	}
	return nil
}

func (f *ForwardOfferRepo) PushReqOrderMessage(key string, msg []byte) {
	logger.Warnf("push order to queue：%s msg: %+v", key, string(msg))
	err := utils.ReTryFunc(10, func() (bool, error) {
		err := f.cacheService.RPush(key, msg).Err()
		if err != nil {
			logger.Errorf("Push msg to queue: %s", err)
		}
		return false, err
	})
	if err != nil {
		logStr := helper.StringJoinString("重试发送订单到队列失败3次, 请检查是否redis出现问题")
		logger.Error(logStr)
	}
}
