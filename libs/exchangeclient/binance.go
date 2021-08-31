package exchangeclient

import (
	"errors"
	"fortune-bd/libs/goex"
	"fortune-bd/libs/goex/binance"
)

type Binance struct {
	ApiClient *binance.Binance
	Ws        *binance.BinanceWs
}

func InitBinance(apiKey, secret string) *Binance {
	binanceClt := binance.NewWithConfig(&goex.APIConfig{
		HttpClient:    client,
		Endpoint:      "",
		ApiKey:        apiKey,
		ApiSecretKey:  secret,
		ApiPassphrase: "",
		ClientId:      "",
	})
	return &Binance{
		ApiClient: binanceClt,
		Ws:        nil,
	}
}

func (b *Binance) GetAccountSpot() (*goex.Account, error) {
	return b.ApiClient.GetAccount()
}

func (b *Binance) GetAccountSpotUsdt() (float64, error) {
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

func (b *Binance) GetAccountSwap() (*goex.Account, error) {
	return nil, errors.New("Not available")
}

func (b *Binance) CheckIfApiValid() error {
	_, err := b.GetAccountSpot()
	return err
}

func (b *Binance) CreateSubAccount() (accId string, err error) {
	return b.ApiClient.CreateSubAccount()
}

func (b *Binance) EnableSubAccountMargin(subAccountId string) (err error) {
	return b.ApiClient.EnableSubAccountMargin(subAccountId)
}

func (b *Binance) ParentTransferToSubAccount(toId, clientTranId, asset, amount string) (data binance.TransFerResp, err error) {
	return b.ApiClient.SubAccountTransfer("", toId, clientTranId, asset, amount)
}

func (b *Binance) SubAccountTransferToParent(fromId, clientTranId, asset, amount string) (data binance.TransFerResp, err error) {
	return b.ApiClient.SubAccountTransfer(fromId, "", clientTranId, asset, amount)
}

func (b *Binance) GetAccountDepositAddress(symbol string) (resp binance.DepositAddrResp, err error) {
	return b.ApiClient.GetAccountDepositAddress(symbol)
}
func (b *Binance) CreateSubAccountApi(subAccountId, canTrade string) (binance.SubAccountApiResp, error) {
	return b.ApiClient.CreateSubAccountApi(subAccountId, canTrade)
}
