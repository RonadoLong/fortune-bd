package exchange

import (
	"fmt"
	"wq-fotune-backend/pkg/goex"
	"wq-fotune-backend/pkg/goex/okex"
)

type OKEXClient struct {
	APIClient *okex.OKEx
	WS        *okex.OKExV3FuturesWs
}

func InitOKEX(info Info) *OKEXClient {
	OKExCli := okex.NewOKEx(&goex.APIConfig{
		HttpClient:    client,
		ApiKey:        info.ApiKey,
		ApiSecretKey:  info.ApiSecretKey,
		ApiPassphrase: info.ApiPassphrase,
		Endpoint:      "https://www.okex.com",
	})

	return &OKEXClient{
		APIClient: OKExCli,
	}
}

func (O *OKEXClient) BuildWS() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("BuildWS panic %+v", e)
		}
	}()
	base := okex.NewOKEx(O.APIClient.Config)
	ws := okex.NewOKExV3FuturesWs(base)
	err = ws.Login()
	if err != nil {
		return
	}
	O.WS = ws
	return
}

func (O OKEXClient) GetExchangeName() string {
	return OKEX
}

func (O OKEXClient) Buy() {
	panic("implement me")
}

func (O OKEXClient) Sell() {
	panic("implement me")
}

func (O OKEXClient) GetAllOrder() {
	panic("implement me")
}

func (O OKEXClient) GetUnfinishOrder() {
	panic("implement me")
}

func (O OKEXClient) SubOrderCallBack() {
	panic("implement me")
}

func (O OKEXClient) SubTraderCallBack() {
	panic("implement me")
}
