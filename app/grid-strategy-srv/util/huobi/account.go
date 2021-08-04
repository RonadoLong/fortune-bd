package huobi

import (
	"errors"
	"fmt"
	"strconv"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi/pkg/client"
)

const (
	host = "api.huobi.pro"

	AccountSpot        = "spot"
	AccountMargin      = "margin"
	AccountSuperMargin = "super-margin"

	// FilledFees 永续合约已成交的手续费
	FilledFees = 0.002
)

// Accounter 通用交易所账号接口
//type Accounter interface {
//	GetCurrencyBalance(currency string) (float64, error)
//	PlaceLimitOrder(side string, symbol string, price string, amount string, clientOrderID string) (string, error)
//	PlaceMarketOrder(side string, symbol string, amount string, clientOrderID string) (string, error)
//	CancelOrder(orderID string) error
//	GetOrderInfo(orderID string) (interface{}, error)
//	GetHistoryOrdersInfo(symbol string, states string, types string) (interface{}, error)
//}

//Account 账号密钥
type Account struct {
	AccessKey   string `json:"accessKey"`
	SecretKey   string `json:"secretKey"`
	AccountID   string `json:"accountId"`
	AccountType string `json:"accountType"`
}

// InitAccount 实例化
func InitAccount(accessKey, secretKey string, accountTypes ...string) (*Account, error) {
	account := &Account{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}

	ids, err := account.GetAccountID()
	if err != nil {
		return account, err
	}

	// 默认为spot账号id
	var accountType string
	if len(accountTypes) == 0 {
		accountType = AccountSpot
	}

	account.AccountType = accountType
	account.AccountID = fmt.Sprintf("%d", ids[accountType])

	return account, nil
}

// GetAccountID 获取账号id，spot:现货用户，margin:逐仓杠杆交易用户，super-margin:全仓杠杆交易用户
func (a *Account) GetAccountID() (map[string]int64, error) {
	accounts := map[string]int64{}

	cli := new(client.AccountClient).Init(a.AccessKey, a.SecretKey, host)
	resp, err := cli.GetAccountInfo()
	if err != nil {
		return accounts, err
	}

	for _, v := range resp {
		if v.State == "working" {
			accounts[v.Type] = v.Id
		}
	}

	if len(accounts) == 0 {
		return accounts, errors.New("working account not fount")
	}

	return accounts, nil
}

// GetAccountBalance 获取账号余额，返回map，如果key不存在表示余额为零，其中key由"货币:类型"组成，类型： trade表示可使用余额，frozen表示冻结余额
func (a *Account) GetAccountBalance() (map[string]float64, error) {
	balances := map[string]float64{}

	if a.AccountID == "" {
		return balances, errors.New("accountID is empty, need init first")
	}

	cli := new(client.AccountClient).Init(a.AccessKey, a.SecretKey, host)
	resp, err := cli.GetAccountBalance(a.AccountID)
	if err != nil {
		return balances, err
	}

	for _, v := range resp.List {
		if v.Balance != "0" {
			balances[v.Currency+":"+v.Type] = str2Float64(v.Balance)
		}
	}

	if len(balances) == 0 {
		return balances, errors.New("all currency balances are empty")
	}

	return balances, nil
}

// GetCurrencyBalance 查询某个货币的余额
func (a *Account) GetCurrencyBalance(currency string) (float64, error) {
	balances, err := a.GetAccountBalance()
	if err != nil {
		return 0, err
	}

	val, ok := balances[currency+":trade"]
	if !ok {
		return 0, nil
	}

	return val, nil
}

func str2Float64(str string) float64 {
	if str == "" {
		return 0
	}

	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}

	return f
}
