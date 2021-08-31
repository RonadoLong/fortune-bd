package auth

import (
	"encoding/json"
	"fortune-bd/app/grid-strategy-svc/util/huobi/pkg/response/base"
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
