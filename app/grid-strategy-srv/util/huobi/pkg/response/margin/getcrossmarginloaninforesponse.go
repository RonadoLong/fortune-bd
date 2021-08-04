package margin

type GetCrossMarginLoanInfoResponse struct {
	Status string                `json:"status"`
	Data   []CrossMarginLoanInfo `json:"data"`
}
type CrossMarginLoanInfo struct {
	Currency     string `json:"currency"`
	InterestRate string `json:"interest-rate"`
	MinLoanAmt   string `json:"min-loan-amt"`
	MaxLoanAmt   string `json:"max-loan-amt"`
	LoanableAmt  string `json:"loanable-amt"`
	ActualRate   string `json:"actual-rate"`
}
