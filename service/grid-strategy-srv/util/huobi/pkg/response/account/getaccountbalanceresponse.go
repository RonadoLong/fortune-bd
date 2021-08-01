package account

type GetAccountBalanceResponse struct {
	Status string          `json:"status"`
	Data   *AccountBalance `json:"data"`
}
