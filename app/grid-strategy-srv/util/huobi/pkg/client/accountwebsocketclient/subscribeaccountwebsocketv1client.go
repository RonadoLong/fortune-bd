package accountwebsocketclient

import (
	"encoding/json"
	"fmt"

	"wq-grid-strategy/util/huobi/pkg/client/websocketclientbase"
	"wq-grid-strategy/util/huobi/pkg/response/account"
)

// Responsible to handle account asset subscription from WebSocket
// This need authentication version 1
type SubscribeAccountWebSocketV1Client struct {
	websocketclientbase.WebSocketV1ClientBase
}

// Initializer
func (p *SubscribeAccountWebSocketV1Client) Init(accessKey string, secretKey string, host string) *SubscribeAccountWebSocketV1Client {
	p.WebSocketV1ClientBase.Init(accessKey, secretKey, host)
	return p
}

// Set callback biz
func (p *SubscribeAccountWebSocketV1Client) SetHandler(
	authHandler websocketclientbase.AuthenticationV1ResponseHandler,
	responseHandler websocketclientbase.ResponseHandler) {
	p.WebSocketV1ClientBase.SetHandler(authHandler, p.handleMessage, responseHandler)
}

// Subscribe all balance updates of the current account
// 1 to include frozen balance
// 0 to not
func (p *SubscribeAccountWebSocketV1Client) Subscribe(mode string, clientId string) error {

	sub := fmt.Sprintf("{ \"op\":\"sub\", \"topic\":\"accounts\", \"mode\": \"%s\", \"cid\": \"%s\"}", mode, clientId)
	return p.Send(sub)
}

// Unsubscribe balance updates
func (p *SubscribeAccountWebSocketV1Client) UnSubscribe(mode string, clientId string) error {
	unsub := fmt.Sprintf("{ \"op\":\"unsub\", \"topic\":\"accounts\", \"mode\": \"%s\", \"cid\": \"%s\" }", mode, clientId)
	return p.Send(unsub)
}
func (p *SubscribeAccountWebSocketV1Client) handleMessage(msg string) (interface{}, error) {
	result := account.SubscribeAccountV1Response{}
	err := json.Unmarshal([]byte(msg), &result)
	return result, err
}
