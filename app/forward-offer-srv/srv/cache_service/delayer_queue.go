package cache_service

import (
	"fmt"
	"strings"
	"wq-fotune-backend/app/forward-offer-srv/global"
	"wq-fotune-backend/app/forward-offer-srv/srv/model"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/libs/utils"

	jsoniter "github.com/json-iterator/go"
)

// PushDelayerInfoToQueue send delayer message to queue
func PushDelayerInfoToQueue(req *model.DelayerReq) {
	defer func() {
		if recover() != nil {
			logger.Errorf("发送到延迟队列失败: %v", req)
		}
	}()
	// 通过已有连接创建客户端
	ret, _ := jsoniter.Marshal(req)
	req.PqTime = global.GetCurrentTime()
	_ = utils.ReTryFunc(10, func() (bool, error) {
		err := global.RedisClient.RPush(global.Topic, ret).Err()
		if err != nil {
			if strings.Contains(err.Error(), "nil") {
				return true, err
			}
			msg := fmt.Sprintf("%s %s", err.Error(), global.StructToJsonStr(req))
			logger.Error(msg)
			return false, err
		}
		return false, nil
	})
}

// GetAutoAddOrderInfoFromQueue get delayer message from queue
func GetAutoAddOrderInfoFromQueue() *model.DelayerReq {
	var message *model.DelayerReq
	_ = utils.ReTryFunc(10, func() (bool, error) {
		ret, err := global.RedisClient.LPop(global.Topic).Result()
		if err != nil {
			if strings.Contains(err.Error(), "nil") {
				return true, err
			}
			msg := fmt.Sprintf("GetAutoAddOrderInfoFromQueue: %+v", err.Error())
			logger.Warn(msg)
			return false, err
		}
		_ = jsoniter.UnmarshalFromString(ret, &message)
		return false, nil
	})
	if message != nil {
		if global.GetCurrentTime().Sub(message.PqTime).Seconds() >= 10 {
			return message
		}
		PushDelayerInfoToQueue(message)
	}
	return nil
}
