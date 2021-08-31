package service

import (
	"context"
	"fortune-bd/api/constant"
	pb "fortune-bd/api/quote/v1"
	"fortune-bd/app/quote-svc/cron"
	"fortune-bd/app/quote-svc/internal/dao"
	"fortune-bd/libs/logger"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type QuoteService struct {
	pb.UnimplementedQuoteServer
	dao *dao.Dao
}

func NewQuoteService() *QuoteService {
	handler := &QuoteService{dao: dao.New()}
	return handler
}

func (s *QuoteService) GetTicks(ctx context.Context, req *pb.GetTicksReq) *pb.TickResp {
	tickArrayAll := cron.BinanceTickArrayAll
	if req.Exchange == constant.HUOBI {
		tickArrayAll = cron.HuobiTickArrayAll
	}
	if req.Exchange == constant.BINANCE {
		tickArrayAll = cron.BinanceTickArrayAll
	}
	if len(tickArrayAll) == 0 {
		time.Sleep(2 * time.Second)
	}
	ticks, err := jsoniter.Marshal(tickArrayAll)
	if err != nil {
		logger.Infof("streamOkexTicks json Marshal err %v", err)
		return nil
	}
	return &pb.TickResp{Ticks: ticks}
}
