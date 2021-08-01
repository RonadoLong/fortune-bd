package requestbuilder

import (
	"encoding/json"
	"testing"
	"time"
	"wq-fotune-backend/service/grid-strategy-srv/util/huobi/internal/model"
)

func TestWebSocketV1RequestBuilder_build_Time_Success(t *testing.T) {
	builder := new(WebSocketV1RequestBuilder).Init("access", "secret", "api.huobi.pro", "/ws/v1")
	utcDate := time.Date(2019, 11, 21, 10, 0, 0, 0, time.UTC)

	str, err := builder.build(utcDate)
	if err != nil {
		t.Error(err)
	}

	authReq := &model.WebSocketV1AuthenticationRequest{}
	err = json.Unmarshal([]byte(str), authReq)
	if err != nil {
		t.Error(err)
	}

	var actual, expected string

	expected = "auth"
	actual = authReq.Op
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}

	expected = "access"
	actual = authReq.AccessKeyId
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}

	expected = "HmacSHA256"
	actual = authReq.SignatureMethod
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}

	expected = "2019-11-21T10:00:00"
	actual = authReq.Timestamp
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}

	expected = "2"
	actual = authReq.SignatureVersion
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}

	expected = "nWj8xkaQ8mWPyvdtRVPFkrX2B8v3mSomAfhXiOGoS3M="
	actual = authReq.Signature
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}
