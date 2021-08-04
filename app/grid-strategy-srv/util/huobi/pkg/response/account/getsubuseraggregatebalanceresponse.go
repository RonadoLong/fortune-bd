package account

type GetSubUserAggregateBalanceResponse struct {
	Status string    `json:"status"`
	Data   []Balance `json:"data"`
}
