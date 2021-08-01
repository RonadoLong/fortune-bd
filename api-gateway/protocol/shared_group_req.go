package protocol

type SharedGroup struct {
	ID             int64        `json:"id"`
	GroupName      string       `json:"groupName"`
	TotalCapitals  string       `json:"totalCapitals"`
	CapitalUnit    string       `json:"capitalUnit"`
	DistributeType int          `json:"distributeType"`
	LeverageRatio  string       `json:"leverageRatio"`
	TotalReturn    string       `json:"total_return"`
	AnnualReturn   string       `json:"annual_return"`
	MaxDdpercent   string       `json:"max_ddpercent"`
	CalmarRatio    string       `json:"calmar_ratio"`
	SharpeRatio    string       `json:"sharpe_ratio"`
	Strategies     []Strategies `json:"strategies"`
}

type Strategies struct {
	ID              int64  `json:"id"`
	StrategyName    string `json:"strategyName"`
	DistributeRatio string `json:"distributeRatio"`
	Exchange        string `json:"exchange"`
	Symbol          string `json:"symbol"`
}
