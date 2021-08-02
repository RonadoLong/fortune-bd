package account

import "github.com/shopspring/decimal"

type GetAccountLedgerResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    []Ledger `json:"data"`
	NextId  int64    `json:"nextId"`
}

type Ledger struct {
	AccountId    int64           `json:"accountId"`
	Currency     string          `json:"currency"`
	TransactAmt  decimal.Decimal `json:"transactAmt"`
	TransactType string          `json:"transactType"`
	TransferType string          `json:"transferType"`
	TransactId   int64           `json:"transactId"`
	TransactTime int64           `json:"transactTime"`
	Transferer   int64           `json:"transferer"`
	Transferee   int64           `json:"transferee"`
}
