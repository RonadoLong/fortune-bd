package order

type GetOrderResponse struct {
	Status       string `json:"status"`
	ErrorCode    string `json:"err-code"`
	ErrorMessage string `json:"err-msg"`
	Data         *struct {
		Id               int64  `json:"id"`
		ClientOrderId    string `json:"client-order-id"`
		Symbol           string `json:"symbol"`
		AccountId        int    `json:"account-id"`
		Amount           string `json:"amount"`
		Price            string `json:"price"`
		CreatedAt        int64  `json:"created-at"`
		Type             string `json:"type"`
		FilledAmount     string `json:"field-amount"`
		FilledCashAmount string `json:"field-cash-amount"`
		FilledFees       string `json:"field-fees"`
		Source           string `json:"source"`
		State            string `json:"state"`
		FinishedAt       int64  `json:"finished-at"`

		StopPrice string `json:"stop-price"` //止盈止损订单触发价格
		Operator  string `json:"operator"`   // 止盈止损订单触发价运算符	gte,lte
	}
}
