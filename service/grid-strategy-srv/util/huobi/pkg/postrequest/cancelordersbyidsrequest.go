package postrequest

type CancelOrdersByIdsRequest struct {
	OrderIds       []string `json:"order-ids,omitempty"`
	ClientOrderIds []string `json:"client-order-ids,omitempty"`
}
