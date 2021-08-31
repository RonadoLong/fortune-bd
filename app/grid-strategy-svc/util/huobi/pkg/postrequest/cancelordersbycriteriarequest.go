package postrequest

type CancelOrdersByCriteriaRequest struct {
	AccountId string `json:"account-id"`
	Symbol    string `json:"symbol,omitempty"`
	Side      string `json:"side,omitempty"`
	Size      int    `json:"size,omitempty"`
}
