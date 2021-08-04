package orderwebsocketclient

import (
	"encoding/json"
	"fmt"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/pkg/client/websocketclientbase"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/pkg/response/order"
)

// Responsible to handle order request from WebSocket
// This need authentication version 1
type RequestOrderWebSocketV1Client struct {
	websocketclientbase.WebSocketV1ClientBase
}

// Initializer
func (p *RequestOrderWebSocketV1Client) Init(accessKey string, secretKey string, host string) *RequestOrderWebSocketV1Client {
	p.WebSocketV1ClientBase.Init(accessKey, secretKey, host)
	return p
}

// Set callback service
func (p *RequestOrderWebSocketV1Client) SetHandler(
	authHandler websocketclientbase.AuthenticationV1ResponseHandler,
	responseHandler websocketclientbase.ResponseHandler) {
	p.WebSocketV1ClientBase.SetHandler(authHandler, p.handleMessage, responseHandler)
}

func (p *RequestOrderWebSocketV1Client) Request(orderId string, clientId string) error {

	req := fmt.Sprintf("{ \"op\":\"req\", \"topic\":\"orders.detail\", \"order-id\": \"%s\",\"cid\": \"%s\"}", orderId, clientId)
	return p.Send(req)
}

func (p *RequestOrderWebSocketV1Client) handleMessage(msg string) (interface{}, error) {
	result := order.RequestOrderV1Response{}
	err := json.Unmarshal([]byte(msg), &result)
	return result, err
}
