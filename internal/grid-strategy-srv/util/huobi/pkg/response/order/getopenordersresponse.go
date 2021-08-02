package order

import "github.com/shopspring/decimal"

type GetOpenOrdersResponse struct {
	Status       string `json:"status"`
	ErrorCode    string `json:"err-code"`
	ErrorMessage string `json:"err-msg"`
	Data         []OpenOrder
}

type OpenOrder struct {
	Id               int64           `json:"id"`
	ClientOrderId    string          `json:"client-order-id"`
	AccountId        int             `json:"account-id"`
	Amount           decimal.Decimal `json:"amount"`
	Symbol           string          `json:"symbol"`
	Price            decimal.Decimal `json:"price"`
	CreatedAt        int64           `json:"created-at"`
	Type             string          `json:"type"`
	FilledAmount     decimal.Decimal `json:"filled-amount"`
	FilledCashAmount decimal.Decimal `json:"filled-cash-amount"`
	FilledFees       decimal.Decimal `json:"filled-fees"`
	Source           string          `json:"source"`
	State            string          `json:"state"`
	StopPrice        decimal.Decimal `json:"stop-price"`
	Operator         string          `json:"operator"`
}
