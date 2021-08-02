package handler

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/empty"
	"time"
	"wq-fotune-backend/libs/logger"
	exchange_info "wq-fotune-backend/pkg/exchange-info"
	"wq-fotune-backend/pkg/response"
	"wq-fotune-backend/internal/quote-srv/cron"
	"wq-fotune-backend/internal/quote-srv/dao"
	fotune_srv_quote "wq-fotune-backend/internal/quote-srv/proto"
)

const (
	errID = "quote"
	USDT  = "usdt"
	BTC   = "btc"
)

type QuoteHandler struct {
	dao *dao.Dao
}

func (q *QuoteHandler) GetTicksWithExchangeSymbol(ctx context.Context, req *fotune_srv_quote.GetTicksSymbolReq, resp *fotune_srv_quote.TickResp) error {
	var tickArrayAll = make([]cron.Ticker, 0)
	//var tickerAll = make(map[string]map[string]interface{})
	if req.Exchange == exchange_info.BINANCE {
		if req.Symbol == USDT {
			tickArrayAll = cron.BinanceTickArrayAll
		}
		if req.Symbol == BTC {
			tickArrayAll = cron.BinaceTickArrayBtc
		}
	}
	if req.Exchange == exchange_info.HUOBI {
		if req.Symbol == USDT {
			tickArrayAll = cron.HuobiTickArrayAll
		}
		if req.Symbol == BTC {
			tickArrayAll = cron.HuobiTickArrayBtc
		}
	}
	if req.Exchange == exchange_info.OKEX {
		if req.Symbol == USDT {
			tickArrayAll = cron.OkexTickArrayAll
		}
		if req.Symbol == BTC {
			tickArrayAll = cron.OkexTickArrayBtc
		}
	}

	ticks, err := json.Marshal(tickArrayAll)
	if err != nil {
		return response.NewInternalServerErrMsg(errID)
	}
	resp.Ticks = ticks
	return nil
}

func (q *QuoteHandler) GetTicksWithExchange(ctx context.Context, req *fotune_srv_quote.GetTicksReq, resp *fotune_srv_quote.TickResp) error {
	//if req.All == false {
	//	if len(cron.OkexTickArray) == 0 {
	//		return response.NewDataNotFound(errID, "行情数据更新失败")
	//	}
	//	ticks, err := json.Marshal(cron.OkexTickArray)
	//	if err != nil {
	//		return response.NewInternalServerErrMsg(errID)
	//	}
	//	resp.Ticks = ticks
	//	return nil
	//}
	//tickArrayAll := cron.OkexTickArrayAll
	//if req.Exchange == exchange_info.HUOBI {
	//	tickArrayAll = cron.HuobiTickArrayAll
	//}
	//if req.Exchange == exchange_info.BINANCE {
	//	tickArrayAll = cron.BinanceTickArrayAll
	//}

	var tickerAll = make(map[string]map[string]interface{})
	tickerAll[exchange_info.BINANCE] = map[string]interface{}{
		"usdt": cron.BinanceTickMapAll,
	}
	ticks, err := json.Marshal(tickerAll)
	if err != nil {
		return response.NewInternalServerErrMsg(errID)
	}
	resp.Ticks = ticks
	return nil
}

func (q *QuoteHandler) GetRate(ctx context.Context, e *empty.Empty, rmb *fotune_srv_quote.RateUsdRmb) error {
	bytes, err := q.dao.RedisCli.Get(cron.RateKey).Bytes()
	if err != nil {
		logger.Warnf("redis取汇率错误 %v", err)
		return response.NewInternalServerErrMsg(errID)
	}
	var rate cron.QuoteRate
	if err := json.Unmarshal(bytes, &rate); err != nil {
		logger.Warnf("解析redis获取的汇率数据失败 %v", err)
		return response.NewInternalServerErrMsg(errID)
	}
	rmb.InstrumentID = rate.InstrumentID
	rmb.Rate = rate.Rate
	rmb.Timestamp = rate.Timestamp
	return nil
}

func NewQuoteHandler() *QuoteHandler {
	handler := &QuoteHandler{dao: dao.New()}
	return handler
}

//func (q *QuoteHandler) GetOkexTicks(ctx context.Context, req *empty.Empty, resp *fotune_srv_quote.OkexTickResp) error {
//}

func (q *QuoteHandler) StreamOkexTicks(ctx context.Context, req *fotune_srv_quote.GetTicksReq, resp fotune_srv_quote.QuoteService_StreamOkexTicksStream) error {
	for {
		tickArrayAll := cron.OkexTickArrayAll
		//tickArrayBtc := cron.OkexTickArrayBtc

		if req.Exchange == exchange_info.HUOBI {
			tickArrayAll = cron.HuobiTickArrayAll
			//tickArrayBtc = cron.HuobiTickArrayBtc
		}
		if req.Exchange == exchange_info.BINANCE {
			tickArrayAll = cron.BinanceTickArrayAll
			//tickArrayBtc = cron.BinaceTickArrayBtc
		}

		//if len(tickArrayAll) == 0 || len(tickArrayBtc) == 0{
		if len(tickArrayAll) == 0 {
			time.Sleep(2 * time.Second)
			continue
		}
		ticks, err := json.Marshal(tickArrayAll)
		if err != nil {
			logger.Infof("streamOkexTicks json Marshal err %v", err)
			return err
			//return  response.NewInternalServerErrMsg(errID)
		}
		//TODO 币本位行情
		//dataMap := make(map[string][]cron.Ticker)
		//dataMap["USDT"] = tickArrayAll
		//dataMap["BTC"] = tickArrayBtc

		//data := make([]map[string][]cron.Ticker, 0)
		//data = append(data, dataMap)
		//ticks, err := json.Marshal(data) //{"usdt":[tickArrayALl的数据]}
		//if err != nil {
		//	logger.Infof("streamOkexTicks json Marshal err %v", err)
		//	return err
		//	//return  response.NewInternalServerErrMsg(errID)
		//}
		err = resp.Send(&fotune_srv_quote.TickResp{Ticks: ticks})
		if err != nil {
			logger.Infof("streamOkexTicks sendMsg err %v", err)
			return err
		}
		time.Sleep(6 * time.Second)
		continue
	}
}
