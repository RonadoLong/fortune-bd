package gateio

import (
	"net/http"
	"testing"
	"wq-fotune-backend/pkg/goex"
)

var gate = New(http.DefaultClient, "", "")

func TestGate_GetTicker(t *testing.T) {
	ticker, err := gate.GetTicker(goex.BTC_USDT)
	t.Log("err=>", err)
	t.Log("ticker=>", ticker)
}

func TestGate_GetDepth(t *testing.T) {
	dep, err := gate.GetDepth(1, goex.BTC_USDT)

	t.Log("err=>", err)
	t.Log("asks=>", dep.AskList)
	t.Log("bids=>", dep.BidList)
}
