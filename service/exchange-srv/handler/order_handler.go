package handler

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	pb "wq-fotune-backend/service/exchange-srv/proto"
	"wq-fotune-backend/service/exchange-srv/service"
)

type ForwardOfferHandle struct {
	forWardOfferService *service.ForwardOfferService
}

func NewForwardOfferHandle() *ForwardOfferHandle {
	return &ForwardOfferHandle{
		forWardOfferService: service.NewForwardOfferService(),
	}
}

func (f *ForwardOfferHandle) PushSwapOrder(ctx context.Context, req *pb.TradeSignal, empty *empty.Empty) error {
	return f.forWardOfferService.PushSwapOrder(req)
}

////PushReqOrderMessage 发送订单
//func (f *ForwardOfferHandle) PushReqOrderMessage(key string, msg []byte) {
//	logger.Warnf("push order to queue：%s msg: %+v", key, string(msg))
//	err := utils.ReTryFunc(10, func() (bool, error) {
//		err := f.RedisClient.RPush(key, msg).Err()
//		if err != nil {
//			logger.Errorf("Push msg to queue: %s", err)
//		}
//		return false, err
//	})
//	if err != nil {
//		logStr := helper.StringJoinString("重试发送订单到队列失败3次, 请检查是否redis出现问题")
//		logger.Error(logStr)
//	}
//}
