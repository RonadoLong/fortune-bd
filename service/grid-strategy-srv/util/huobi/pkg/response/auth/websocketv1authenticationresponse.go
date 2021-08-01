package auth

import "encoding/json"

type WebSocketV1AuthenticationResponse struct {
	Op        string `json:"op"`
	Timestamp int64  `json:"ts"`
	ErrorCode int    `json:"err-code"`
	Data      *struct {
		UserId int `json:"user-id"`
	}
}

func (p *WebSocketV1AuthenticationResponse) IsAuth() bool {
	return p.Op == "auth" && p.ErrorCode == 0
}

func ParseWSV1AuthResp(message string) *WebSocketV1AuthenticationResponse {
	result := &WebSocketV1AuthenticationResponse{}
	err := json.Unmarshal([]byte(message), result)
	if err != nil {
		return nil
	}

	return result
}
