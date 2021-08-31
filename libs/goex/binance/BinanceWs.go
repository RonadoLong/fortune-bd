package binance

import (
	"bytes"
	"errors"
	"fmt"
	"fortune-bd/libs/helper"
	"github.com/json-iterator/go"
	"log"
	"strconv"
	"strings"
	"time"
	"unsafe"
	. "fortune-bd/libs/goex"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type BinanceWs struct {
	baseURL         string
	combinedBaseURL string
	proxyUrl        string
	tickerCallback  func(*Ticker)
	depthCallback   func(*Depth)
	tradeCallback   func(*Trade)
	klineCallback   func(*Kline, int)
	wsConns         []*WsConn
}

type AggTrade struct {
	Trade
	FirstBreakdownTradeID int64 `json:"f"`
	LastBreakdownTradeID  int64 `json:"l"`
	TradeTime             int64 `json:"T"`
}

type RawTrade struct {
	Trade
	BuyerOrderID  int64 `json:"b"`
	SellerOrderID int64 `json:"a"`
}

type DiffDepth struct {
	Depth
	UpdateID      int64 `json:"u"`
	FirstUpdateID int64 `json:"U"`
}

var _INERNAL_KLINE_PERIOD_REVERTER = map[string]int{
	"1m":  KLINE_PERIOD_1MIN,
	"3m":  KLINE_PERIOD_3MIN,
	"5m":  KLINE_PERIOD_5MIN,
	"15m": KLINE_PERIOD_15MIN,
	"30m": KLINE_PERIOD_30MIN,
	"1h":  KLINE_PERIOD_60MIN,
	"2h":  KLINE_PERIOD_2H,
	"4h":  KLINE_PERIOD_4H,
	"6h":  KLINE_PERIOD_6H,
	"8h":  KLINE_PERIOD_8H,
	"12h": KLINE_PERIOD_12H,
	"1d":  KLINE_PERIOD_1DAY,
	"3d":  KLINE_PERIOD_3DAY,
	"1w":  KLINE_PERIOD_1WEEK,
	"1M":  KLINE_PERIOD_1MONTH,
}

func NewBinanceWs() *BinanceWs {
	bnWs := &BinanceWs{}
	bnWs.baseURL = "wss://stream.binance.com:9443/ws"
	bnWs.combinedBaseURL = "wss://stream.binance.com:9443/stream?streams="
	return bnWs
}

func (bnWs *BinanceWs) ProxyUrl(proxyUrl string) {
	bnWs.proxyUrl = proxyUrl
}

func (bnWs *BinanceWs) SetBaseUrl(baseURL string) {
	bnWs.baseURL = baseURL
}

func (bnWs *BinanceWs) SetCombinedBaseURL(combinedBaseURL string) {
	bnWs.combinedBaseURL = combinedBaseURL
}

func (bnWs *BinanceWs) SetCallbacks(
	tickerCallback func(*Ticker),
	depthCallback func(*Depth),
	tradeCallback func(*Trade),
	klineCallback func(*Kline, int),
) {
	bnWs.tickerCallback = tickerCallback
	bnWs.depthCallback = depthCallback
	bnWs.tradeCallback = tradeCallback
	bnWs.klineCallback = klineCallback
}

func (bnWs *BinanceWs) Subscribe(endpoint string, handle func(msg []byte) error) *WsConn {
	wsConn := NewWsBuilder().
		WsUrl(endpoint).
		AutoReconnect().
		ProtoHandleFunc(handle).
		ProxyUrl(bnWs.proxyUrl).
		ReconnectInterval(time.Millisecond * 5).
		Build()
	bnWs.wsConns = append(bnWs.wsConns, wsConn)
	go bnWs.exitHandler(wsConn)
	return wsConn
}

func (bnWs *BinanceWs) Close() {
	for _, con := range bnWs.wsConns {
		con.CloseWs()
	}
}

func (bnWs *BinanceWs) SubscribeDepth(pair CurrencyPair, size int) error {
	if bnWs.depthCallback == nil {
		return errors.New("please set depth callback func")
	}
	if size != 5 && size != 10 && size != 20 {
		return errors.New("please set depth size as 5 / 10 / 20")
	}
	endpoint := fmt.Sprintf("%s/%s@depth%d@100ms", bnWs.baseURL, strings.ToLower(pair.ToSymbol("")), size)

	handle := func(msg []byte) error {
		rawDepth := struct {
			LastUpdateID int64           `json:"lastUpdateId"`
			Bids         [][]interface{} `json:"bids"`
			Asks         [][]interface{} `json:"asks"`
		}{}
		err := json.Unmarshal(msg, &rawDepth)
		if err != nil {
			fmt.Println("json unmarshal error for ", string(msg))
			return err
		}
		depth := bnWs.parseDepthData(rawDepth.Bids, rawDepth.Asks)
		depth.Pair = pair
		depth.UTime = time.Now()
		bnWs.depthCallback(depth)
		return nil
	}
	bnWs.Subscribe(endpoint, handle)
	return nil
}

func (bnWs *BinanceWs) SubscribeTicker(pair CurrencyPair) error {
	if bnWs.tickerCallback == nil {
		return errors.New("please set ticker callback func")
	}
	endpoint := fmt.Sprintf("%s/%s@ticker", bnWs.baseURL, strings.ToLower(pair.ToSymbol("")))

	handle := func(msg []byte) error {
		datamap := make(map[string]interface{})
		err := json.Unmarshal(msg, &datamap)
		if err != nil {
			fmt.Println("json unmarshal error for ", string(msg))
			return err
		}

		msgType, isOk := datamap["e"].(string)
		if !isOk {
			return errors.New("no message type")
		}

		switch msgType {
		case "24hrTicker":
			tick := bnWs.parseTickerData(datamap)
			tick.Pair = pair
			bnWs.tickerCallback(tick)
			return nil
		default:
			return errors.New("unknown message " + msgType)
		}
	}
	bnWs.Subscribe(endpoint, handle)
	return nil
}

func (bnWs *BinanceWs) SubscribeTrade(pair CurrencyPair) error {
	if bnWs.tradeCallback == nil {
		return errors.New("please set trade callback func")
	}
	endpoint := fmt.Sprintf("%s/%s@trade", bnWs.baseURL, strings.ToLower(pair.ToSymbol("")))

	handle := func(msg []byte) error {
		datamap := make(map[string]interface{})
		err := json.Unmarshal(msg, &datamap)
		if err != nil {
			fmt.Println("json unmarshal error for ", string(msg))
			return err
		}

		msgType, isOk := datamap["e"].(string)
		if !isOk {
			return errors.New("no message type")
		}

		switch msgType {
		case "trade":
			side := BUY
			if datamap["m"].(bool) == false {
				side = SELL
			}
			trade := &RawTrade{
				Trade: Trade{
					Tid:    int64(ToUint64(datamap["t"])),
					Type:   TradeSide(side),
					Amount: ToFloat64(datamap["q"]),
					Price:  ToFloat64(datamap["p"]),
					Date:   int64(ToUint64(datamap["T"])),
				},
				BuyerOrderID:  ToInt64(datamap["b"]),
				SellerOrderID: ToInt64(datamap["a"]),
			}
			trade.Pair = pair
			bnWs.tradeCallback((*Trade)(unsafe.Pointer(trade)))
			return nil
		default:
			return errors.New("unknown message " + msgType)
		}
	}
	bnWs.Subscribe(endpoint, handle)
	return nil
}

func (bnWs *BinanceWs) SubscribeKline(pair CurrencyPair, period int) error {
	if bnWs.klineCallback == nil {
		return errors.New("place set kline callback func")
	}
	periodS, isOk := _INERNAL_KLINE_PERIOD_CONVERTER[period]
	if isOk != true {
		periodS = "M1"
	}
	endpoint := fmt.Sprintf("%s/%s@kline_%s", bnWs.baseURL, strings.ToLower(pair.ToSymbol("")), periodS)

	handle := func(msg []byte) error {
		datamap := make(map[string]interface{})
		err := json.Unmarshal(msg, &datamap)
		if err != nil {
			fmt.Println("json unmarshal error for ", string(msg))
			return err
		}

		msgType, isOk := datamap["e"].(string)
		if !isOk {
			return errors.New("no message type")
		}

		switch msgType {
		case "kline":
			k := datamap["k"].(map[string]interface{})
			period := _INERNAL_KLINE_PERIOD_REVERTER[k["i"].(string)]
			kline := bnWs.parseKlineData(k)
			kline.Pair = pair
			bnWs.klineCallback(kline, period)
			return nil
		default:
			return errors.New("unknown message " + msgType)
		}
	}
	bnWs.Subscribe(endpoint, handle)
	return nil
}

func (bnWs *BinanceWs) parseTickerData(tickmap map[string]interface{}) *Ticker {
	t := new(Ticker)
	t.Date = ToUint64(tickmap["E"])
	t.Last = ToFloat64(tickmap["c"])
	t.Vol = ToFloat64(tickmap["v"])
	t.Low = ToFloat64(tickmap["l"])
	t.High = ToFloat64(tickmap["h"])
	t.Buy = ToFloat64(tickmap["b"])
	t.Sell = ToFloat64(tickmap["a"])

	return t
}

func (bnWs *BinanceWs) parseDepthData(bids, asks [][]interface{}) *Depth {
	depth := new(Depth)
	for _, v := range bids {
		depth.BidList = append(depth.BidList, DepthRecord{ToFloat64(v[0]), ToFloat64(v[1])})
	}

	for _, v := range asks {
		depth.AskList = append(depth.AskList, DepthRecord{ToFloat64(v[0]), ToFloat64(v[1])})
	}
	return depth
}

func (bnWs *BinanceWs) parseKlineData(k map[string]interface{}) *Kline {
	kline := &Kline{
		Timestamp: int64(ToInt(k["t"])) / 1000,
		Open:      ToFloat64(k["o"]),
		Close:     ToFloat64(k["c"]),
		High:      ToFloat64(k["h"]),
		Low:       ToFloat64(k["l"]),
		Vol:       ToFloat64(k["v"]),
	}
	return kline
}

func (bnWs *BinanceWs) SubscribeAggTrade(pair CurrencyPair, tradeCallback func(*Trade)) error {
	if tradeCallback == nil {
		return errors.New("please set trade callback func")
	}
	endpoint := fmt.Sprintf("%s/%s@aggTrade", bnWs.baseURL, strings.ToLower(pair.ToSymbol("")))

	handle := func(msg []byte) error {
		datamap := make(map[string]interface{})
		err := json.Unmarshal(msg, &datamap)
		if err != nil {
			fmt.Println("json unmarshal error for ", string(msg))
			return err
		}

		msgType, isOk := datamap["e"].(string)
		if !isOk {
			return errors.New("no message type")
		}

		switch msgType {
		case "aggTrade":
			side := BUY
			if datamap["m"].(bool) == false {
				side = SELL
			}
			aggTrade := &AggTrade{
				Trade: Trade{
					Tid:    int64(ToUint64(datamap["a"])),
					Type:   TradeSide(side),
					Amount: ToFloat64(datamap["q"]),
					Price:  ToFloat64(datamap["p"]),
					Date:   int64(ToUint64(datamap["E"])),
				},
				FirstBreakdownTradeID: int64(ToUint64(datamap["f"])),
				LastBreakdownTradeID:  int64(ToUint64(datamap["l"])),
				TradeTime:             int64(ToUint64(datamap["T"])),
			}
			aggTrade.Pair = pair
			tradeCallback((*Trade)(unsafe.Pointer(aggTrade)))
			return nil
		default:
			return errors.New("unknown message " + msgType)
		}
	}
	bnWs.Subscribe(endpoint, handle)
	return nil
}

func (bnWs *BinanceWs) SubscribeDiffDepth(pair CurrencyPair, depthCallback func(*Depth)) error {
	if depthCallback == nil {
		return errors.New("please set depth callback func")
	}
	endpoint := fmt.Sprintf("%s/%s@depth", bnWs.baseURL, strings.ToLower(pair.ToSymbol("")))

	handle := func(msg []byte) error {
		rawDepth := struct {
			Type     string          `json:"e"`
			Time     int64           `json:"E"`
			Symbol   string          `json:"s"`
			UpdateID int             `json:"u"`
			Bids     [][]interface{} `json:"b"`
			Asks     [][]interface{} `json:"a"`
		}{}

		err := json.Unmarshal(msg, &rawDepth)
		if err != nil {
			fmt.Println("json unmarshal error for ", string(msg))
			return err
		}
		diffDepth := new(DiffDepth)
		for _, v := range rawDepth.Bids {
			diffDepth.BidList = append(diffDepth.BidList, DepthRecord{ToFloat64(v[0]), ToFloat64(v[1])})
		}

		for _, v := range rawDepth.Asks {
			diffDepth.AskList = append(diffDepth.AskList, DepthRecord{ToFloat64(v[0]), ToFloat64(v[1])})
		}

		diffDepth.Pair = pair
		diffDepth.UTime = time.Unix(0, rawDepth.Time*int64(time.Millisecond))
		depthCallback((*Depth)(unsafe.Pointer(diffDepth)))
		return nil
	}
	bnWs.Subscribe(endpoint, handle)
	return nil
}

// SubscribeExecutionReport executionReport
func (bnWs *BinanceWs) SubscribeExecutionReport(lKey string, depthCallback func(*Order)) error {
	if depthCallback == nil {
		return errors.New("please set depth callback func")
	}
	endpoint := fmt.Sprintf("%s/%s", bnWs.baseURL, lKey)
	handle := func(msg []byte) error {
		defer func() {
			if e := recover(); e != nil {
				log.Println("panic err: ", e)
			}
		}()
		if bytes.Contains(msg, []byte("executionReport")) {
			var order = new(Order)
			var data = map[string]interface{}{}
			_ = jsoniter.Unmarshal(msg, &data)
			if _, ok := data["e"]; ok {
				if price, ok := data["p"]; ok {
					order.Price = helper.StringToFloat64(price.(string))
				}
				if amount, ok := data["q"]; ok {
					order.Amount = helper.StringToFloat64(amount.(string))
				}
				if avgPrice, ok := data["p"]; ok {
					order.AvgPrice = helper.StringToFloat64(avgPrice.(string))
				}
				if dealAmount, ok := data["z"]; ok {
					order.DealAmount = helper.StringToFloat64(dealAmount.(string))
				}
				if Fee, ok := data["n"]; ok {
					order.Fee = helper.StringToFloat64(Fee.(string))
				}
				if Cid, ok := data["C"]; ok {
					order.Cid = Cid.(string)
				}
				if OrderID2, ok := data["i"]; ok {
					order.OrderID2 = fmt.Sprint(OrderID2)
				}
				//if OrderID2, ok := data["i"]; ok {
				//	order.OrderID = OrderID2.(int)
				//}
				if s, ok := data["s"]; ok {
					order.Symbol = s.(string)
				}
				if side, ok := data["S"]; ok {
					if side.(string) == "SELL" {
						order.Side = SELL
					} else {
						order.Side = BUY
					}
				}

				if o, ok := data["o"]; ok {
					order.Type = o.(string)
				}
				if o, ok := data["O"]; ok {
					order.OrderTime = ToInt(o)
				}
				if t, ok := data["T"]; ok {
					order.FinishedTime = ToInt64(t)
				}
				if status, ok := data["X"]; ok {
					switch status.(string) {
					case "NEW":
						order.Status = ORDER_UNFINISH
					case "FILLED":
						order.Status = ORDER_FINISH
					case "PARTIALLY_FILLED":
						order.Status = ORDER_PART_FINISH
					case "CANCELED":
						order.Status = ORDER_CANCEL
					case "PENDING_CANCEL":
						order.Status = ORDER_CANCEL_ING
					case "REJECTED":
						order.Status = ORDER_REJECT
					}
				}
			}

			depthCallback(order)
		}
		return nil
	}
	bnWs.Subscribe(endpoint, handle)
	return nil
}

func (bnWs *BinanceWs) exitHandler(c *WsConn) {
	pingTicker := time.NewTicker(10 * time.Minute)
	pongTicker := time.NewTicker(time.Second)
	defer pingTicker.Stop()
	defer pongTicker.Stop()
	defer c.CloseWs()

	for {
		select {
		case t := <-pingTicker.C:
			c.SendPingMessage([]byte(strconv.Itoa(int(t.UnixNano() / int64(time.Millisecond)))))
		case t := <-pongTicker.C:
			c.SendPongMessage([]byte(strconv.Itoa(int(t.UnixNano() / int64(time.Millisecond)))))
		}
	}
}
