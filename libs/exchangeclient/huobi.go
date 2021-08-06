package exchangeclient

import (
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/pkg/goex"
	"wq-fotune-backend/pkg/goex/huobi"
)

var (
	TypeSpot        = "spot"
	TypeMargin      = "margin"
	TypeSuperMargin = "super-margin"
)

type HuobiClient struct {
	APIClient        *huobi.HuoBiPro
	WS               *huobi.SpotWs
	accIdSpot        string //现货账户id
	accIdMargin      string //逐仓杠杆账户id
	accIdSuperMargin string //全仓杠杆账户id  需要在initClient赋值
}

func InitHuobi(apiKey, apiSecretKey string, initID bool) *HuobiClient {
	huobiPro := huobi.NewHuoBiPro(client, apiKey, apiSecretKey, "")
	huobiClt := &HuobiClient{
		APIClient: huobiPro,
	}
	if !initID {
		return huobiClt
	}
	spotAcc, err := huobiClt.GetAccountId(TypeSpot)
	if err != nil {
		logger.Warnf("apiKey%s 获取现货accountId 失败 %v", apiKey, err)
	} else {
		huobiClt.accIdSpot = spotAcc.Id
	}

	marginAcc, err := huobiClt.GetAccountId(TypeMargin)
	if err != nil {
		logger.Warnf("apiKey%s 获取逐仓账户accountId 失败 %v", apiKey, err)
	} else {
		huobiClt.accIdMargin = marginAcc.Id
	}
	superMarginAcc, err := huobiClt.GetAccountId(TypeSuperMargin)
	if err != nil {
		logger.Warnf("apiKey%s 获取全仓账户accountId 失败 %v", apiKey, err)
	} else {
		huobiClt.accIdSuperMargin = superMarginAcc.Id
	}
	return huobiClt
}

func (h *HuobiClient) GetAccountSpot() (*goex.Account, error) {
	h.APIClient.AccountId = h.accIdSpot
	return h.APIClient.GetAccount()
}

func (h *HuobiClient) GetAccountSwap() (*goex.Account, error) {
	h.APIClient.AccountId = h.accIdSuperMargin
	return h.APIClient.GetAccount()
}

func (h *HuobiClient) GetAccountId(_type string) (acc huobi.AccountInfo, err error) {
	return h.APIClient.GetAccountInfo(_type)
}

func (h *HuobiClient) CheckIfApiValid() error {
	_, err := h.GetAccountId(TypeSpot)
	return err
}
