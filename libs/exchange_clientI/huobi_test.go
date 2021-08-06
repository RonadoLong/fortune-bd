package exchange_clientI

import (
	"log"
	"testing"
)

func TestInitClient(t *testing.T) {

	Clt := InitClient("26324e15-bgrveg5tmn-6e3a6eb8-09a97", "cf629bd3-9b712f10-6f5b240b-334c0", true)
	log.Printf("%+v", Clt)
}

func TestGetAccount(t *testing.T) {
	Clt := InitClient("26324e15-bgrveg5tmn-6e3a6eb8-09a97", "cf629bd3-9b712f10-6f5b240b-334c", true)
	account, err := Clt.GetAccountSpot()
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range account.SubAccounts {
		log.Printf("%+v", v)
	}
}

func TestGetAccountSwap(t *testing.T) {
	Clt := InitClient("26324e15-bgrveg5tmn-6e3a6eb8-09a97", "cf629bd3-9b712f10-6f5b240b-334c0", true)
	account, err := Clt.GetAccountSwap()
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range account.SubAccounts {
		log.Printf("%+v", v)
	}
}
