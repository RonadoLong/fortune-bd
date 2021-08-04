package common

type GetV2ReferenceCurrenciesResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    []CurrencyChain `json:"data"`
}

type CurrencyChain struct {
	Currency   string   `json:"currency"`
	InstStatus string   `json:"instStatus"`
	Chains     []Chains `json:"chains"`
}

type Chains struct {
	Chain                   string `json:"chain"`
	BaseChain               string `json:"baseChain"`
	BaseChainProtocol       string `json:"baseChainProtocol"`
	IsDynamic               bool   `json:"isDynamic"`
	NumOfConfirmations      int    `json:"numOfConfirmations"`
	NumOfFastConfirmations  int    `json:"numOfFastConfirmations"`
	MinDepositAmt           string `json:"minDepositAmt"`
	DepositStatus           string `json:"depositStatus"`
	MinWithdrawAmt          string `json:"minWithdrawAmt"`
	MaxWithdrawAmt          string `json:"maxWithdrawAmt"`
	WithdrawQuotaPerDay     string `json:"withdrawQuotaPerDay"`
	WithdrawQuotaPerYear    string `json:"withdrawQuotaPerYear"`
	WithdrawQuotaTotal      string `json:"withdrawQuotaTotal"`
	WithdrawPrecision       int    `json:"withdrawPrecision"`
	WithdrawFeeType         string `json:"withdrawFeeType"`
	TransactFeeWithdraw     string `json:"transactFeeWithdraw"`
	MinTransactFeeWithdraw  string `json:"minTransactFeeWithdraw"`
	MaxTransactFeeWithdraw  string `json:"maxTransactFeeWithdraw"`
	TransactFeeRateWithdraw string `json:"transactFeeRateWithdraw"`
	WithdrawStatus          string `json:"withdrawStatus"`
}
