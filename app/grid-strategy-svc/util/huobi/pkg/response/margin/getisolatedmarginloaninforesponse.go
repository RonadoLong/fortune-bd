package margin

type GetIsolatedMarginLoanInfoResponse struct {
	Status string                   `json:"status"`
	Data   []IsolatedMarginLoanInfo `json:"data"`
}
type IsolatedMarginLoanInfo struct {
	Symbol     string `json:"symbol"`
	Currencies []struct {
		Currency     string `json:"currency"`
		InterestRate string `json:"interest-rate"`
		MinLoanAmt   string `json:"min-loan-amt"`
		MaxLoanAmt   string `json:"max-loan-amt"`
		LoanableAmt  string `json:"loanable-amt"`
		ActualRate   string `json:"actual-rate"`
	}
}
