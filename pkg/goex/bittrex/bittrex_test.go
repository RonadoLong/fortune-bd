package bittrex

import (
	"net/http"
	"testing"
	"wq-fotune-backend/pkg/goex"
)

var b = New(http.DefaultClient, "", "")

func TestBittrex_GetTicker(t *testing.T) {
	ticker, err := b.GetTicker(goex.BTC_USDT)
	t.Log("err=>", err)
	t.Log("ticker=>", ticker)
}

func TestBittrex_GetDepth(t *testing.T) {
	dep, err := b.GetDepth(1, goex.BTC_USDT)
	t.Log("err=>", err)
	t.Log("ask=>", dep.AskList)
	t.Log("bid=>", dep.BidList)
}
