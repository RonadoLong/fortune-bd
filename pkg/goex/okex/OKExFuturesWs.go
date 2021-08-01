package okex

import (
	"encoding/json"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"sort"
	"strconv"
	"strings"
	"time"
	. "wq-fotune-backend/pkg/goex"
	"wq-fotune-backend/pkg/goex/internal/logger"
)

type OKExV3FuturesWs struct {
	base            *OKEx
	v3Ws            *OKExV3Ws
	tickerCallback  func(*FutureTicker)
	depthCallback   func(*Depth)
	orderCallback   func(*FutureOrder, string)
	tradeCallback   func(*Trade, string)
	klineCallback   func(*FutureKline, int)
	ErrorHandleFunc func(error)
}

func NewOKExV3FuturesWs(base *OKEx) *OKExV3FuturesWs {
	okV3Ws := &OKExV3FuturesWs{
		base: base,
	}
	okV3Ws.v3Ws = NewOKExV3Ws(base, okV3Ws.handle)
	return okV3Ws
}

func (okV3Ws *OKExV3FuturesWs) TickerCallback(tickerCallback func(*FutureTicker)) {
	okV3Ws.tickerCallback = tickerCallback
}

func (okV3Ws *OKExV3FuturesWs) DepthCallback(depthCallback func(*Depth)) {
	okV3Ws.depthCallback = depthCallback
}

func (okV3Ws *OKExV3FuturesWs) TradeCallback(tradeCallback func(*Trade, string)) {
	okV3Ws.tradeCallback = tradeCallback
}

func (okV3Ws *OKExV3FuturesWs) OrderCallback(orderCallback func(*FutureOrder, string)) {
	okV3Ws.orderCallback = orderCallback
}

func (okV3Ws *OKExV3FuturesWs) ErrorCallbacks(errCallback func(err error)) {
	okV3Ws.ErrorHandleFunc = errCallback
}

func (okV3Ws *OKExV3FuturesWs) KlineCallback(klineCallback func(*FutureKline, int)) {
	okV3Ws.klineCallback = klineCallback
}

func (okV3Ws *OKExV3FuturesWs) SetCallbacks(tickerCallback func(*FutureTicker),
	depthCallback func(*Depth),
	tradeCallback func(*Trade, string),
	klineCallback func(*FutureKline, int)) {
	okV3Ws.tickerCallback = tickerCallback
	okV3Ws.depthCallback = depthCallback
	okV3Ws.tradeCallback = tradeCallback
	okV3Ws.klineCallback = klineCallback
}

func (okV3Ws *OKExV3FuturesWs) getChannelName(currencyPair CurrencyPair, contractType string) string {
	var (
		prefix      string
		contractId  string
		channelName string
	)

	if contractType == SWAP_CONTRACT {
		prefix = "swap"
		contractId = fmt.Sprintf("%s-SWAP", currencyPair.ToSymbol("-"))
	} else {
		prefix = "futures"
		contractId = okV3Ws.base.OKExFuture.GetFutureContractId(currencyPair, contractType)
		logger.Info("contractid=", contractId)
	}

	channelName = prefix + "/%s:" + contractId

	return channelName
}

func (okV3Ws *OKExV3FuturesWs) SubscribeDepth(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.depthCallback == nil {
		return errors.New("please set depth callback func")
	}

	chName := okV3Ws.getChannelName(currencyPair, contractType)

	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf(chName, "depth5")}})
}

func (okV3Ws *OKExV3FuturesWs) SubscribeTicker(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.tickerCallback == nil {
		return errors.New("please set ticker callback func")
	}
	chName := okV3Ws.getChannelName(currencyPair, contractType)
	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf(chName, "ticker")}})
}

func (okV3Ws *OKExV3FuturesWs) SubscribeTrade(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.tradeCallback == nil {
		return errors.New("please set trade callback func")
	}
	chName := okV3Ws.getChannelName(currencyPair, contractType)
	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf(chName, "trade")}})
}

func (okV3Ws *OKExV3FuturesWs) SubscribeKline(currencyPair CurrencyPair, contractType string, period int) error {
	if okV3Ws.klineCallback == nil {
		return errors.New("place set kline callback func")
	}

	seconds := adaptKLinePeriod(KlinePeriod(period))
	if seconds == -1 {
		return fmt.Errorf("unsupported kline period %d in okex", period)
	}

	chName := okV3Ws.getChannelName(currencyPair, contractType)
	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf(chName, fmt.Sprintf("candle%ds", seconds))}})
}

func (okV3Ws *OKExV3FuturesWs) getContractAliasAndCurrencyPairFromInstrumentId(instrumentId string) (alias string, pair CurrencyPair) {
	if strings.HasSuffix(instrumentId, "SWAP") {
		ar := strings.Split(instrumentId, "-")
		return instrumentId, NewCurrencyPair2(fmt.Sprintf("%s_%s", ar[0], ar[1]))
	} else {
		contractInfo, err := okV3Ws.base.OKExFuture.GetContractInfo(instrumentId)
		if err != nil {
			logger.Error("instrument id invalid:", err)
			return "", UNKNOWN_PAIR
		}
		alias = contractInfo.Alias
		pair = NewCurrencyPair2(fmt.Sprintf("%s_%s", contractInfo.UnderlyingIndex, contractInfo.QuoteCurrency))
		return alias, pair
	}
}

func (okV3Ws *OKExV3FuturesWs) SubscribeOrder(symbol string, contractType string) error {
	if okV3Ws.orderCallback == nil {
		return errors.New("place set order callback func")
	}
	chName := fmt.Sprintf("%s/order:%s", okV3Ws.v3Ws.getTablePrefix(CurrencyPair{}, contractType), symbol)
	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{chName}})
}

func (okV3Ws *OKExV3FuturesWs) handle(channel string, data json.RawMessage) error {
	var (
		err           error
		ch            string
		tickers       []tickerResponse
		depthResp     []depthResponse
		dep           Depth
		tradeResponse []struct {
			Side         string  `json:"side"`
			TradeId      int64   `json:"trade_id,string"`
			Price        float64 `json:"price,string"`
			Qty          float64 `json:"qty,string"`
			InstrumentId string  `json:"instrument_id"`
			Timestamp    string  `json:"timestamp"`
		}
		klineResponse []struct {
			Candle       []string `json:"candle"`
			InstrumentId string   `json:"instrument_id"`
		}
	)

	if strings.Contains(channel, "error") {
		okV3Ws.ErrorHandleFunc(errors.New(string(data)))
		return nil
	}

	if strings.Contains(channel, "futures/candle") {
		ch = "candle"
	} else {
		ch, err = okV3Ws.v3Ws.parseChannel(channel)
		if err != nil {
			logger.Errorf("[%s] parse channel err=%s ,  originChannel=%s", okV3Ws.base.GetExchangeName(), err, ch)
			return nil
		}
	}

	switch ch {
	case "ticker":
		err = json.Unmarshal(data, &tickers)
		if err != nil {
			return err
		}

		for _, t := range tickers {
			alias, pair := okV3Ws.getContractAliasAndCurrencyPairFromInstrumentId(t.InstrumentId)
			date, _ := time.Parse(time.RFC3339, t.Timestamp)
			okV3Ws.tickerCallback(&FutureTicker{
				Ticker: &Ticker{
					Pair: pair,
					Last: t.Last,
					Buy:  t.BestBid,
					Sell: t.BestAsk,
					High: t.High24h,
					Low:  t.Low24h,
					Vol:  t.Volume24h,
					Date: uint64(date.UnixNano() / int64(time.Millisecond)),
				},
				ContractType: alias,
			})
		}
		return nil
	case "candle":
		err = json.Unmarshal(data, &klineResponse)
		if err != nil {
			return err
		}

		for _, t := range klineResponse {
			_, pair := okV3Ws.getContractAliasAndCurrencyPairFromInstrumentId(t.InstrumentId)
			ts, _ := time.Parse(time.RFC3339, t.Candle[0])
			//granularity := adaptKLinePeriod(KlinePeriod(period))
			okV3Ws.klineCallback(&FutureKline{
				Kline: &Kline{
					Pair:      pair,
					High:      ToFloat64(t.Candle[2]),
					Low:       ToFloat64(t.Candle[3]),
					Timestamp: ts.Unix(),
					Open:      ToFloat64(t.Candle[1]),
					Close:     ToFloat64(t.Candle[4]),
					Vol:       ToFloat64(t.Candle[5]),
				},
				Vol2: ToFloat64(t.Candle[6]),
			}, 1)
		}
		return nil
	case "depth5":
		err := json.Unmarshal(data, &depthResp)
		if err != nil {
			logger.Error(err)
			return err
		}
		if len(depthResp) == 0 {
			return nil
		}
		alias, pair := okV3Ws.getContractAliasAndCurrencyPairFromInstrumentId(depthResp[0].InstrumentId)
		dep.Pair = pair
		dep.ContractType = alias
		dep.UTime, _ = time.Parse(time.RFC3339, depthResp[0].Timestamp)
		for _, itm := range depthResp[0].Asks {
			dep.AskList = append(dep.AskList, DepthRecord{
				Price:  ToFloat64(itm[0]),
				Amount: ToFloat64(itm[1])})
		}
		for _, itm := range depthResp[0].Bids {
			dep.BidList = append(dep.BidList, DepthRecord{
				Price:  ToFloat64(itm[0]),
				Amount: ToFloat64(itm[1])})
		}
		sort.Sort(sort.Reverse(dep.AskList))
		//call back func
		okV3Ws.depthCallback(&dep)
		return nil
	case "trade":
		err := json.Unmarshal(data, &tradeResponse)
		if err != nil {
			logger.Error("unmarshal error :", err)
			return err
		}

		for _, resp := range tradeResponse {
			alias, pair := okV3Ws.getContractAliasAndCurrencyPairFromInstrumentId(resp.InstrumentId)

			tradeSide := SELL
			switch resp.Side {
			case "buy":
				tradeSide = BUY
			}

			t, err := time.Parse(time.RFC3339, resp.Timestamp)
			if err != nil {
				logger.Warn("parse timestamp error:", err)
			}

			okV3Ws.tradeCallback(&Trade{
				Tid:    resp.TradeId,
				Type:   tradeSide,
				Amount: resp.Qty,
				Price:  resp.Price,
				Date:   t.Unix(),
				Pair:   pair,
			}, alias)
		}
		return nil
	case "order":
		var orders []map[string]interface{}
		err = jsoniter.Unmarshal(data, &orders)
		if err != nil {
			logger.Error(err)
			return err
		}
		for _, order := range orders {
			futureOrder, s, err := okV3Ws.parseFutureOrder(order)
			if err != nil {
				return err
			}
			if okV3Ws.orderCallback != nil {
				okV3Ws.orderCallback(futureOrder, s)
			}
		}
		return nil
	}

	return fmt.Errorf("[%s] unknown websocket message: %s", ch, string(data))
}

func (okV3Ws *OKExV3FuturesWs) getKlinePeriodFormChannel(channel string) int {
	metas := strings.Split(channel, ":")
	if len(metas) != 2 {
		return 0
	}
	i, _ := strconv.ParseInt(metas[1], 10, 64)
	return int(i)
}

func getSign(apiSecretKey, timestamp, method, url, body string) (string, error) {
	data := timestamp + method + url + body
	return GetParamHmacSHA256Base64Sign(apiSecretKey, data)
}

func (okV3Ws *OKExV3FuturesWs) getTimestamp() string {
	seconds := float64(time.Now().UTC().UnixNano()) / float64(time.Second)
	return fmt.Sprintf("%.3f", seconds)
}

func (okV3Ws *OKExV3FuturesWs) Login() error {
	okV3Ws.v3Ws.ConnectWs()
	apiKey := okV3Ws.base.Config.ApiKey
	passphrase := okV3Ws.base.Config.ApiPassphrase
	apiSecretKey := okV3Ws.base.Config.ApiSecretKey
	//clear last login result
	timestamp := okV3Ws.getTimestamp()
	method := "GET"
	url := "/users/self/verify"
	sign, _ := getSign(apiSecretKey, timestamp, method, url, "")
	op := map[string]interface{}{
		"op":   "login",
		"args": []string{apiKey, passphrase, timestamp, sign}}
	return okV3Ws.v3Ws.Subscribe(op)
}

func (okV3Ws *OKExV3FuturesWs) parseFutureOrder(data interface{}) (*FutureOrder, string, error) {
	var fallback *FutureOrder
	switch v := data.(type) {
	case map[string]interface{}:
		var err error
		contractID := v["instrument_id"].(string)
		futureOrder := &FutureOrder{}
		futureOrder.ClientOid = v["client_oid"].(string)
		futureOrder.OrderID, _ = strconv.ParseInt(v["order_id"].(string), 10, 64)
		futureOrder.OrderID2 = v["order_id"].(string)
		futureOrder.Amount, err = strconv.ParseFloat(v["size"].(string), 64)
		if err != nil {
			return fallback, "", err
		}
		futureOrder.Price, err = strconv.ParseFloat(v["price"].(string), 64)
		if err != nil {
			return fallback, "", err
		}
		futureOrder.AvgPrice, err = strconv.ParseFloat(v["price_avg"].(string), 64)
		if err != nil {
			return fallback, "", err
		}
		futureOrder.DealAmount, err = strconv.ParseFloat(v["filled_qty"].(string), 64)
		if err != nil {
			return fallback, "", err
		}
		futureOrder.Fee, err = strconv.ParseFloat(v["fee"].(string), 64)
		if err != nil {
			return fallback, "", err
		}
		if i, err := strconv.ParseInt(v["type"].(string), 10, 64); err == nil {
			futureOrder.OType = int(i)
		} else {
			return fallback, "", err
		}
		futureOrder.OrderTime, err = timeStringToInt64(v["timestamp"].(string))
		if err != nil {
			return fallback, "", err
		}
		if v["leverage"] != nil {
			i, err := strconv.ParseFloat(v["leverage"].(string), 64)
			if err != nil {
				return fallback, "", err
			}
			futureOrder.LeverRate = int(i)
		}
		futureOrder.LastFillPrice, err = strconv.ParseFloat(v["last_fill_px"].(string), 64)
		if err != nil {
			return fallback, "", err
		}
		futureOrder.LastFillAmount, err = strconv.ParseFloat(v["last_fill_qty"].(string), 64)
		if err != nil {
			return fallback, "", err
		}
		futureOrder.LastFillTime, err = timeStringToTime(v["last_fill_time"].(string))
		if err != nil {
			return fallback, "", err
		}
		futureOrder.ContractName = v["instrument_id"].(string)
		//futureOrder.Currency = currencyPair
		if tId, ok := v["last_fill_id"]; ok {
			futureOrder.LastFillId = tId.(string)
		}
		state, err := strconv.ParseInt(v["state"].(string), 10, 64)
		if err != nil {
			return fallback, "", err
		}
		switch state {
		case 0:
			futureOrder.Status = ORDER_UNFINISH
		case 1:
			futureOrder.Status = ORDER_PART_FINISH
		case 2:
			futureOrder.Status = ORDER_FINISH
		case 4:
			futureOrder.Status = ORDER_CANCEL_ING
		case -1:
			futureOrder.Status = ORDER_CANCEL
		case 3:
			futureOrder.Status = ORDER_UNFINISH
		case -2:
			futureOrder.Status = ORDER_REJECT
		default:
			return fallback, contractID, fmt.Errorf("unknown order status: %v", v)
		}
		return futureOrder, contractID, nil
	}

	return fallback, "", fmt.Errorf("unknown FutureOrder data: %v", data)
}

func timeStringToTime(t string) (time.Time, error) {
	timestamp, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return time.Time{}, err
	}
	return timestamp, nil
}
