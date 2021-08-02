package postrequest

type SwapRequest struct {
	EtfName string `json:"etf_name"`
	Amount  int64  `json:"amount"`
}
