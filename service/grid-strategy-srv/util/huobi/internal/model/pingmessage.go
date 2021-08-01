package model

import "encoding/json"

type PingMessage struct {
	Ping int64 `json:"ping"`
}

func ParsePingMessage(message string) *PingMessage {
	result := PingMessage{}
	err := json.Unmarshal([]byte(message), &result)
	if err != nil {
		return nil
	}

	return &result
}
