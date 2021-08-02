package order

type RequestOrdersV1Response struct {
	Op        string  `json:"op"`
	Timestamp int64   `json:"ts"`
	Topic     string  `json:"topic"`
	ErrorCode string  `json:"err-code"`
	ClientId  string  `json:"cid"`
	Data      []Order `json:"data"`
}
