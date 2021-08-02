package account

import "wq-fotune-backend/internal/grid-strategy-srv/util/huobi/pkg/response/base"

type SubscribeAccountV2Response struct {
	base.WebSocketV2ResponseBase
	Data *struct {
		Currency    string `json:"currency"`
		AccountId   int    `json:"accountId"`
		Balance     string `json:"balance"`
		ChangeType  string `json:"changeType"`
		AccountType string `json:"accountType"`
		ChangeTime  int64  `json:"changeTime"`
	}
}
