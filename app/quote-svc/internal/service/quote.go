package service

import (
	"context"
	"fortune-bd/api/constant"
	"fortune-bd/api/response"
	"fortune-bd/app/quote-svc/cron"
	"fortune-bd/libs/logger"
	jsoniter "github.com/json-iterator/go"
	"time"

	pb "fortune-bd/api/quote/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	errID = "quote"
	USDT  = "usdt"
	BTC   = "btc"
)

func (s *QuoteService) GetTicksWithExchange(ctx context.Context, req *pb.GetTicksReq) (*pb.TickResp, error) {
	var resp *pb.TickResp
	var tickerAll = make(map[string]map[string]interface{})
	tickerAll[constant.BINANCE] = map[string]interface{}{
		"usdt": cron.BinanceTickMapAll,
	}
	ticks, err := jsoniter.Marshal(tickerAll)
	if err != nil {
		return nil, response.NewInternalServerErrMsg(errID)
	}
	resp.Ticks = ticks
	return resp, nil
}

func (s *QuoteService) GetTicksWithExchangeSymbol(ctx context.Context, req *pb.GetTicksSymbolReq) (*pb.TickResp, error) {
	var resp *pb.TickResp
	logger.Infof("GetTicksWithExchangeSymbol: %+v", req)
	var tickArrayAll = make([]cron.Ticker, 0)
	if req.Exchange == constant.BINANCE {
		if req.Symbol == USDT {
			tickArrayAll = cron.BinanceTickArrayAll
		}
		if req.Symbol == BTC {
			tickArrayAll = cron.BinaceTickArrayBtc
		}
	}
	if req.Exchange == constant.HUOBI {
		if req.Symbol == USDT {
			tickArrayAll = cron.HuobiTickArrayAll
		}
		if req.Symbol == BTC {
			tickArrayAll = cron.HuobiTickArrayBtc
		}
	}
	ticks, err := jsoniter.Marshal(tickArrayAll)
	if err != nil {
		return nil, response.NewInternalServerErrMsg(errID)
	}
	resp.Ticks = ticks
	logger.Infof("获取行情接口: %+v", string(resp.Ticks))
	return resp, nil
}

func (s *QuoteService) StreamTicks(req *pb.GetTicksReq, conn pb.Quote_StreamTicksServer) error {
	for {
		tickArrayAll := cron.BinanceTickArrayAll
		if req.Exchange == constant.HUOBI {
			tickArrayAll = cron.HuobiTickArrayAll
		}
		if req.Exchange == constant.BINANCE {
			tickArrayAll = cron.BinanceTickArrayAll
		}
		if len(tickArrayAll) == 0 {
			time.Sleep(2 * time.Second)
			continue
		}
		ticks, err := jsoniter.Marshal(tickArrayAll)
		if err != nil {
			logger.Infof("streamOkexTicks json Marshal err %v", err)
			return err
		}
		err = conn.Send(&pb.TickResp{Ticks: ticks})
		if err != nil {
			logger.Infof("streamOkexTicks sendMsg err %v", err)
			return err
		}
		time.Sleep(6 * time.Second)
		continue
	}
}

func (s *QuoteService) GetRate(ctx context.Context, req *emptypb.Empty) (*pb.RateUsdRmb, error) {
	var resp *pb.RateUsdRmb
	bytes, err := s.dao.RedisCli.Get(cron.RateKey).Bytes()
	if err != nil {
		logger.Warnf("redis取汇率错误 %v", err)
		return nil, response.NewInternalServerErrMsg(errID)
	}
	var rate cron.QuoteRate
	if err := jsoniter.Unmarshal(bytes, &rate); err != nil {
		logger.Warnf("解析redis获取的汇率数据失败 %v", err)
		return nil, response.NewInternalServerErrMsg(errID)
	}
	resp.InstrumentId = rate.InstrumentID
	resp.Rate = rate.Rate
	resp.Timestamp = rate.Timestamp
	return resp, nil
}
