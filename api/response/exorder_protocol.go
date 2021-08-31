package response

import "errors"

type ExchangeApiReq struct {
	ExchangeID int64  `json:"exchange_id"`
	ApiKey     string `json:"api_key"`
	Secret     string `json:"secret"`
	Passphrase string `json:"passphrase"`
}

func (req *ExchangeApiReq) CheckNotNull() error {
	switch "" {
	case req.ApiKey:
		return errors.New("Apikey 不能为空")
	case req.Secret:
		return errors.New("Secret 不能为空")
	}
	return nil
}

type UpdateExchangeApiReq struct {
	ExchangeApiReq
	ApiId int64 `json:"api_id"`
}

type SetUserStrategyApiReq struct {
	StrategyID string `json:"strategy_id"`
	ApiKey     string `json:"api_key"`
}

func (req *SetUserStrategyApiReq) CheckNotNull() error {
	switch "" {
	case req.ApiKey:
		return errors.New("Apikey不能为空")
	case req.StrategyID:
		return errors.New("策略id不能为空")
	}
	return nil
}

type SetUserStrategyBalanceReq struct {
	StrategyID string  `json:"strategy_id"`
	Balance    float32 `json:"balance"`
}

func (req *SetUserStrategyBalanceReq) CheckParam() error {
	if req.StrategyID == "" {
		return errors.New("策略id不能为空")
	}
	if req.Balance == 0.0 {
		return errors.New("资金不能为0")
	}
	return nil
}

type BuyStrategyReq struct {
	Id      int64   `json:"id"`
	Balance float32 `json:"balance"`
}
