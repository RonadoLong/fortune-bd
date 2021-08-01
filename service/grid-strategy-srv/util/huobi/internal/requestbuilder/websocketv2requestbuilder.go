package requestbuilder

import (
	"time"
	"wq-fotune-backend/service/grid-strategy-srv/util/huobi/internal/model"
	"wq-fotune-backend/service/grid-strategy-srv/util/huobi/pkg/getrequest"
	"wq-fotune-backend/service/grid-strategy-srv/util/huobi/pkg/postrequest"
)

type WebSocketV2RequestBuilder struct {
	akKey   string
	akValue string
	smKey   string
	smValue string
	svKey   string
	svValue string
	tKey    string
	tValue  string

	host string
	path string

	signer *Signer
}

func (p *WebSocketV2RequestBuilder) Init(accessKey string, secretKey string, host string, path string) *WebSocketV2RequestBuilder {
	p.akKey = "accessKey"
	p.akValue = accessKey
	p.smKey = "signatureMethod"
	p.smValue = "HmacSHA256"
	p.svKey = "signatureVersion"
	p.svValue = "2.1"
	p.tKey = "timestamp"

	p.host = host
	p.path = path

	p.signer = new(Signer).Init(secretKey)

	return p
}

func (p *WebSocketV2RequestBuilder) Build() (string, error) {
	time := time.Now().UTC()
	return p.build(time)
}

func (p *WebSocketV2RequestBuilder) build(utcDate time.Time) (string, error) {
	time := utcDate.Format("2006-01-02T15:04:05")

	req := new(getrequest.GetRequest).Init()
	req.AddParam(p.akKey, p.akValue)
	req.AddParam(p.smKey, p.smValue)
	req.AddParam(p.svKey, p.svValue)
	req.AddParam(p.tKey, time)

	signature := p.signer.Sign("GET", p.host, p.path, req.BuildParams())

	auth := new(model.WebSocketV2AuthenticationRequest).Init()
	auth.Params.AccessKey = p.akValue
	auth.Params.Timestamp = time
	auth.Params.Signature = signature

	result, err := postrequest.ToJson(auth)
	return result, err
}
