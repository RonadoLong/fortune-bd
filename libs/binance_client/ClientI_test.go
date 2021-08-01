package apiBinance

import (
	"github.com/shopspring/decimal"
	"testing"
)

func TestBinanceClient_GetAccountSpot(t *testing.T) {
	b := InitClient("p3nbQkUhKsD2vD6Nt6tsv5OQ8OK8IJVGrjDD6ZDx28Iganzha3gVIN6UOPTIWXR2", "HvgxFc2dtMYYKmgWjm90E7mEWrzbJfNyZ4yhPwbW0n0VBom4l9iJHlB96HPLcWq3")
	got, err := b.GetAccountSpot()
	if err != nil {
		t.Errorf("fail %v", err)
		return
	}
	for currency, account := range got.SubAccounts {
		if account.Amount != 0 {
			t.Logf("%s-%+v", currency, account)
		}
	}
}

func TestBinanceClient_CreateSubAccount(t *testing.T) {
	b := InitClient("lfNGLnHexoDNXEYeQGApIWb75ItHm7w7zOCJpxp1vvODIQFOFwChmuHxhvoleb1d", "8G3X3a3NxsZAyh3ZmEYRIX3d5DKK6PyqXyC6JylA0CQiQtafMZ8AUa8v8gRq43Sz")
	got, err := b.CreateSubAccount()
	if err != nil {
		t.Errorf("fail %v", err)
		return
	}
	t.Log(got)
}

func TestBinanceClient_EnableSubAccountMargin(t *testing.T) {
	b := InitClient("lfNGLnHexoDNXEYeQGApIWb75ItHm7w7zOCJpxp1vvODIQFOFwChmuHxhvoleb1d", "8G3X3a3NxsZAyh3ZmEYRIX3d5DKK6PyqXyC6JylA0CQiQtafMZ8AUa8v8gRq43Sz")
	err := b.EnableSubAccountMargin("502409971729264640")
	if err != nil {
		t.Errorf("fail %v", err)
		return
	}
}

func TestBinanceClient_ParentTransferToSubAccount(t *testing.T) {
	b := InitClient("lfNGLnHexoDNXEYeQGApIWb75ItHm7w7zOCJpxp1vvODIQFOFwChmuHxhvoleb1d", "8G3X3a3NxsZAyh3ZmEYRIX3d5DKK6PyqXyC6JylA0CQiQtafMZ8AUa8v8gRq43Sz")
	resp, err := b.ParentTransferToSubAccount("506724279742025728", "", "USDT", "10.65")
	if err != nil {
		t.Errorf("fail %v", err)
		return
	}
	t.Logf("%+v", resp)
}

func TestBinanceClient_SubAccountTransferToParent(t *testing.T) {
	b := InitClient("lfNGLnHexoDNXEYeQGApIWb75ItHm7w7zOCJpxp1vvODIQFOFwChmuHxhvoleb1d", "8G3X3a3NxsZAyh3ZmEYRIX3d5DKK6PyqXyC6JylA0CQiQtafMZ8AUa8v8gRq43Sz")
	resp, err := b.SubAccountTransferToParent("506724279742025728", "", "USDT", "5.65")
	if err != nil {
		t.Errorf("fail %v", err)
		return
	}
	t.Logf("%+v", resp)
}

func TestBinanceClient_GetTicks(t *testing.T) {
	//不要乱动这里
	//db := dbclient.NewDB("root:WQabc123@tcp(47.57.169.103:13306)/wq_fotune?charset=utf8mb4&parseTime=True&loc=Local")
	//b := InitClient("", "")
	//got, err := b.ApiClient.GetTickers()
	//if err != nil {
	//	t.Errorf("fail %v", err)
	//	return
	//}
	//for _, ticker := range got {
	//	if strings.HasSuffix(ticker.Symbol, "BTC") {
	//		sym := strings.ReplaceAll(ticker.Symbol, "BTC", "-BTC")
	//		t.Logf("%+v\n", sym)
	//		exec := db.Exec("INSERT INTO `wq_symbol` (`symbol`, `exchange`, `state`, `unit`) VALUE (?,?,?,?)", sym, "binance", 1, "btc")
	//		if exec.Error != nil {
	//			t.Errorf("错误 %v", exec.Error)
	//		}
	//	}
	//}
	//db.Close()
}

//
//func TestBinanceClient_GetSubAccountDepositAddress(t *testing.T) {
//	b := InitClient("lfNGLnHexoDNXEYeQGApIWb75ItHm7w7zOCJpxp1vvODIQFOFwChmuHxhvoleb1d", "8G3X3a3NxsZAyh3ZmEYRIX3d5DKK6PyqXyC6JylA0CQiQtafMZ8AUa8v8gRq43Sz")
//	resp, err := b.GetSubAccountDepositAddress("bikuanziying_47354259_47383985_brokersubuser@163.com", "USDT")
//	if err != nil {
//		t.Errorf("fail %v", err)
//		return
//	}
//	t.Logf("%+v", resp)
//}
//获取账户重置地址
func TestBinanceClient_GetAccountDepositAddress(t *testing.T) {
	b := InitClient("ybRzNzXcT4wy3wfm5BrZCaoesHuVkfL5eKkRGMPEDK6uV0OPuKLhz9CuaMBYmxJv", "GiyStrsnE8gpfFFqye0TDG3jQCm6bN18LLJjFEkeVrupBlKZHdBG9ZvMIoCvrC0F")
	resp, err := b.GetAccountDepositAddress("USDT")
	if err != nil {
		t.Errorf("fail %v", err)
		return
	}
	t.Logf("%+v", resp)
}

func TestBinanceClient_CreateSubAccountApi(t *testing.T) {
	b := InitClient("lfNGLnHexoDNXEYeQGApIWb75ItHm7w7zOCJpxp1vvODIQFOFwChmuHxhvoleb1d", "8G3X3a3NxsZAyh3ZmEYRIX3d5DKK6PyqXyC6JylA0CQiQtafMZ8AUa8v8gRq43Sz")
	resp, err := b.CreateSubAccountApi("502409971729264640", "true")
	if err != nil {
		t.Errorf("fail %v", err)
		return
	}
	t.Logf("%+v", resp)
}

//ybRzNzXcT4wy3wfm5BrZCaoesHuVkfL5eKkRGMPEDK6uV0OPuKLhz9CuaMBYmxJv
//GiyStrsnE8gpfFFqye0TDG3jQCm6bN18LLJjFEkeVrupBlKZHdBG9ZvMIoCvrC0F

func TestMy(t *testing.T) {
	gh := decimal.NewFromFloat(1)
	t.Log(gh.Cmp(decimal.NewFromFloat(-4)))
}
