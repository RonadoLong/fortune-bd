package order

type CancelOrderByClientResponse struct {
	Status       string `json:"status"`
	Data         int    `json:"data"`
	ErrorCode    string `json:"err-code"`
	ErrorMessage string `json:"err-msg"`
}
