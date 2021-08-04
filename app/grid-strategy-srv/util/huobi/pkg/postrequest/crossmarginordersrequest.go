package postrequest

type CrossMarginOrdersRequest struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}
