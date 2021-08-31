package binance

import (
	"errors"
	"fmt"
	"fortune-bd/libs/env"
	"log"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	. "fortune-bd/app/grid-strategy-svc/util/goex"
	"github.com/json-iterator/go"
	"github.com/shopspring/decimal"
)

const (
	GLOBAL_API_BASE_URL = "https://api.binance.com"
	US_API_BASE_URL     = "https://api.binance.us"
	JE_API_BASE_URL     = "https://api.binance.je"
	//API_V1       = API_BASE_URL + "api/v1/"
	//API_V3       = API_BASE_URL + "api/v3/"

	TICKER_URI             = "ticker/24hr?symbol=%s"
	TICKERS_URI            = "ticker/allBookTickers"
	DEPTH_URI              = "depth?symbol=%s&limit=%d"
	ACCOUNT_URI            = "account?"
	ORDER_URI              = "order"
	UNFINISHED_ORDERS_INFO = "openOrders?"
	KLINE_URI              = "klines"
	SERVER_TIME_URL        = "time"

	UserDataStream = "userDataStream?"

	// OrderStateSubmitted 已委托
	OrderStateSubmitted = "submitted"
	// OrderStateSubmitted 已取消
	OrderStateCanceled = "canceled"
	// OrderStateFilled 已成交
	OrderStateFilled = "filled"
	// FilledFees 永续合约已成交的手续费
	FilledFees = 0.001
	// BrokerID 经纪商id
	BrokerID = "x-H1DQZX35"
)

var _INERNAL_KLINE_PERIOD_CONVERTER = map[int]string{
	KLINE_PERIOD_1MIN:   "1m",
	KLINE_PERIOD_3MIN:   "3m",
	KLINE_PERIOD_5MIN:   "5m",
	KLINE_PERIOD_15MIN:  "15m",
	KLINE_PERIOD_30MIN:  "30m",
	KLINE_PERIOD_60MIN:  "1h",
	KLINE_PERIOD_1H:     "1h",
	KLINE_PERIOD_2H:     "2h",
	KLINE_PERIOD_4H:     "4h",
	KLINE_PERIOD_6H:     "6h",
	KLINE_PERIOD_8H:     "8h",
	KLINE_PERIOD_12H:    "12h",
	KLINE_PERIOD_1DAY:   "1d",
	KLINE_PERIOD_3DAY:   "3d",
	KLINE_PERIOD_1WEEK:  "1w",
	KLINE_PERIOD_1MONTH: "1M",
}

type Filter struct {
	FilterType          string  `json:"filterType"`
	MaxPrice            float64 `json:"maxPrice,string"`
	MinPrice            float64 `json:"minPrice,string"`
	TickSize            float64 `json:"tickSize,string"`
	MultiplierUp        float64 `json:"multiplierUp,string"`
	MultiplierDown      float64 `json:"multiplierDown,string"`
	AvgPriceMins        int     `json:"avgPriceMins"`
	MinQty              float64 `json:"minQty,string"`
	MaxQty              float64 `json:"maxQty,string"`
	StepSize            float64 `json:"stepSize,string"`
	MinNotional         float64 `json:"minNotional,string"`
	ApplyToMarket       bool    `json:"applyToMarket"`
	Limit               int     `json:"limit"`
	MaxNumAlgoOrders    int     `json:"maxNumAlgoOrders"`
	MaxNumIcebergOrders int     `json:"maxNumIcebergOrders"`
	MaxNumOrders        int     `json:"maxNumOrders"`
}

type RateLimit struct {
	Interval      string `json:"interval"`
	IntervalNum   int64  `json:"intervalNum"`
	Limit         int64  `json:"limit"`
	RateLimitType string `json:"rateLimitType"`
}

type TradeSymbol struct {
	Symbol                     string   `json:"symbol"`
	Status                     string   `json:"status"`
	BaseAsset                  string   `json:"baseAsset"`
	BaseAssetPrecision         int      `json:"baseAssetPrecision"`
	QuoteAsset                 string   `json:"quoteAsset"`
	QuotePrecision             int      `json:"quotePrecision"`
	BaseCommissionPrecision    int      `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision   int      `json:"quoteCommissionPrecision"`
	Filters                    []Filter `json:"filters"`
	IcebergAllowed             bool     `json:"icebergAllowed"`
	IsMarginTradingAllowed     bool     `json:"isMarginTradingAllowed"`
	IsSpotTradingAllowed       bool     `json:"isSpotTradingAllowed"`
	OcoAllowed                 bool     `json:"ocoAllowed"`
	QuoteOrderQtyMarketAllowed bool     `json:"quoteOrderQtyMarketAllowed"`
	OrderTypes                 []string `json:"orderTypes"`
}

func (ts TradeSymbol) GetMinAmount() float64 {
	for _, v := range ts.Filters {
		if v.FilterType == "LOT_SIZE" {
			return v.MinQty
		}
	}
	return 0
}

func (ts TradeSymbol) GetAmountPrecision() int {
	for _, v := range ts.Filters {
		if v.FilterType == "LOT_SIZE" {
			step := strconv.FormatFloat(v.StepSize, 'f', -1, 64)
			pres := strings.Split(step, ".")
			if len(pres) == 1 {
				return 0
			}
			return len(pres[1])
		}
	}
	return 0
}

func (ts TradeSymbol) GetMinPrice() float64 {
	for _, v := range ts.Filters {
		if v.FilterType == "PRICE_FILTER" {
			return v.MinPrice
		}
	}
	return 0
}

func (ts TradeSymbol) GetMinValue() float64 {
	for _, v := range ts.Filters {
		if v.FilterType == "MIN_NOTIONAL" {
			return v.MinNotional
		}
	}
	return 0
}

func (ts TradeSymbol) GetPricePrecision() int {
	for _, v := range ts.Filters {
		if v.FilterType == "PRICE_FILTER" {
			step := strconv.FormatFloat(v.TickSize, 'f', -1, 64)
			pres := strings.Split(step, ".")
			if len(pres) == 1 {
				return 0
			}
			return len(pres[1])
		}
	}
	return 0
}

type ExchangeInfo struct {
	Timezone        string        `json:"timezone"`
	ServerTime      int           `json:"serverTime"`
	ExchangeFilters []interface{} `json:"exchangeFilters,omitempty"`
	RateLimits      []RateLimit   `json:"rateLimits"`
	Symbols         []TradeSymbol `json:"symbols"`
}

type Binance struct {
	accessKey  string
	secretKey  string
	baseUrl    string
	apiV1      string
	apiV3      string
	httpClient *http.Client
	timeOffset int64 //nanosecond
	*ExchangeInfo
}

func (bn *Binance) buildParamsSigned(postForm *url.Values) error {
	postForm.Set("recvWindow", "60000")
	tonce := strconv.FormatInt(time.Now().UnixNano()+bn.timeOffset, 10)[0:13]
	postForm.Set("timestamp", tonce)
	payload := postForm.Encode()
	sign, _ := GetParamHmacSHA256Sign(bn.secretKey, payload)
	postForm.Set("signature", sign)
	return nil
}

func New(client *http.Client, api_key, secret_key string) *Binance {
	return NewWithConfig(&APIConfig{
		HttpClient:   client,
		Endpoint:     GLOBAL_API_BASE_URL,
		ApiKey:       api_key,
		ApiSecretKey: secret_key})
}

func NewWithConfig(config *APIConfig) *Binance {
	if config.Endpoint == "" {
		config.Endpoint = GLOBAL_API_BASE_URL
	}

	bn := &Binance{
		baseUrl:    config.Endpoint,
		apiV1:      config.Endpoint + "/api/v1/",
		apiV3:      config.Endpoint + "/api/v3/",
		accessKey:  config.ApiKey,
		secretKey:  config.ApiSecretKey,
		httpClient: config.HttpClient}
	bn.setTimeOffset()
	return bn
}

func (bn *Binance) GetExchangeName() string {
	return BINANCE
}

func (bn *Binance) Ping() bool {
	_, err := HttpGet(bn.httpClient, bn.apiV3+"ping")
	if err != nil {
		return false
	}
	return true
}

func (bn *Binance) setTimeOffset() error {
	respmap, err := HttpGet(bn.httpClient, bn.apiV3+SERVER_TIME_URL)
	if err != nil {
		return err
	}

	stime := int64(ToInt(respmap["serverTime"]))
	st := time.Unix(stime/1000, 1000000*(stime%1000))
	lt := time.Now()
	offset := st.Sub(lt).Nanoseconds()
	bn.timeOffset = int64(offset)
	return nil
}

func (bn *Binance) GetTicker(currency CurrencyPair) (*Ticker, error) {
	tickerUri := bn.apiV3 + fmt.Sprintf(TICKER_URI, currency.ToSymbol(""))
	tickerMap, err := HttpGet(bn.httpClient, tickerUri)

	if err != nil {
		return nil, err
	}

	var ticker Ticker
	ticker.Pair = currency
	t, _ := tickerMap["closeTime"].(float64)
	ticker.Date = uint64(t / 1000)
	ticker.Last = ToFloat64(tickerMap["lastPrice"])
	ticker.Buy = ToFloat64(tickerMap["bidPrice"])
	ticker.Sell = ToFloat64(tickerMap["askPrice"])
	ticker.Low = ToFloat64(tickerMap["lowPrice"])
	ticker.High = ToFloat64(tickerMap["highPrice"])
	ticker.Vol = ToFloat64(tickerMap["volume"])
	return &ticker, nil
}

func (bn *Binance) GetDepth(size int, currencyPair CurrencyPair) (*Depth, error) {
	if size <= 5 {
		size = 5
	} else if size <= 10 {
		size = 10
	} else if size <= 20 {
		size = 20
	} else if size <= 50 {
		size = 50
	} else if size <= 100 {
		size = 100
	} else if size <= 500 {
		size = 500
	} else {
		size = 1000
	}

	apiUrl := fmt.Sprintf(bn.apiV3+DEPTH_URI, currencyPair.ToSymbol(""), size)
	resp, err := HttpGet(bn.httpClient, apiUrl)
	if err != nil {
		return nil, err
	}

	if _, isok := resp["code"]; isok {
		return nil, errors.New(resp["msg"].(string))
	}

	bids := resp["bids"].([]interface{})
	asks := resp["asks"].([]interface{})

	depth := new(Depth)
	depth.Pair = currencyPair
	depth.UTime = time.Now()
	n := 0
	for _, bid := range bids {
		_bid := bid.([]interface{})
		amount := ToFloat64(_bid[1])
		price := ToFloat64(_bid[0])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.BidList = append(depth.BidList, dr)
		n++
		if n == size {
			break
		}
	}

	n = 0
	for _, ask := range asks {
		_ask := ask.([]interface{})
		amount := ToFloat64(_ask[1])
		price := ToFloat64(_ask[0])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.AskList = append(depth.AskList, dr)
		n++
		if n == size {
			break
		}
	}

	return depth, nil
}

func (bn *Binance) placeOrder(amount, price string, pair CurrencyPair, orderType, orderSide, newClientOrderId string) (*Order, error) {
	path := bn.apiV3 + ORDER_URI
	params := url.Values{}
	params.Set("symbol", pair.ToSymbol(""))
	params.Set("side", orderSide)
	params.Set("type", orderType)
	params.Set("newOrderRespType", "ACK")
	params.Set("quantity", amount)
	params.Set("newClientOrderId", newClientOrderId)

	switch orderType {
	case "LIMIT":
		params.Set("timeInForce", "GTC")
		params.Set("price", price)
	case "MARKET":
		params.Set("newOrderRespType", "RESULT")
	}

	bn.buildParamsSigned(&params)

	resp, err := HttpPostForm2(bn.httpClient, path, params,
		map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return nil, err
	}

	respmap := make(map[string]interface{})
	err = jsoniter.Unmarshal(resp, &respmap)
	if err != nil {
		return nil, err
	}

	orderId := ToInt(respmap["orderId"])
	if orderId <= 0 {
		return nil, errors.New(string(resp))
	}

	side := BUY
	if orderSide == "SELL" {
		side = SELL
	}

	dealAmount := ToFloat64(respmap["executedQty"])
	cummulativeQuoteQty := ToFloat64(respmap["cummulativeQuoteQty"])
	avgPrice := 0.0
	if cummulativeQuoteQty > 0 && dealAmount > 0 {
		avgPrice = cummulativeQuoteQty / dealAmount
	}

	return &Order{
		Currency:   pair,
		OrderID:    orderId,
		OrderID2:   strconv.Itoa(orderId),
		Price:      ToFloat64(price),
		Amount:     ToFloat64(amount),
		DealAmount: dealAmount,
		AvgPrice:   avgPrice,
		Side:       TradeSide(side),
		Status:     ORDER_UNFINISH,
		OrderTime:  ToInt(respmap["transactTime"])}, nil
}

func (bn *Binance) GetAccount() (*Account, error) {
	params := url.Values{}
	bn.buildParamsSigned(&params)
	path := bn.apiV3 + ACCOUNT_URI + params.Encode()
	respmap, err := HttpGet2(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return nil, err
	}
	if _, isok := respmap["code"]; isok == true {
		return nil, errors.New(respmap["msg"].(string))
	}
	acc := Account{}
	acc.Exchange = bn.GetExchangeName()
	acc.SubAccounts = make(map[Currency]SubAccount)

	balances := respmap["balances"].([]interface{})
	for _, v := range balances {
		vv := v.(map[string]interface{})
		currency := NewCurrency(vv["asset"].(string), "").AdaptBccToBch()
		balance, _ := decimal.NewFromFloat(ToFloat64(vv["free"])).Add(decimal.NewFromFloat(ToFloat64(vv["locked"]))).Float64()
		acc.SubAccounts[currency] = SubAccount{
			Currency:     currency,
			Amount:       ToFloat64(vv["free"]),
			ForzenAmount: ToFloat64(vv["locked"]),
			Balance:      balance,
		}
	}

	return &acc, nil
}

func (bn *Binance) LimitBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return bn.placeOrder(amount, price, currencyPair, "LIMIT", "BUY", "")
}

func (bn *Binance) LimitSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return bn.placeOrder(amount, price, currencyPair, "LIMIT", "SELL", "")
}

func (bn *Binance) MarketBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return bn.placeOrder(amount, price, currencyPair, "MARKET", "BUY", "")
}

func (bn *Binance) MarketSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return bn.placeOrder(amount, price, currencyPair, "MARKET", "SELL", "")
}

func (bn *Binance) CancelOrder(orderId string, currencyPair CurrencyPair) (bool, error) {
	path := bn.apiV3 + ORDER_URI
	params := url.Values{}
	params.Set("symbol", currencyPair.ToSymbol(""))
	params.Set("orderId", orderId)

	bn.buildParamsSigned(&params)

	resp, err := HttpDeleteForm(bn.httpClient, path, params, map[string]string{"X-MBX-APIKEY": bn.accessKey})

	if err != nil {
		return false, err
	}

	respmap := make(map[string]interface{})
	err = jsoniter.Unmarshal(resp, &respmap)
	if err != nil {
		return false, err
	}

	orderIdCanceled := ToInt(respmap["orderId"])
	if orderIdCanceled <= 0 {
		return false, errors.New(string(resp))
	}

	return true, nil
}

func (bn *Binance) GetOneOrder(orderId string, currencyPair CurrencyPair) (*Order, error) {
	params := url.Values{}
	params.Set("symbol", currencyPair.ToSymbol(""))
	if orderId != "" {
		params.Set("orderId", orderId)
	}
	params.Set("orderId", orderId)

	bn.buildParamsSigned(&params)
	path := bn.apiV3 + ORDER_URI + "?" + params.Encode()

	respmap, err := HttpGet2(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return nil, err
	}

	status := respmap["status"].(string)
	side := respmap["side"].(string)

	ord := Order{}
	ord.Currency = currencyPair
	ord.OrderID = ToInt(orderId)
	ord.OrderID2 = orderId
	ord.Cid, _ = respmap["clientOrderId"].(string)
	ord.Type = respmap["type"].(string)

	if side == "SELL" {
		ord.Side = SELL
	} else {
		ord.Side = BUY
	}

	switch status {
	case "NEW":
		ord.Status = ORDER_UNFINISH
	case "FILLED":
		ord.Status = ORDER_FINISH
	case "PARTIALLY_FILLED":
		ord.Status = ORDER_PART_FINISH
	case "CANCELED":
		ord.Status = ORDER_CANCEL
	case "PENDING_CANCEL":
		ord.Status = ORDER_CANCEL_ING
	case "REJECTED":
		ord.Status = ORDER_REJECT
	}

	ord.Amount = ToFloat64(respmap["origQty"].(string))
	ord.Price = ToFloat64(respmap["price"].(string))
	ord.DealAmount = ToFloat64(respmap["executedQty"])
	ord.AvgPrice = ord.Price // response no avg price ， fill price
	ord.OrderTime = ToInt(respmap["time"])

	cummulativeQuoteQty := ToFloat64(respmap["cummulativeQuoteQty"])
	if cummulativeQuoteQty > 0 {
		ord.AvgPrice = cummulativeQuoteQty / ord.DealAmount
	}

	return &ord, nil
}

func (bn *Binance) GetUnfinishOrders(currencyPair CurrencyPair) ([]Order, error) {
	params := url.Values{}
	params.Set("symbol", currencyPair.ToSymbol(""))

	bn.buildParamsSigned(&params)
	path := bn.apiV3 + UNFINISHED_ORDERS_INFO + params.Encode()

	respmap, err := HttpGet3(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return nil, err
	}

	orders := make([]Order, 0)
	for _, v := range respmap {
		ord := v.(map[string]interface{})
		side := ord["side"].(string)
		orderSide := SELL
		if side == "BUY" {
			orderSide = BUY
		}
		ordId := ToInt(ord["orderId"])
		orders = append(orders, Order{
			OrderID:   ordId,
			OrderID2:  strconv.Itoa(ordId),
			Currency:  currencyPair,
			Price:     ToFloat64(ord["price"]),
			Amount:    ToFloat64(ord["origQty"]),
			Side:      TradeSide(orderSide),
			Status:    ORDER_UNFINISH,
			OrderTime: ToInt(ord["time"])})
	}
	return orders, nil
}

func (bn *Binance) GetAllUnfinishOrders() ([]Order, error) {
	params := url.Values{}

	bn.buildParamsSigned(&params)
	path := bn.apiV3 + UNFINISHED_ORDERS_INFO + params.Encode()

	respmap, err := HttpGet3(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return nil, err
	}

	orders := make([]Order, 0)
	for _, v := range respmap {
		ord := v.(map[string]interface{})
		side := ord["side"].(string)
		orderSide := SELL
		if side == "BUY" {
			orderSide = BUY
		}

		ordId := ToInt(ord["orderId"])
		orders = append(orders, Order{
			OrderID:   ToInt(ord["orderId"]),
			OrderID2:  strconv.Itoa(ordId),
			Currency:  bn.toCurrencyPair(ord["symbol"].(string)),
			Price:     ToFloat64(ord["price"]),
			Amount:    ToFloat64(ord["origQty"]),
			Side:      TradeSide(orderSide),
			Status:    ORDER_UNFINISH,
			OrderTime: ToInt(ord["time"])})
	}
	return orders, nil
}

func (bn *Binance) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	params := url.Values{}
	params.Set("symbol", currency.ToSymbol(""))
	params.Set("interval", _INERNAL_KLINE_PERIOD_CONVERTER[period])
	if since > 0 {
		params.Set("startTime", strconv.Itoa(since))
	}
	//params.Set("endTime", strconv.Itoa(int(time.Now().UnixNano()/1000000)))
	params.Set("limit", fmt.Sprintf("%d", size))

	klineUrl := bn.apiV3 + KLINE_URI + "?" + params.Encode()
	klines, err := HttpGet3(bn.httpClient, klineUrl, nil)
	if err != nil {
		return nil, err
	}
	var klineRecords []Kline

	for _, _record := range klines {
		r := Kline{Pair: currency}
		record := _record.([]interface{})
		r.Timestamp = int64(record[0].(float64)) / 1000 //to unix timestramp
		r.Open = ToFloat64(record[1])
		r.High = ToFloat64(record[2])
		r.Low = ToFloat64(record[3])
		r.Close = ToFloat64(record[4])
		r.Vol = ToFloat64(record[5])

		klineRecords = append(klineRecords, r)
	}

	return klineRecords, nil

}

//非个人，整个交易所的交易记录
//注意：since is fromId
func (bn *Binance) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	param := url.Values{}
	param.Set("symbol", currencyPair.ToSymbol(""))
	param.Set("limit", "500")
	if since > 0 {
		param.Set("fromId", strconv.Itoa(int(since)))
	}
	apiUrl := bn.apiV3 + "historicalTrades?" + param.Encode()
	resp, err := HttpGet3(bn.httpClient, apiUrl, map[string]string{
		"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return nil, err
	}

	var trades []Trade
	for _, v := range resp {
		m := v.(map[string]interface{})
		ty := SELL
		if m["isBuyerMaker"].(bool) {
			ty = BUY
		}
		trades = append(trades, Trade{
			Tid:    ToInt64(m["id"]),
			Type:   ty,
			Amount: ToFloat64(m["qty"]),
			Price:  ToFloat64(m["price"]),
			Date:   ToInt64(m["time"]),
			Pair:   currencyPair,
		})
	}

	return trades, nil
}

func (bn *Binance) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	params := url.Values{}
	params.Set("symbol", currency.ToSymbol(""))

	bn.buildParamsSigned(&params)
	path := bn.apiV3 + "allOrders?" + params.Encode()

	respmap, err := HttpGet3(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return nil, err
	}

	orders := make([]Order, 0)
	for _, v := range respmap {
		ord := v.(map[string]interface{})
		side := ord["side"].(string)
		orderSide := SELL
		if side == "BUY" {
			orderSide = BUY
		}
		ordId := ToInt(ord["orderId"])
		orders = append(orders, Order{
			OrderID:   ToInt(ord["orderId"]),
			OrderID2:  strconv.Itoa(ordId),
			Currency:  currency,
			Price:     ToFloat64(ord["price"]),
			Amount:    ToFloat64(ord["origQty"]),
			Side:      TradeSide(orderSide),
			Status:    ORDER_UNFINISH,
			OrderTime: ToInt(ord["time"])})
	}
	return orders, nil

}

// OrderInfo 订单信息
type OrderInfo struct {
	ID               string // 订单id
	ClientOrderID    string // 用户自定义订单id
	Exchange         string // 交易所
	Symbol           string // 品种
	Price            string // 价格
	Amount           string // 交易数量
	Type             string // 订单类型，limit:限价单, market:市价单
	Side             string // 订单方向，buy:买入，sell:卖出
	FilledAmount     string // 已成交数量
	FilledCashAmount string // 已成交额
	FilledFees       string // 已成交手续费
	AnchorSymbol     string // 锚定币，例如usdt
	CreatedAt        int64  // 创建订单时间时间戳，精确到毫秒
	UpdateAt         int64  // 修改订单状态时间戳
	State            string // 状态
	StopPrice        string // 止盈止损
	Operator         string // 操作符
}

func convertState(state string) string {
	switch state {
	case "NEW":
		return "submitted"
	case "CANCELED":
		return "canceled"
	case "FILLED":
		return "filled"
	case "PARTIALLY_FILLED", "PENDING_CANCEL", "REJECTED":
		return strings.ToLower(state)
	}

	return state
}

// GetOrderInfo 获取订单信息
func (bn *Binance) GetOrderInfo(orderId string, currencyPair CurrencyPair) (*OrderInfo, error) {
	params := url.Values{}
	params.Set("symbol", currencyPair.ToSymbol(""))
	if orderId != "" {
		params.Set("orderId", orderId)
	}
	params.Set("orderId", orderId)

	bn.buildParamsSigned(&params)
	path := bn.apiV3 + ORDER_URI + "?" + params.Encode()

	oi, err := HttpGet2(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return nil, err
	}

	return &OrderInfo{
		ID:            strconv.Itoa(ToInt(oi["orderId"])),
		ClientOrderID: ToString(oi["clientOrderId"]),
		Exchange:      "binance",
		Symbol:        strings.ToLower(ToString(oi["symbol"])),
		//Price:            fmt.Sprintf("%v", ToFloat64(oi["price"])),
		Price:            strconv.FormatFloat(ToFloat64(oi["price"]), 'f', -1, 64),
		Amount:           fmt.Sprintf("%v", ToFloat64(oi["origQty"])),
		Type:             strings.ToLower(ToString(oi["type"])),
		Side:             strings.ToLower(ToString(oi["side"])),
		FilledAmount:     ToString(oi["executedQty"]),
		FilledCashAmount: ToString(oi["cumQuote"]),
		CreatedAt:        int64(ToInt(oi["time"])),
		UpdateAt:         int64(ToInt(oi["updateTime"])),
		State:            convertState(ToString(oi["status"])),
	}, nil
}

// GetSubmittedOrders 获取品种的全部委托挂单
func (bn *Binance) GetSubmittedOrders(currencyPair CurrencyPair) ([]*OrderInfo, error) {
	params := url.Values{}
	params.Set("symbol", currencyPair.ToSymbol(""))

	bn.buildParamsSigned(&params)
	path := bn.apiV3 + UNFINISHED_ORDERS_INFO + params.Encode()

	respmap, err := HttpGet3(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return nil, err
	}

	orders := []*OrderInfo{}
	for _, v := range respmap {
		oi := v.(map[string]interface{})
		orders = append(orders, &OrderInfo{
			ID:            strconv.Itoa(ToInt(oi["orderId"])),
			ClientOrderID: ToString(oi["clientOrderId"]),
			Exchange:      "binance",
			Symbol:        strings.ToLower(ToString(oi["symbol"])),
			//Price:            fmt.Sprintf("%v", ToFloat64(oi["price"])),
			Price:            strconv.FormatFloat(ToFloat64(oi["price"]), 'f', -1, 64),
			Amount:           fmt.Sprintf("%v", ToFloat64(oi["origQty"])),
			Type:             strings.ToLower(ToString(oi["type"])),
			Side:             strings.ToLower(ToString(oi["side"])),
			FilledAmount:     ToString(oi["executedQty"]),
			FilledCashAmount: ToString(oi["cumQuote"]),
			CreatedAt:        int64(ToInt(oi["time"])),
			UpdateAt:         int64(ToInt(oi["updateTime"])),
			State:            convertState(ToString(oi["status"])),
		})
	}

	return orders, nil
}

// Get30DaysOrders 获取最近30天的已成交和委托挂单
func (bn *Binance) Get30DaysOrders(currency CurrencyPair, currentPage, pageSize int) ([]*OrderInfo, error) {
	params := url.Values{}
	params.Set("symbol", currency.ToSymbol(""))

	bn.buildParamsSigned(&params)
	path := bn.apiV3 + "allOrders?" + params.Encode()

	respmap, err := HttpGet3(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return nil, err
	}

	orders := []*OrderInfo{}
	for _, v := range respmap {
		oi := v.(map[string]interface{})
		orders = append(orders, &OrderInfo{
			ID:            strconv.Itoa(ToInt(oi["orderId"])),
			ClientOrderID: ToString(oi["clientOrderId"]),
			Exchange:      "binance",
			Symbol:        strings.ToLower(ToString(oi["symbol"])),
			//Price:            fmt.Sprintf("%v", ToFloat64(oi["price"])),
			Price:            strconv.FormatFloat(ToFloat64(oi["price"]), 'f', -1, 64),
			Amount:           fmt.Sprintf("%v", ToFloat64(oi["origQty"])),
			Type:             strings.ToLower(ToString(oi["type"])),
			Side:             strings.ToLower(ToString(oi["side"])),
			FilledAmount:     ToString(oi["executedQty"]),
			FilledCashAmount: ToString(oi["cumQuote"]),
			CreatedAt:        int64(ToInt(oi["time"])),
			UpdateAt:         int64(ToInt(oi["updateTime"])),
			State:            convertState(ToString(oi["status"])),
		})
	}

	return orders, nil
}

func (bn *Binance) GetUserDataStream() (string, error) {
	params := url.Values{}
	_ = bn.buildParamsSigned(&params)
	path := bn.apiV3 + "userDataStream"
	respmap, err := HttpPostForm2(bn.httpClient, path, nil, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return "", err
	}
	var resp map[string]interface{}
	_ = jsoniter.Unmarshal(respmap, &resp)
	if _, isok := resp["code"]; isok == true {
		return "", errors.New(resp["msg"].(string))
	}
	return resp["listenKey"].(string), nil
}

func (bn *Binance) PutUserDataStream(listenKey string) (bool, error) {
	path := bn.apiV3 + "userDataStream?" + "listenKey=" + listenKey
	respmap, err := HttpPutData(bn.httpClient, path, "", map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return false, err
	}
	var resp map[string]interface{}
	_ = jsoniter.Unmarshal(respmap, &resp)
	if _, isok := resp["code"]; isok == true {
		return false, errors.New(resp["msg"].(string))
	}
	return true, nil
}

func (bn *Binance) toCurrencyPair(symbol string) CurrencyPair {
	if bn.ExchangeInfo == nil {
		var err error
		bn.ExchangeInfo, err = bn.GetExchangeInfo()
		if err != nil {
			return CurrencyPair{}
		}
	}
	for _, v := range bn.ExchangeInfo.Symbols {
		if v.Symbol == symbol {
			return NewCurrencyPair2(v.BaseAsset + "_" + v.QuoteAsset)
		}
	}
	return CurrencyPair{}
}

func (bn *Binance) GetExchangeInfo() (*ExchangeInfo, error) {
	resp, err := HttpGet5(bn.httpClient, bn.apiV3+"exchangeInfo", nil)
	if err != nil {
		return nil, err
	}
	info := &ExchangeInfo{}
	err = jsoniter.Unmarshal(resp, info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (bn *Binance) GetTradeSymbol(currencyPair CurrencyPair) (*TradeSymbol, error) {
	if bn.ExchangeInfo == nil {
		var err error
		bn.ExchangeInfo, err = bn.GetExchangeInfo()
		if err != nil {
			return nil, err
		}
	}
	for k, v := range bn.ExchangeInfo.Symbols {
		if v.Symbol == currencyPair.ToSymbol("") {
			return &bn.ExchangeInfo.Symbols[k], nil
		}
	}
	return nil, errors.New("symbol not found")
}

func (bn *Binance) CancelAllOrders(currencyPair CurrencyPair) (bool, error) {
	params := url.Values{}
	params.Set("symbol", currencyPair.ToSymbol(""))
	params.Set("timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	_ = bn.buildParamsSigned(&params)

	path := bn.apiV3 + "openOrders"
	respmap, err := HttpDeleteForm(bn.httpClient, path, params, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return false, err
	}
	var resp map[string]interface{}
	_ = jsoniter.Unmarshal(respmap, &resp)
	if _, isok := resp["code"]; isok == true {
		return false, errors.New(resp["msg"].(string))
	}
	log.Printf("%+v", resp)
	return true, nil
}

// -------------------------------------------------------------------------------------------------

type BinanceAccount struct {
	Binance *Binance
}

// InitAccount 实例化
func InitAccount(accessKey, secretKey string, proxyAddr ...string) (*BinanceAccount, error) {
	api := &APIConfig{
		HttpClient:   &http.Client{Timeout: 10 * time.Second},
		Endpoint:     GLOBAL_API_BASE_URL,
		ApiKey:       accessKey,
		ApiSecretKey: secretKey,
	}
	api.HttpClient.Transport = getTransport(proxyAddr...)

	bn := &Binance{
		baseUrl:    api.Endpoint,
		apiV1:      api.Endpoint + "/api/v1/",
		apiV3:      api.Endpoint + "/api/v3/",
		accessKey:  api.ApiKey,
		secretKey:  api.ApiSecretKey,
		httpClient: api.HttpClient}
	bn.setTimeOffset()

	return &BinanceAccount{bn}, nil
}

// GetAccountBalance 获取账号余额，返回map，如果key不存在表示余额为零，其中key由"货币:类型"组成，类型： trade表示可使用余额，frozen表示冻结余额
func (b *BinanceAccount) GetAccountBalance() (map[string]float64, error) {
	balances := map[string]float64{}

	account, err := b.Binance.GetAccount()
	if err != nil {
		return balances, err
	}

	for _, v := range account.SubAccounts {
		if v.Amount > 0.0 {
			balances[strings.ToLower(v.Currency.Symbol)+":trade"] = v.Amount
		}
		if v.ForzenAmount > 0.0 {
			balances[strings.ToLower(v.Currency.Symbol)+":frozen"] = v.ForzenAmount
		}
	}

	return balances, nil
}

// GetCurrencyBalance 查询某个货币的可用余额
func (b *BinanceAccount) GetCurrencyBalance(currency string) (float64, error) {
	balances, err := b.GetAccountBalance()
	if err != nil {
		return 0, err
	}

	currency = strings.ToLower(currency)
	val, ok := balances[currency+":trade"]
	if !ok {
		return 0, nil
	}

	return val, nil
}

func getCurrencyPair(symbol string) CurrencyPair {
	cpCache := GetCurrencyPair(symbol)
	if cpCache.CurrencyB.Symbol != "" {
		return *cpCache
	}

	symbol = strings.ToUpper(symbol)
	cp, ok := CurrencyPairMap[symbol]
	if !ok {
		return UNKNOWN_PAIR
	}

	return cp
}

// PlaceLimitOrder 买入、卖出限价单
func (b *BinanceAccount) PlaceLimitOrder(side string, symbol string, price string, amount string, clientOrderID string) (string, error) {
	order, err := &Order{}, error(nil)
	cp := getCurrencyPair(symbol)

	side = strings.ToUpper(side)
	if side != "BUY" && side != "SELL" {
		return "", fmt.Errorf("unknown side=%s", side)
	}

	order, err = b.Binance.placeOrder(amount, price, cp, "LIMIT", side, clientOrderID)
	if err != nil {
		return "", err
	}

	return order.OrderID2, nil
}

// PlaceMarketOrder 买入、卖出市价单
func (b *BinanceAccount) PlaceMarketOrder(side string, symbol string, amount string, clientOrderID string) (string, error) {
	order, err := &Order{}, error(nil)
	cp := getCurrencyPair(symbol)

	side = strings.ToUpper(side)
	if side != "BUY" && side != "SELL" {
		return "", fmt.Errorf("unknown side=%s", side)
	}

	order, err = b.Binance.placeOrder(amount, "", cp, "MARKET", side, clientOrderID)
	if err != nil {
		return "", err
	}

	return order.OrderID2, nil
}

// CancelOrder 取消委托订单
func (b *BinanceAccount) CancelOrder(orderID string, symbol string) error {
	cp := getCurrencyPair(symbol)
	ok, err := b.Binance.CancelOrder(orderID, cp)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("cancel order failed")
	}

	return nil
}

// GetHistoryOrdersInfo 获取历史订单列表
func (b *BinanceAccount) GetOrderInfo(orderID string, symbol string) (interface{}, error) {
	cp := getCurrencyPair(symbol)
	//return b.Binance.GetOneOrder(orderID, cp)
	return b.Binance.GetOrderInfo(orderID, cp)
}

// GetHistoryOrdersInfo 获取历史订单列表，所有挂单、最近30天的成交订单、取消订单，参数states值有submitted:已挂单, filled:已成交 canceled:已取消
func (b *BinanceAccount) GetHistoryOrdersInfo(symbol string, states string, types string) (interface{}, error) {
	cp := getCurrencyPair(symbol)

	switch states {
	case "submitted": // 挂单
		return b.Binance.GetSubmittedOrders(cp)

	case "filled", "canceled": // 已成交
		orderInfos, err := b.Binance.Get30DaysOrders(cp, 0, 20)
		if err != nil {
			return nil, err
		}

		orders := []*OrderInfo{}
		for _, v := range orderInfos {
			if v.State == states {
				orders = append(orders, v)
			}
		}
		return orders, nil

	default:
		return nil, fmt.Errorf("unknown states %s", states)
	}
}

// -------------------------------------------------------------------------------------------------

// SymbolLimit 品种限制信息
type SymbolLimit struct {
	Symbol            string  // 品种
	QuantityLimit     float64 // 买卖最小交易所数量
	VolumeLimit       float64 // 买卖最小交易额(单位是品种后面的货币)
	PricePrecision    int     // 价格精度
	QuantityPrecision int     // 数量精度
}

// GetSelectSymbols 获取指定品种信息
func GetSelectSymbols(symbols []string, proxyAddr ...string) ([]*SymbolLimit, error) {
	limitSymbols := make([]*SymbolLimit, 0)
	b := NewWithConfig(&APIConfig{
		HttpClient: &http.Client{Transport: getTransport(proxyAddr...), Timeout: 10 * time.Second},
		Endpoint:   GLOBAL_API_BASE_URL,
	})

	if b.ExchangeInfo == nil {
		var err error
		b.ExchangeInfo, err = b.GetExchangeInfo()
		if err != nil {
			return nil, err
		}
	}

	for _, v := range b.ExchangeInfo.Symbols {
		symbol := strings.ToLower(v.Symbol)
		for _, selectSymbol := range symbols {
			if selectSymbol == symbol {
				limitSymbols = append(limitSymbols, &SymbolLimit{
					Symbol:            symbol,
					QuantityLimit:     v.GetMinAmount(),
					VolumeLimit:       v.GetMinValue(),
					PricePrecision:    v.GetPricePrecision(),
					QuantityPrecision: v.GetAmountPrecision(),
				})
			}
		}
	}

	return limitSymbols, nil
}

// GetAnchorCurrencySymbols 获取某个币本位下的所有品种信息
func GetAnchorCurrencySymbols(anchorCurrency string, proxyAddr ...string) ([]*SymbolLimit, error) {
	limitSymbols := make([]*SymbolLimit, 0)
	if anchorCurrency == "" {
		return nil, errors.New("anchorCurrency is empty")
	}

	b := NewWithConfig(&APIConfig{
		HttpClient: &http.Client{Transport: getTransport(proxyAddr...), Timeout: 10 * time.Second},
		Endpoint:   GLOBAL_API_BASE_URL,
	})

	if b.ExchangeInfo == nil {
		var err error
		b.ExchangeInfo, err = b.GetExchangeInfo()
		if err != nil {
			return nil, err
		}
	}

	acLen := len(anchorCurrency)
	for _, v := range b.ExchangeInfo.Symbols {
		symbol := strings.ToLower(v.Symbol)
		if len(symbol) > acLen {
			if symbol[len(symbol)-acLen:] == anchorCurrency {
				limitSymbols = append(limitSymbols, &SymbolLimit{
					Symbol:            symbol,
					QuantityLimit:     v.GetMinAmount(),
					VolumeLimit:       v.GetMinValue(),
					PricePrecision:    v.GetPricePrecision(),
					QuantityPrecision: v.GetAmountPrecision(),
				})
			}
		}
	}

	return limitSymbols, nil
}

// GetSymbolLimit 获取品种限制值
func GetSymbolLimit(symbol string, proxyAddr ...string) (*SymbolLimit, error) {
	b := NewWithConfig(&APIConfig{
		HttpClient: &http.Client{Transport: getTransport(proxyAddr...), Timeout: 10 * time.Second},
		Endpoint:   GLOBAL_API_BASE_URL,
	})

	ts, err := b.GetTradeSymbol(getCurrencyPair(symbol))
	if err != nil {
		return &SymbolLimit{}, err
	}

	return &SymbolLimit{
		Symbol:            strings.ToLower(symbol),
		QuantityLimit:     ts.GetMinAmount(),
		VolumeLimit:       ts.GetMinValue(),
		PricePrecision:    ts.GetPricePrecision(),
		QuantityPrecision: ts.GetAmountPrecision(),
	}, nil
}

// GetLatestPrice 获取品种的最新价格
func GetLatestPrice(symbol string, proxyAddr ...string) (float64, error) {
	b := NewWithConfig(&APIConfig{
		HttpClient: env.GetProxyHttpClient(),
		Endpoint:   GLOBAL_API_BASE_URL,
	})

	tk, err := b.GetTicker(getCurrencyPair(symbol))
	if err != nil {
		return 0, err
	}

	if tk.Last <= 0.0 {
		return 0.0, fmt.Errorf("price is zero, price = %v", tk.Last)
	}

	return tk.Last, nil
}

// GetKlines 获取最近一个周期的K线数据
func GetKlines(symbol string, interval int, size int, proxyAddr ...string) ([]Kline, error) {
	b := NewWithConfig(&APIConfig{
		HttpClient: &http.Client{Transport: getTransport(proxyAddr...), Timeout: 10 * time.Second},
		Endpoint:   GLOBAL_API_BASE_URL,
	})

	klines, err := b.GetKlineRecords(getCurrencyPair(symbol), interval, size, 0)
	if err != nil {
		return nil, err
	}

	return klines, nil
}

// Get3CyclePriceRange 获取最近3个周期最低价格
func Get3CyclePriceRange(symbol string, interval int, proxyAddr ...string) (float64, float64, error) {
	klines, err := GetKlines(symbol, interval, 3, proxyAddr...)
	if err != nil {
		return 0, 0, err
	}

	lowerPrices := []float64{}
	highPrices := []float64{}
	for _, kl := range klines {
		lowerPrices = append(lowerPrices, kl.Low)
		highPrices = append(highPrices, kl.High)
	}
	sort.Float64s(lowerPrices)
	sort.Float64s(highPrices)

	lowerPrice, highPrice := 0.0, 0.0
	lowLen := len(lowerPrices)
	if lowLen > 1 {
		lowerPrice = (lowerPrices[0] + lowerPrices[1]) / 2
	} else if lowLen == 1 {
		lowerPrice = lowerPrices[0]
	} else {
		return 0.0, 0.0, errors.New("lowPrices is nil")
	}

	highLen := len(highPrices)
	if highLen == 3 {
		highPrice = (highPrices[1] + highPrices[2]) / 2
	} else if highLen == 2 {
		highPrice = (highPrices[0] + highPrices[1]) / 2
	} else if highLen == 1 {
		highPrice = highPrices[0]
	} else {
		return 0.0, 0.0, errors.New("highPrices is nil")
	}

	return lowerPrice, highPrice, nil
}

func getTransport(proxyAddr ...string) *http.Transport {
	transport := &http.Transport{
		Dial: (&net.Dialer{Timeout: 10 * time.Second}).Dial,
	}
	if len(proxyAddr) > 0 && proxyAddr[0] != "" {
		transport = &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) { return url.Parse(proxyAddr[0]) },
			Dial:  (&net.Dialer{Timeout: 10 * time.Second}).Dial,
		}
	}

	return transport
}

func Str2Int64(s string) int64 {
	if s == "" {
		return 0
	}
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return value
}

// 字符串转float64，非法字符串的返回值默认为0.0
func Str2Float64(str string) float64 {
	if str == "" {
		return 0
	}

	str = strings.Replace(str, " ", "", -1)
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Printf("ParseFloat error, err=%s\n", err.Error())
		return 0.0
	}

	return f
}
