package etf

import "github.com/shopspring/decimal"

type GetSwapListResponse struct {
	Code int         `json:"code"`
	Data []*SwapList `json:"data"`
}
type SwapList struct {
	Id         int64           `json:"id"`
	GmtCreated int64           `json:"gmt_created"`
	Currency   string          `json:"currency"`
	Amount     decimal.Decimal `json:"amount"`
	Type       int             `json:"type"`
	Status     int             `json:"status"`
	Detail     *struct {
		UsedCurrencyList   []*UnitPrice    `json:"used_ currency_list"`
		Rate               decimal.Decimal `json:"rate"`
		Fee                decimal.Decimal `json:"fee"`
		PointCardAmount    decimal.Decimal `json:"point_card_amount"`
		ObtainCurrencyList []*UnitPrice    `json:"obtain_currency_list"`
	}
}
