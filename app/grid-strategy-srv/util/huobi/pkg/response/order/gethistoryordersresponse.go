package order

type GetHistoryOrdersResponse struct {
	Status       string `json:"status"`
	ErrorCode    string `json:"err-code"`
	ErrorMessage string `json:"err-msg"`
	Data         []struct {
		Id               int64  `json:"id"`
		ClientOrderId    string `json:"client-order-id"`
		AccountId        int    `json:"account-id"`
		UserId           int    `json:"user-id"`
		Amount           string `json:"amount"`
		Symbol           string `json:"symbol"`
		Price            string `json:"price"`
		CreatedAt        int64  `json:"created-at"`
		CanceledAt       int64  `json:"canceled-at"`
		FinishedAt       int64  `json:"finished-at"`
		Type             string `json:"type"`
		FilledAmount     string `json:"field-amount"`
		FilledCashAmount string `json:"field-cash-amount"`
		FilledFees       string `json:"field-fees"`
		Source           string `json:"source"`
		State            string `json:"state"`
		Exchange         string `json:"exchange"`
		Batch            string `json:"batch"`
		StopPrice        string `json:"stop-price"`
		Operator         string `json:"operator"`
	}
}
