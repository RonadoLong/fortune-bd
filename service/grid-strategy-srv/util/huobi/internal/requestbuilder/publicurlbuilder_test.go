package requestbuilder

import (
	"testing"
	"wq-fotune-backend/service/grid-strategy-srv/util/huobi/pkg/getrequest"
)

func TestPublicUrlBuilder_Build_NoRequestParameter_Success(t *testing.T) {
	builder := new(PublicUrlBuilder).Init("api.huobi.pro")

	result := builder.Build("/v1/common/symbols", nil)

	expected := "https://api.huobi.pro/v1/common/symbols"
	if result != expected {
		t.Errorf("expected: %s, actual: %s", expected, result)
	}
}

func TestPublicUrlBuilder_Build_HasRequestParameter_Success(t *testing.T) {
	builder := new(PublicUrlBuilder).Init("api.huobi.pro")
	reqParams := new(getrequest.GetRequest).Init()
	reqParams.AddParam("symbol", "btcusdt")
	reqParams.AddParam("period", "1min")
	reqParams.AddParam("size", "1")

	result := builder.Build("/market/history/kline", reqParams)

	expected := "https://api.huobi.pro/market/history/kline?period=1min&size=1&symbol=btcusdt"
	if result != expected {
		t.Errorf("expected: %s, actual: %s", expected, result)
	}
}
