package global

import (
	"wq-fotune-backend/libs/helper"
	"wq-fotune-backend/libs/logger"
)

// GetReqOrderMessageFromQueue 消费队列数据
func GetReqOrderMessageFromQueue() string {
	ret, err := RedisClient.LPop(Queue).Result()
	if err != nil && err.Error() != "redis: nil" {
		logger.Errorf("GetReqOrderMessageFromQueue err: %s", err.Error())
	}
	return ret
}

func PushReqOrderMessage(key string, msg []byte) {
	logger.Warnf("push order to queue：%s msg: %+v", key, string(msg))
	err := RunRetry(3, func() error {
		err := RedisClient.RPush(key, msg).Err()
		if err != nil {
			logger.Errorf("Push msg to queue: %s", err)
		}
		return err
	})
	if err != nil {
		logStr := helper.StringJoinString("重试发送订单到队列失败3次, 请检查是否redis出现问题")
		logger.Error(logStr)
	}
}
