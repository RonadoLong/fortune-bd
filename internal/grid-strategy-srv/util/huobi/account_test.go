// https://github.com/huobirdcenter/huobi_golang

package huobi

import (
	"fmt"
	"testing"

	"github.com/k0kubun/pp"
)

var (
	account   *Account
	accessKey = "dbuqg6hkte-23df26bc-aa64857f-7930e"
	secretKey = "bf7c125a-1ed0243f-716143f5-bdf6c"
)

func init() {
	var err error
	account, err = InitAccount(accessKey, secretKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("initAccount success")
}

func TestGetAccountID(t *testing.T) {
	ids, err := account.GetAccountID()
	if err != nil {
		t.Error(err)
		return
	}

	pp.Println(ids)
}

func TestGetAccountBalance(t *testing.T) {
	resp, err := account.GetAccountBalance()
	if err != nil {
		t.Error(err)
		return
	}

	pp.Println(resp)
}

func TestGetCurrencyBalance(t *testing.T) {
	resp, err := account.GetCurrencyBalance("usdt")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(resp)
}
