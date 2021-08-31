package postrequest

import "github.com/shopspring/decimal"

type FuturesTransferRequest struct {
	Currency string          `json:"currency"`
	Amount   decimal.Decimal `json:"amount"`
	Type     string          `json:"type"`
}
