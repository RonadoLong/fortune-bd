package response

type SharedGroup struct {
	ID             int64        `json:"id"`
	GroupName      string       `json:"group_name"`
	TotalCapitals  string       `json:"total_capitals"`
	CapitalUnit    string       `json:"capital_unit"`
	DistributeType int          `json:"distribute_type"`
	LeverageRatio  string       `json:"leverage_ratio"`
	TotalReturn    string       `json:"total_return"`
	AnnualReturn   string       `json:"annual_return"`
	MaxDdpercent   string       `json:"max_ddpercent"`
	CalmarRatio    string       `json:"calmar_ratio"`
	SharpeRatio    string       `json:"sharpe_ratio"`
	Strategies     []Strategies `json:"strategies"`
}

type Strategies struct {
	ID              int64  `json:"id"`
	StrategyName    string `json:"strategy_name"`
	DistributeRatio string `json:"distribute_ratio"`
	Exchange        string `json:"exchange"`
	Symbol          string `json:"symbol"`
}
