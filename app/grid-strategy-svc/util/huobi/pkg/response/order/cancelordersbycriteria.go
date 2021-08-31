package order

type CancelOrdersByCriteriaResponse struct {
	Status       string `json:"status"`
	ErrorCode    string `json:"err-code"`
	ErrorMessage string `json:"err-msg"`
	Data         *struct {
		SuccessCount int `json:"success-count"`
		FailedCount  int `json:"failed-count"`
		NextId       int `json:"next-id"`
	}
}
