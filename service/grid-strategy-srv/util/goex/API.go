package goex

// api interface

type API interface {
	LimitBuy(amount, price string, currency CurrencyPair) (*Order, error)
	LimitSell(amount, price string, currency CurrencyPair) (*Order, error)
	MarketBuy(amount, price string, currency CurrencyPair) (*Order, error)
	MarketSell(amount, price string, currency CurrencyPair) (*Order, error)
	CancelOrder(orderId string, currency CurrencyPair) (bool, error)
	GetOneOrder(orderId string, currency CurrencyPair) (*Order, error)
	GetUnfinishOrders(currency CurrencyPair) ([]Order, error)
	GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error)
	GetAccount() (*Account, error)

	GetTicker(currency CurrencyPair) (*Ticker, error)
	GetDepth(size int, currency CurrencyPair) (*Depth, error)
	GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error)
	//非个人，整个交易所的交易记录
	GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error)

	GetExchangeName() string
}

// Accounter 各个交易所账号公共接口
type Accounter interface {
	GetCurrencyBalance(currency string) (float64, error)
	PlaceLimitOrder(side string, symbol string, price string, amount string, clientOrderID string) (string, error)
	PlaceMarketOrder(side string, symbol string, amount string, clientOrderID string) (string, error)
	CancelOrder(orderID string, symbol string) error
	GetOrderInfo(orderID string, symbol string) (interface{}, error)
	GetHistoryOrdersInfo(symbol string, states string, types string) (interface{}, error)
}

// Processer 处理网格接口
type Processer interface {
	UpdateGridOrder(tradeTime int64, tradeAmount float64, orderID, clientOrderID, orderStatus string) error
}
