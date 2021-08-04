package market

type GetHistoricalTradeResponse struct {
	Status string      `json:"status"`
	Ch     string      `json:"ch"`
	Ts     int64       `json:"ts"`
	Data   []TradeTick `json:"data"`
}
