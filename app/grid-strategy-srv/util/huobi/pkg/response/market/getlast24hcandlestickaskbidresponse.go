package market

import "github.com/shopspring/decimal"

type GetLast24hCandlestickAskBidResponse struct {
	Status string             `json:"status"`
	Ch     string             `json:"ch"`
	Ts     int64              `json:"ts"`
	Tick   *CandlestickAskBid `json:"tick"`
}
type CandlestickAskBid struct {
	Amount  decimal.Decimal   `json:"amount"`
	Open    decimal.Decimal   `json:"open"`
	Close   decimal.Decimal   `json:"close"`
	High    decimal.Decimal   `json:"high"`
	Id      int64             `json:"id"`
	Count   int64             `json:"count"`
	Low     decimal.Decimal   `json:"low"`
	Vol     decimal.Decimal   `json:"vol"`
	Version int64             `json:"version"`
	Bid     []decimal.Decimal `json:"bid"`
	Ask     []decimal.Decimal `json:"ask"`
}
