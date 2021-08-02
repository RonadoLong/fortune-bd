package account

type GetAccountHistoryResponse struct {
	Status string           `json:"status"`
	Data   []AccountHistory `json:"data"`
}
type AccountHistory struct {
	AccountId    int64  `json:"account-id"`
	Currency     string `json:"currency"`
	TransactAmt  string `json:"transact-amt"`
	TransactType string `json:"transact-type"`
	RecordId     int64  `json:"record-id"`
	AvailBalance string `json:"avail-balance"`
	AcctBalance  string `json:"acct-balance"`
	TransactTime int64  `json:"transact-time"`
}
