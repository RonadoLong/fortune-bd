package requestbuilder

import (
	"testing"
	"time"
	"wq-fotune-backend/internal/grid-strategy-srv/util/huobi/pkg/getrequest"
)

func TestPrivateUrlBuilder_Build_NoRequestParameter_Success(t *testing.T) {
	builder := new(PrivateUrlBuilder).Init("access", "secret", "api.huobi.pro")
	utcDate := time.Date(2019, 11, 21, 10, 0, 0, 0, time.UTC)

	result := builder.BuildWithTime("GET", "/v1/account/accounts", utcDate, nil)

	expected := "https://api.huobi.pro/v1/account/accounts?AccessKeyId=access&SignatureMethod=HmacSHA256&SignatureVersion=2&Timestamp=2019-11-21T10%3A00%3A00&Signature=rWnLcMt3XBAsmXoNHtTQVpvMbH%2FcE1PXFwQAGeYwt3s%3D"
	if result != expected {
		t.Errorf("expected: %s, actual: %s", expected, result)
	}
}

func TestPrivateUrlBuilder_Build_HasRequestParameter_Success(t *testing.T) {
	builder := new(PrivateUrlBuilder).Init("access", "secret", "api.huobi.pro")
	utcDate := time.Date(2019, 11, 21, 10, 0, 0, 0, time.UTC)
	reqParams := new(getrequest.GetRequest).Init()
	reqParams.AddParam("account-id", "123")
	reqParams.AddParam("currency", "btc")

	result := builder.BuildWithTime("GET", "/v1/account/history", utcDate, reqParams)

	expected := "https://api.huobi.pro/v1/account/history?AccessKeyId=access&SignatureMethod=HmacSHA256&SignatureVersion=2&Timestamp=2019-11-21T10%3A00%3A00&account-id=123&currency=btc&Signature=SGZYJ9Ub%2FhFerEBbSWsCxl8Djk%2BLRBgEZOB4fLc4T9Q%3D"
	if result != expected {
		t.Errorf("expected: %s, actual: %s", expected, result)
	}
}
