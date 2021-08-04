package market

import "github.com/shopspring/decimal"

type GetLatestTradeResponse struct {
	Status string     `json:"status"`
	Ch     string     `json:"ch"`
	Ts     int64      `json:"ts"`
	Tick   *TradeTick `json:"tick"`
}
type TradeTick struct {
	Id   int64 `json:"id"`
	Ts   int64 `json:"ts"`
	Data []struct {
		Amount    decimal.Decimal `json:"amount"`
		TradeId   int64           `json:"trade-id"`
		Ts        int64           `json:"ts"`
		Id        int64           `json:"id"`
		Price     decimal.Decimal `json:"price"`
		Direction string          `json:"direction"`
	}
}
