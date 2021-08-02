package base

type WebSocketResponseBase struct {
	Status    string `json:"status"`
	Channel   string `json:"ch"`
	Timestamp int64  `json:"ts"`
}
