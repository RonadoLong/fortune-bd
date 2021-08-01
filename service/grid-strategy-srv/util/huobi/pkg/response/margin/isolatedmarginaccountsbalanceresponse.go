package margin

type IsolatedMarginAccountsBalanceResponse struct {
	Status string                          `json:"status"`
	Data   []IsolatedMarginAccountsBalance `json:"data"`
}
type IsolatedMarginAccountsBalance struct {
	Id       int64  `json:"id"`
	Type     string `json:"type"`
	State    string `json:"state"`
	Symbol   string `json:"symbol"`
	FlPrice  string `json:"fl-price"`
	FlType   string `json:"fl-type"`
	RiskRate string `json:"risk-rate"`
	List     []struct {
		Currency string `json:"currency"`
		Type     string `json:"type"`
		Balance  string `json:"balance"`
	}
}
