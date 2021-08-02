package apiBinance

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/pkg/goex"
	"wq-fotune-backend/pkg/goex/binance"
)

type BinanceClient struct {
	ApiClient *binance.Binance
	Ws        *binance.BinanceWs
}

//todo 代理
var (
	client = &http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return &url.URL{
					Scheme: "socks5",
					Host:   strings.Split(env.ProxyAddr, "//")[1],
				}, nil
			},
		},
	}
)

func InitClient(apiKey, secret string) *BinanceClient {
	binanceClt := binance.NewWithConfig(&goex.APIConfig{
		HttpClient:    client,
		Endpoint:      "",
		ApiKey:        apiKey,
		ApiSecretKey:  secret,
		ApiPassphrase: "",
		ClientId:      "",
	})
	return &BinanceClient{
		ApiClient: binanceClt,
		Ws:        nil,
	}
}

func (b *BinanceClient) GetAccountSpot() (*goex.Account, error) {
	return b.ApiClient.GetAccount()
}

func (b *BinanceClient) GetAccountSpotUsdt() (float64, error) {
	account, err := b.ApiClient.GetAccount()
	if err != nil {
		return 0, err
	}
	for key, value := range account.SubAccounts {
		if key.Symbol == "USDT" {
			return value.Amount, nil
		}
	}
	return 0, nil
}

func (b *BinanceClient) GetAccountSwap() (*goex.Account, error) {
	return nil, errors.New("Not available")
}

func (b *BinanceClient) CheckIfApiValid() error {
	_, err := b.GetAccountSpot()
	return err
}

func (b *BinanceClient) CreateSubAccount() (accId string, err error) {
	return b.ApiClient.CreateSubAccount()
}

func (b *BinanceClient) EnableSubAccountMargin(subAccountId string) (err error) {
	return b.ApiClient.EnableSubAccountMargin(subAccountId)
}

func (b *BinanceClient) ParentTransferToSubAccount(toId, clientTranId, asset, amount string) (data binance.TransFerResp, err error) {
	return b.ApiClient.SubAccountTransfer("", toId, clientTranId, asset, amount)
}

func (b *BinanceClient) SubAccountTransferToParent(fromId, clientTranId, asset, amount string) (data binance.TransFerResp, err error) {
	return b.ApiClient.SubAccountTransfer(fromId, "", clientTranId, asset, amount)
}

//func (b *BinanceClient) GetSubAccountDepositAddress(email, symbol string) (resp binance.DepositAddrResp, err error) {
//	return b.ApiClient.GetSubAccountDepositAddress(email, symbol)
//}
//GetAccountDepositAddress
func (b *BinanceClient) GetAccountDepositAddress(symbol string) (resp binance.DepositAddrResp, err error) {
	return b.ApiClient.GetAccountDepositAddress(symbol)
}
func (b *BinanceClient) CreateSubAccountApi(subAccountId, canTrade string) (binance.SubAccountApiResp, error) {
	return b.ApiClient.CreateSubAccountApi(subAccountId, canTrade)
}
