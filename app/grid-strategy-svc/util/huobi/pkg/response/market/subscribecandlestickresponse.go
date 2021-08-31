package market

import (
	"fortune-bd/app/grid-strategy-svc/util/huobi/pkg/response/base"
	"github.com/shopspring/decimal"
)

type SubscribeCandlestickResponse struct {
	base base.WebSocketResponseBase
	Tick *Tick
	Data []Tick
}
type Tick struct {
	Id     int64           `json:"id"`
	Amount decimal.Decimal `json:"amount"`
	Count  int             `json:"count"`
	Open   decimal.Decimal `json:"open"`
	Close  decimal.Decimal `json:"close"`
	Low    decimal.Decimal `json:"low"`
	High   decimal.Decimal `json:"high"`
	Vol    decimal.Decimal `json:"vol"`
}
