package order

type SubscribeOrderV1Response struct {
	Op        string `json:"op"`
	Timestamp int64  `json:"ts"`
	Topic     string `json:"topic"`
	Data      struct {
		MatchId          int    `json:"match-id"`
		OrderId          int64  `json:"order-id"`
		Symbol           string `json:"symbol"`
		OrderState       string `json:"order-state"`
		Role             string `json:"role"`
		Price            string `json:"price"`
		FilledAmount     string `json:"filled-amount"`
		FilledCashAmount string `json:"filled-cash-amount"`
		UnfilledAmount   string `json:"unfilled-amount"`
		ClientOrderId    string `json:"client-order-id"`
		OrderType        string `json:"order-type"`
	}
}
