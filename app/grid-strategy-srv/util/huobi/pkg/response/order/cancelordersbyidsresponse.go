package order

type CancelOrdersByIdsResponse struct {
	Status       string `json:"status"`
	ErrorCode    string `json:"err-code"`
	ErrorMessage string `json:"err-msg"`
	Data         *struct {
		Success []string `json:"success"`
		Failed  []struct {
			OrderId       string `json:"order-id"`
			ClientOrderId string `json:"client-order-id"`
			ErrorCode     string `json:"err-code"`
			ErrorMessage  string `json:"err-msg"`
		}
	}
}
