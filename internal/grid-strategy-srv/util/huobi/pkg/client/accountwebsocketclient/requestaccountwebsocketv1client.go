package accountwebsocketclient

import (
	"encoding/json"
	"fmt"

	"wq-grid-strategy/util/huobi/pkg/client/websocketclientbase"
	"wq-grid-strategy/util/huobi/pkg/response/account"
)

// Responsible to handle account asset request from WebSocket
// This need authentication version 1
type RequestAccountWebSocketV1Client struct {
	websocketclientbase.WebSocketV1ClientBase
}

// Initializer
func (p *RequestAccountWebSocketV1Client) Init(accessKey string, secretKey string, host string) *RequestAccountWebSocketV1Client {
	p.WebSocketV1ClientBase.Init(accessKey, secretKey, host)
	return p
}

// Set callback handler
func (p *RequestAccountWebSocketV1Client) SetHandler(
	authHandler websocketclientbase.AuthenticationV1ResponseHandler,
	responseHandler websocketclientbase.ResponseHandler) {
	p.WebSocketV1ClientBase.SetHandler(authHandler, p.handleMessage, responseHandler)
}

// Request all account data of the current user.
func (p *RequestAccountWebSocketV1Client) Request(clientId string) error {

	req := fmt.Sprintf("{ \"op\":\"req\", \"topic\":\"accounts.list\", \"cid\": \"%s\"}", clientId)
	return p.Send(req)
}

func (p *RequestAccountWebSocketV1Client) handleMessage(msg string) (interface{}, error) {
	result := account.RequestAccountV1Response{}
	err := json.Unmarshal([]byte(msg), &result)
	return result, err
}
