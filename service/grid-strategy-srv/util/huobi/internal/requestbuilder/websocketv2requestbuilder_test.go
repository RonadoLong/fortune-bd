package requestbuilder

import (
	"encoding/json"
	"testing"
	"time"

	"wq-grid-strategy/util/huobi/internal/model"
)

func TestWebSocketV2RequestBuilder_build_Time_Success(t *testing.T) {
	builder := new(WebSocketV2RequestBuilder).Init("access", "secret", "api.huobi.pro", "/ws/v2")
	utcDate := time.Date(2019, 11, 21, 10, 0, 0, 0, time.UTC)

	str, err := builder.build(utcDate)
	if err != nil {
		t.Error(err)
	}

	authReq := &model.WebSocketV2AuthenticationRequest{}
	err = json.Unmarshal([]byte(str), authReq)
	if err != nil {
		t.Error(err)
	}

	var actual, expected string

	expected = "req"
	actual = authReq.Action
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}

	expected = "auth"
	actual = authReq.Ch
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}

	expected = "api"
	actual = authReq.Params.AuthType
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}

	expected = "access"
	actual = authReq.Params.AccessKey
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}

	expected = "HmacSHA256"
	actual = authReq.Params.SignatureMethod
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}

	expected = "2.1"
	actual = authReq.Params.SignatureVersion
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}

	expected = "2019-11-21T10:00:00"
	actual = authReq.Params.Timestamp
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}

	expected = "1/d+cUIEh4tC0aXho86zu5QAxVzJaTe56mUiB275T0E="
	actual = authReq.Params.Signature
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}
