package market

import (
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/pkg/response/base"
)

type SubscribeDepthResponse struct {
	base.WebSocketResponseBase
	Data *Depth
	Tick *Depth
}
