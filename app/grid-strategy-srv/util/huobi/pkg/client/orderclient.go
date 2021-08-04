package client

import (
	"encoding/json"
	"fmt"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/internal"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/internal/requestbuilder"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/pkg/getrequest"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/pkg/postrequest"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/pkg/response/order"
)

// Responsible to operate on order
type OrderClient struct {
	privateUrlBuilder *requestbuilder.PrivateUrlBuilder
}

// Initializer
func (p *OrderClient) Init(accessKey string, secretKey string, host string) *OrderClient {
	p.privateUrlBuilder = new(requestbuilder.PrivateUrlBuilder).Init(accessKey, secretKey, host)
	return p
}

// Place a new order and send to the exchange to be matched.
func (p *OrderClient) PlaceOrder(request *postrequest.PlaceOrderRequest) (*order.PlaceOrderResponse, error) {
	postBody, jsonErr := postrequest.ToJson(request)

	url := p.privateUrlBuilder.Build("POST", "/v1/order/orders/place", nil)
	postResp, postErr := internal.HttpPost(url, postBody)
	if postErr != nil {
		return nil, postErr
	}

	result := order.PlaceOrderResponse{}
	jsonErr = json.Unmarshal([]byte(postResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &result, nil
}

// Place multipler orders (at most 10 orders)
func (p *OrderClient) PlaceOrders(request []postrequest.PlaceOrderRequest) (*order.PlaceOrdersResponse, error) {

	postBody, jsonErr := postrequest.ToJson(request)

	url := p.privateUrlBuilder.Build("POST", "/v1/order/batch-orders", nil)
	postResp, postErr := internal.HttpPost(url, string(postBody))
	if postErr != nil {
		return nil, postErr
	}

	result := order.PlaceOrdersResponse{}
	jsonErr = json.Unmarshal([]byte(postResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &result, nil
}

// Cancel an order by order id
func (p *OrderClient) CancelOrderById(orderId string) (*order.CancelOrderByIdResponse, error) {
	path := fmt.Sprintf("/v1/order/orders/%s/submitcancel", orderId)
	url := p.privateUrlBuilder.Build("POST", path, nil)
	postResp, postErr := internal.HttpPost(url, "")
	if postErr != nil {
		return nil, postErr
	}

	result := order.CancelOrderByIdResponse{}
	jsonErr := json.Unmarshal([]byte(postResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &result, nil
}

// Cancel an order by client order id
func (p *OrderClient) CancelOrderByClientOrderId(clientOrderId string) (*order.CancelOrderByClientResponse, error) {
	url := p.privateUrlBuilder.Build("POST", "/v1/order/orders/submitCancelClientOrder", nil)
	body := fmt.Sprintf("{\"client-order-id\":\"%s\"}", clientOrderId)
	postResp, postErr := internal.HttpPost(url, body)
	if postErr != nil {
		return nil, postErr
	}

	result := order.CancelOrderByClientResponse{}
	jsonErr := json.Unmarshal([]byte(postResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &result, nil
}

// Returns all open orders which have not been filled completely.
func (p *OrderClient) GetOpenOrders(request *getrequest.GetRequest) (*order.GetOpenOrdersResponse, error) {
	url := p.privateUrlBuilder.Build("GET", "/v1/order/openOrders", request)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := order.GetOpenOrdersResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &result, nil
}

// Submit cancellation for multiple orders at once with given criteria.
func (p *OrderClient) CancelOrdersByCriteria(request *postrequest.CancelOrdersByCriteriaRequest) (*order.CancelOrdersByCriteriaResponse, error) {
	postBody, jsonErr := postrequest.ToJson(request)

	url := p.privateUrlBuilder.Build("POST", "/v1/order/orders/batchCancelOpenOrders", nil)
	postResp, postErr := internal.HttpPost(url, string(postBody))
	if postErr != nil {
		return nil, postErr
	}

	result := order.CancelOrdersByCriteriaResponse{}
	jsonErr = json.Unmarshal([]byte(postResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &result, nil
}

// Submit cancellation for multiple orders at once with given ids
func (p *OrderClient) CancelOrdersByIds(request *postrequest.CancelOrdersByIdsRequest) (*order.CancelOrdersByIdsResponse, error) {
	postBody, jsonErr := postrequest.ToJson(request)

	url := p.privateUrlBuilder.Build("POST", "/v1/order/orders/batchcancel", nil)
	postResp, postErr := internal.HttpPost(url, string(postBody))
	if postErr != nil {
		return nil, postErr
	}

	result := order.CancelOrdersByIdsResponse{}
	jsonErr = json.Unmarshal([]byte(postResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &result, nil
}

// Returns the detail of one order by order id
func (p *OrderClient) GetOrderById(orderId string) (*order.GetOrderResponse, error) {
	path := fmt.Sprintf("/v1/order/orders/%s", orderId)
	url := p.privateUrlBuilder.Build("GET", path, nil)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := order.GetOrderResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &result, nil
}

// Returns the detail of one order by client order id
func (p *OrderClient) GetOrderByCriteria(request *getrequest.GetRequest) (*order.GetOrderResponse, error) {
	url := p.privateUrlBuilder.Build("GET", "/v1/order/orders/getClientOrder", request)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := order.GetOrderResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &result, nil
}

// Returns the match result of an order.
func (p *OrderClient) GetMatchResultsById(orderId string) (*order.GetMatchResultsResponse, error) {
	path := fmt.Sprintf("/v1/order/orders/%s/matchresults", orderId)
	url := p.privateUrlBuilder.Build("GET", path, nil)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := order.GetMatchResultsResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &result, nil
}

// Returns orders based on a specific searching criteria.
func (p *OrderClient) GetHistoryOrders(request *getrequest.GetRequest) (*order.GetHistoryOrdersResponse, error) {
	url := p.privateUrlBuilder.Build("GET", "/v1/order/orders", request)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := order.GetHistoryOrdersResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &result, nil
}

// Returns orders based on a specific searching criteria.
func (p *OrderClient) GetLast48hOrders(request *getrequest.GetRequest) (*order.GetHistoryOrdersResponse, error) {
	url := p.privateUrlBuilder.Build("GET", "/v1/order/history", request)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := order.GetHistoryOrdersResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &result, nil
}

// Returns the match results of past and open orders based on specific search criteria.
func (p *OrderClient) GetMatchResultsByCriteria(request *getrequest.GetRequest) (*order.GetMatchResultsResponse, error) {
	url := p.privateUrlBuilder.Build("GET", "/v1/order/matchresults", request)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := order.GetMatchResultsResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &result, nil
}

// Returns the current transaction fee rate applied to the user.
func (p *OrderClient) GetTransactFeeRate(request *getrequest.GetRequest) (*order.GetTransactFeeRateResponse, error) {
	url := p.privateUrlBuilder.Build("GET", "/v2/reference/transact-fee-rate", request)
	getResp, getErr := internal.HttpGet(url)
	if getErr != nil {
		return nil, getErr
	}

	result := order.GetTransactFeeRateResponse{}
	jsonErr := json.Unmarshal([]byte(getResp), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &result, nil
}
