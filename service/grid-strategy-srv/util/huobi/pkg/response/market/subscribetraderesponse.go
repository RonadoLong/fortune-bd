package market

import (

	"github.com/shopspring/decimal"
	"wq-fotune-backend/service/grid-strategy-srv/util/huobi/pkg/response/base"
)

type SubscribeTradeResponse struct {
	base.WebSocketResponseBase
	Data []Trade
	Tick *struct {
		Id        int64 `json:"id"`
		Timestamp int64 `json:"ts"`
		Data      []Trade
	}
}

type Trade struct {
	TradeId   int64           `json:"tradeId"`
	Amount    decimal.Decimal `json:"amount"`
	Price     decimal.Decimal `json:"price"`
	Timestamp int64           `json:"ts"`
	Direction string          `json:"direction"`
}
