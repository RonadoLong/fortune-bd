package common

type GetCurrenciesResponse struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}
