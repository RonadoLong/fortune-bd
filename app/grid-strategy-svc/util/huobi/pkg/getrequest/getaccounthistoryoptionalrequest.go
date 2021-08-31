package getrequest

type GetAccountHistoryOptionalRequest struct {
	Currency      string
	TransactTypes string
	StartTime     int64
	EndTime       int64
	Sort          string
	Size          int
}
