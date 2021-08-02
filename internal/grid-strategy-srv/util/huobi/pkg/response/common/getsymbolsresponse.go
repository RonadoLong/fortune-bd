package common

import "github.com/shopspring/decimal"

type GetSymbolsResponse struct {
	Status string   `json:"status"`
	Data   []Symbol `json:"data"`
}

type Symbol struct {
	BaseCurrency    string          `json:"base-currency"`
	QuoteCurrency   string          `json:"quote-currency"`
	PricePrecision  int             `json:"price-precision"`
	AmountPrecision int             `json:"amount-precision"`
	SymbolPartition string          `json:"symbol-partition"`
	Symbol          string          `json:"symbol"`
	State           string          `json:"state"`
	ValuePrecision  int             `json:"value-precision"`
	MinOrderAmt     decimal.Decimal `json:"min-order-amt"`
	MaxOrderAmt     decimal.Decimal `json:"max-order-amt"`
	MinOrderValue   decimal.Decimal `json:"min-order-value"`
	LeverageRatio   decimal.Decimal `json:"leverage-ratio"`
}
