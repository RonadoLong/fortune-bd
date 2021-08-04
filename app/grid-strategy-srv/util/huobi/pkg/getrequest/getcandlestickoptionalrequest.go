package getrequest

const (
	MIN1  = "1min"
	MIN5  = "5min"
	MIN15 = "15min"
	MIN30 = "30min"
	MIN60 = "60min"
	HOUR4 = "4hour"
	DAY1  = "1day"
	MON1  = "1mon"
	WEEK1 = "1week"
	YEAR1 = "1year"
)

type GetCandlestickOptionalRequest struct {
	Period string
	Size   int
}
