package getrequest

type QuerySubUserDepositHistoryOptionalRequest struct {
	Currency  string `json:"currency"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
	Sort      string `json:"sort"`
	Limit     string `json:"limit"`
	FromId    int64  `json:"fromId"`
}
