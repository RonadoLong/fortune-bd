package getrequest

type GetAccountLedgerOptionalRequest struct {
	Currency string `json:"currency"`
	TransactTypes string `json:"transactTypes"`
	StartTime int64 `json:"startTime"`
	EndTime int64 `json:"endTime"`
	Sort string `json:"sort"`
	Limit int `json:"limit"`
	FromId int64 `json:"fromId"`
}