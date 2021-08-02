package postrequest

type CreateWithdrawRequest struct {
	Address  string `json:"address"`
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
	Fee      string `json:"fee"`
	Chain    string `json:"chain,omitempty"`
	AddrTag  string `json:"addr-tag,omitempty"`
}
