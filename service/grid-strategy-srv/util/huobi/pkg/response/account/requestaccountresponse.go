package account

type RequestAccountV1Response struct {
	Timestamp int64            `json:"ts"`
	Op        string           `json:"op"`
	Topic     string           `json:"topic"`
	ErrorCode int              `json:"err-code"`
	ClientId  string           `json:"cid"`
	Data      []AccountBalance `json:"data"`
}

type AccountBalance struct {
	Id    int       `json:"id"`
	Type  string    `json:"type"`
	State string    `json:"state"`
	List  []Balance `json:"list"`
}

type Balance struct {
	Currency string `json:"currency"`
	Type     string `json:"type"`
	Balance  string `json:"balance"`
}
