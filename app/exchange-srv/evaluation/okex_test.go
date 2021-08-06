package evaluation

import (
	"log"
	"testing"
)

func TestOkex_CalProfit(t *testing.T) {
	o := &Okex{
		Symbol:      "BTC-SWAP",
		ShortVolume: 100,
		OpenPrice:   9510.22,
		Commission:  0.81457456,
		ClosePrice:  9710.22,
		Principal:   2000,
		Direction:   CloseBuy,
		Type:        ContractTradingUsdt,
	}
	got := o.CalProfit()
	log.Println(got)
	got2 := o.CalCommission()
	log.Println(got2)
	got3 := o.RateReturn()
	log.Println(got3)
}
