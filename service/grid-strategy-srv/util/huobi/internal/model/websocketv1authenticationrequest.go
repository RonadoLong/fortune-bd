package model

type WebSocketV1AuthenticationRequest struct {
	Op               string `json:"op"`
	AccessKeyId      string
	SignatureMethod  string
	SignatureVersion string
	Timestamp        string
	Signature        string
}

func (p *WebSocketV1AuthenticationRequest) Init() *WebSocketV1AuthenticationRequest {
	p.Op = "auth"
	p.SignatureMethod = "HmacSHA256"
	p.SignatureVersion = "2"

	return p
}
