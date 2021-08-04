package account

type GetAccountInfoResponse struct {
	Status string        `json:"status"`
	Data   []AccountInfo `json:"data"`
}
type AccountInfo struct {
	Id      int64  `json:"id"`
	Type    string `json:"type"`
	Subtype string `json:"subtype"`
	State   string `json:"state"`
}
