package order

type GetMatchResultsResponse struct {
	Status       string `json:"status"`
	ErrorCode    string `json:"err-code"`
	ErrorMessage string `json:"err-msg"`
	Data         []struct {
		Id                int64  `json:"id"`
		OrderId           int64  `json:"order-id"`
		MatchId           int64  `json:"match-id"`
		TradeId           int64  `json:"trade-id"`
		Symbol            string `json:"symbol"`
		Price             string `json:"price"`
		CreatedAt         int64  `json:"created-at"`
		Type              string `json:"type"`
		FilledAmount      string `json:"filled-amount"`
		FilledFees        string `json:"filled-fees"`
		Source            string `json:"source"`
		Role              string `json:"role"`
		FilledPoints      string `json:"filled-points"`
		FeeDeductCurrency string `json:"fee-deduct-currency"`
	}
}
