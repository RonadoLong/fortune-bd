package order

type CancelOrderByIdResponse struct {
	Status       string `json:"status"`
	Data         string `json:"data"`
	ErrorCode    string `json:"err-code"`
	ErrorMessage string `json:"err-msg"`
}
