package market

type GetLast24hCandlestick struct {
	Status string       `json:"status"`
	Ch     string       `json:"ch"`
	Ts     int64        `json:"ts"`
	Tick   *Candlestick `json:"tick"`
}
