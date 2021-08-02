package market

import "github.com/shopspring/decimal"

type GetCandlestickResponse struct {
	Status string        `json:"status"`
	Ch     string        `json:"ch"`
	Ts     int64         `json:"ts"`
	Data   []Candlestick `json:"data"`
}
type Candlestick struct {
	Amount decimal.Decimal `json:"amount"`
	Open   decimal.Decimal `json:"open"`
	Close  decimal.Decimal `json:"close"`
	High   decimal.Decimal `json:"high"`
	Id     int64           `json:"id"`
	Count  int64           `json:"count"`
	Low    decimal.Decimal `json:"low"`
	Vol    decimal.Decimal `json:"vol"`
}
