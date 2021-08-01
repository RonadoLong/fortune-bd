package margin

type CrossMarginLoanOrdersResponse struct {
	Status string                 `json:"status"`
	Data   []CrossMarginLoanOrder `json:"data"`
}
type CrossMarginLoanOrder struct {
	Id              int64  `json:"id"`
	UserId          int64  `json:"user-id"`
	AccountId       int64  `json:"account-id"`
	Currency        string `json:"currency"`
	LoanAmount      string `json:"loan-amount"`
	LoanBalance     string `json:"loan-balance"`
	InterestAmount  string `json:"interest-amount"`
	InterestBalance string `json:"interest-balance"`
	CreatedAt       int64  `json:"created-at"`
	AccruedAt       int64  `json:"accrued-at"`
	State           string `json:"state"`
	FilledPoints    string `json:"filled-points"`
	FilledHt        string `json:"filled-ht"`
}
