package wallet

type GetWithdrawQuotaResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    *WithdrawQuota `json:"data"`
}
type WithdrawQuota struct {
	Currency string `json:"currency"`
	Chains   []struct {
		Chain                      string `json:"chain"`
		MaxWithdrawAmt             string `json:"maxWithdrawAmt"`
		WithdrawQuotaPerDay        string `json:"withdrawQuotaPerDay"`
		RemainWithdrawQuotaPerDay  string `json:"remainWithdrawQuotaPerDay"`
		WithdrawQuotaPerYear       string `json:"withdrawQuotaPerYear"`
		RemainWithdrawQuotaPerYear string `json:"remainWithdrawQuotaPerYear"`
		WithdrawQuotaTotal         string `json:"withdrawQuotaTotal"`
		RemainWithdrawQuotaTotal   string `json:"remainWithdrawQuotaTotal"`
	} `json:"chains"`
}
