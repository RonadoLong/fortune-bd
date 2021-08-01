package auth

import (
	"encoding/json"
	"wq-fotune-backend/service/grid-strategy-srv/util/huobi/pkg/response/base"
)

type WebSocketV2AuthenticationResponse struct {
	base.WebSocketV2ResponseBase
}

func ParseWSV2AuthResp(message []byte) *WebSocketV2AuthenticationResponse {
	result := &WebSocketV2AuthenticationResponse{}
	err := json.Unmarshal(message, result)
	if err != nil {
		return nil
	}

	return result
}
