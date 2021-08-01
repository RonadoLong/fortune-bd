package market

import (
	"wq-fotune-backend/service/grid-strategy-srv/util/huobi/pkg/response/base"
)

type SubscribeLast24hCandlestickResponse struct {
	base.WebSocketResponseBase
	Data *Candlestick
	Tick *Candlestick
}
