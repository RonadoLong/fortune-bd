package model

type WebSocketV2AuthenticationRequest struct {
	Action string  `json:"action"`
	Ch     string  `json:"ch"`
	Params *Params `json:"params"`
}

type Params struct {
	AuthType         string `json:"authType"`
	AccessKey        string `json:"accessKey"`
	SignatureMethod  string `json:"signatureMethod"`
	SignatureVersion string `json:"signatureVersion"`
	Timestamp        string `json:"timestamp"`
	Signature        string `json:"signature"`
}

func (p *WebSocketV2AuthenticationRequest) Init() *WebSocketV2AuthenticationRequest {

	p.Action = "req"
	p.Ch = "auth"
	p.Params = new(Params)
	p.Params.AuthType = "api"
	p.Params.SignatureMethod = "HmacSHA256"
	p.Params.SignatureVersion = "2.1"

	return p
}
