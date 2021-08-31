package cron

import "encoding/json"

var RateKey = "rate:usd-rmb"

type Ticker struct {
	Symbol string  `json:"symbol"`
	Last   float64 `json:"last"`
	Buy    float64 `json:"buy"`
	Open   float64 `json:"open"`
	Sell   float64 `json:"sell"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Vol    float64 `json:"vol"`
	Change string  `json:"change"`
	Date   uint64  `json:"date"` // 单位:ms
}

func (tk Ticker) MarshalBinary() ([]byte, error) {
	return json.Marshal(tk)
}


type QuoteRate struct {
	InstrumentID string `json:"instrument_id"`
	Rate         string `json:"rate"`
	Timestamp    string `json:"timestamp"`
}

func (rate QuoteRate) MarshalBinary() ([]byte, error) {
	return json.Marshal(rate)
}
