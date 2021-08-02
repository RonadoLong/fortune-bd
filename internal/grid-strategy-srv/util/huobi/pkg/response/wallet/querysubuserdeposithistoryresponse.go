package wallet

import "github.com/shopspring/decimal"

type QuerySubUserDepositHistoryResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    []DepositHistory `json:"data"`
	NextId  int64           `json:"nextId"`
}

type DepositHistory struct {
	Id              int64           `json:"id"`
	Currency        string          `json:"currency"`
	TransactionHash string          `json:"txHash"`
	Chain           string          `json:"chain"`
	Amount          decimal.Decimal `json:"amount"`
	Address         string          `json:"address"`
	AddressTag      string          `json:"addressTag"`
	State           string          `json:"state"`
	CreateTime      int64           `json:"createTime"`
	UpdateTime      int64           `json:"updateTime"`
}