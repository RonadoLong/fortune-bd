package order

type RequestOrderV1Response struct {
	Op        string `json:"op"`
	Timestamp int64  `json:"ts"`
	Topic     string `json:"topic"`
	ErrorCode string `json:"err-code"`
	ClientId  string `json:"cid"`
	Data      Order  `json:"data"`
}

type Order struct {
	AccountId        int    `json:"account-id"`
	Amount           string `json:"amount"`
	Id               int64  `json:"id"`
	Symbol           string `json:"symbol"`
	Price            string `json:"price"`
	CreatedAt        int64  `json:"created-at"`
	Type             string `json:"type"`
	FilledAmount     string `json:"filled-amount"`
	FilledCashAmount string `json:"filled-cash-amount"`
	FilledFees       string `json:"filled-fees"`
	Source           string `json:"source"`
	State            string `json:"state"`
}
