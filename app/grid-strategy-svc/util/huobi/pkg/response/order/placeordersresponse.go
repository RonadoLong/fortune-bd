package order

type PlaceOrdersResponse struct {
	Status       string `json:"status"`
	ErrorCode    string `json:"err-code"`
	ErrorMessage string `json:"err-msg"`
	Data         []PlaceOrderResult
}

type PlaceOrderResult struct {
	OrderId       int64  `json:"order-id"`
	ClientOrderId string `json:"client-order-id"`
	ErrorCode     string `json:"err-code"`
	ErrorMessage  string `json:"err-msg"`
}
