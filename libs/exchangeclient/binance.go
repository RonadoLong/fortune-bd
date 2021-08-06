package exchangeclient

import (
	"errors"
	"wq-fotune-backend/pkg/goex"
	"wq-fotune-backend/pkg/goex/binance"
)

type BinanceClient struct {
	ApiClient *binance.Binance
	Ws        *binance.BinanceWs
}

func InitBinance(apiKey, secret string) *BinanceClient {
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

func (b *BinanceClient) GetAccountDepositAddress(symbol string) (resp binance.DepositAddrResp, err error) {
	return b.ApiClient.GetAccountDepositAddress(symbol)
}
func (b *BinanceClient) CreateSubAccountApi(subAccountId, canTrade string) (binance.SubAccountApiResp, error) {
	return b.ApiClient.CreateSubAccountApi(subAccountId, canTrade)
}
