package margin

type CrossMarginAccountsBalanceResponse struct {
	Status string                      `json:"status"`
	Data   *CrossMarginAccountsBalance `json:"data"`
}
type CrossMarginAccountsBalance struct {
	Id             int64  `json:"id"`
	Type           string `json:"type"`
	State          string `json:"state"`
	AcctBalanceSum string `json:"acct-balance-sum"`
	DebtBalanceSum string `json:"debt-balance-sum"`
	RiskRate       string `json:"risk-rate"`
	List           []struct {
		Currency string `json:"currency"`
		Type     string `json:"type"`
		Balance  string `json:"balance"`
	}
}
