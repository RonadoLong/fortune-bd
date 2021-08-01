package requestbuilder

import (
	"fmt"
	"net/url"
	"time"
	"wq-fotune-backend/service/grid-strategy-srv/util/huobi/pkg/getrequest"
)

type PrivateUrlBuilder struct {
	host    string
	akKey   string
	akValue string
	smKey   string
	smValue string
	svKey   string
	svValue string
	tKey    string

	signer *Signer
}

func (p *PrivateUrlBuilder) Init(accessKey string, secretKey string, host string) *PrivateUrlBuilder {
	p.akKey = "AccessKeyId"
	p.akValue = accessKey
	p.smKey = "SignatureMethod"
	p.smValue = "HmacSHA256"
	p.svKey = "SignatureVersion"
	p.svValue = "2"
	p.tKey = "Timestamp"

	p.host = host
	p.signer = new(Signer).Init(secretKey)

	return p
}

func (p *PrivateUrlBuilder) Build(method string, path string, request *getrequest.GetRequest) string {
	time := time.Now().UTC()

	return p.BuildWithTime(method, path, time, request)
}

func (p *PrivateUrlBuilder) BuildWithTime(method string, path string, utcDate time.Time, request *getrequest.GetRequest) string {
	time := utcDate.Format("2006-01-02T15:04:05")

	req := new(getrequest.GetRequest).InitFrom(request)
	req.AddParam(p.akKey, p.akValue)
	req.AddParam(p.smKey, p.smValue)
	req.AddParam(p.svKey, p.svValue)
	req.AddParam(p.tKey, time)

	parameters := req.BuildParams()

	signature := p.signer.Sign(method, p.host, path, parameters)

	url := fmt.Sprintf("https://%s%s?%s&Signature=%s", p.host, path, parameters, url.QueryEscape(signature))

	return url
}
