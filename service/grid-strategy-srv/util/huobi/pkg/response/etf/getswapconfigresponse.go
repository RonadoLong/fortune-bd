package etf

import "github.com/shopspring/decimal"

type GetSwapConfigResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    *SwapConfig `json:"data"`
}
type SwapConfig struct {
	PurchaseMinAmount   int64           `json:"purchase_min_amount"`
	PurchaseMaxAmount   int64           `json:"purchase_max_amount"`
	RedemptionMinAmount int64           `json:"redemption_min_amount"`
	RedemptionMaxAmount int64           `json:"redemption_max_amount"`
	PurchaseFeeRate     decimal.Decimal `json:"purchase_fee_rate"`
	RedemptionFeeRate   decimal.Decimal `json:"redemption_fee_rate"`
	EtfStatus           int             `json:"etf_status"`
	Unitprice           []*UnitPrice    `json:"unit_price"`
}
type UnitPrice struct {
	Currency string          `json:"currency"`
	Amount   decimal.Decimal `json:"amount"`
}
