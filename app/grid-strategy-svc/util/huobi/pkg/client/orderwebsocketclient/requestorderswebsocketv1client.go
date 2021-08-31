package orderwebsocketclient

import (
	"encoding/json"
	"fortune-bd/app/grid-strategy-svc/util/huobi/pkg/client/websocketclientbase"
	"fortune-bd/app/grid-strategy-svc/util/huobi/pkg/getrequest"
	"fortune-bd/app/grid-strategy-svc/util/huobi/pkg/postrequest"
	"fortune-bd/app/grid-strategy-svc/util/huobi/pkg/response/order"
)

// Responsible to handle orders request from WebSocket
// This need authentication version 1
type RequestOrdersWebSocketV1Client struct {
	websocketclientbase.WebSocketV1ClientBase
}

// Initializer
func (p *RequestOrdersWebSocketV1Client) Init(accessKey string, secretKey string, host string) *RequestOrdersWebSocketV1Client {
	p.WebSocketV1ClientBase.Init(accessKey, secretKey, host)
	return p
}

// Set callback biz
func (p *RequestOrdersWebSocketV1Client) SetHandler(
	authHandler websocketclientbase.AuthenticationV1ResponseHandler,
	responseHandler websocketclientbase.ResponseHandler) {
	p.WebSocketV1ClientBase.SetHandler(authHandler, p.handleMessage, responseHandler)
}

// Search past and open orders based on searching criteria.
func (p *RequestOrdersWebSocketV1Client) Request(req getrequest.RequestOrdersRequest) error {

	reqString, _ := postrequest.ToJson(req)
	return p.Send(reqString)
}

func (p *RequestOrdersWebSocketV1Client) handleMessage(msg string) (interface{}, error) {
	result := order.RequestOrdersV1Response{}
	err := json.Unmarshal([]byte(msg), &result)
	return result, err
}
