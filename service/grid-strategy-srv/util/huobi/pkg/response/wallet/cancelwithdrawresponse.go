package wallet

type CancelWithdrawResponse struct {
	Status string `json:"status"`
	Data   int64  `json:"data"`
}
