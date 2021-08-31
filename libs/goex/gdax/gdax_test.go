package gdax

import (
	"net/http"
	"testing"

	"wq-fotune-backend/pkg/goex"
	"wq-fotune-backend/pkg/goex/internal/logger"
)

var gdax = New(http.DefaultClient, "", "")

func TestGdax_GetTicker(t *testing.T) {
	ticker, err := gdax.GetTicker(goex.BTC_USD)
	t.Log("err=>", err)
	t.Log("ticker=>", ticker)
}

func TestGdax_Get24HStats(t *testing.T) {
	stats, err := gdax.Get24HStats(goex.BTC_USD)
	t.Log("err=>", err)
	t.Log("stats=>", stats)
}

func TestGdax_GetDepth(t *testing.T) {
	dep, err := gdax.GetDepth(2, goex.BTC_USD)
	t.Log("err=>", err)
	t.Log("bids=>", dep.BidList)
	t.Log("asks=>", dep.AskList)
}

func TestGdax_GetKlineRecords(t *testing.T) {
	logger.SetLevel(logger.DEBUG)
	t.Log(gdax.GetKlineRecords(goex.BTC_USD, goex.KLINE_PERIOD_1DAY, 0, 0))
}
