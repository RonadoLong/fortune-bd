package requestbuilder

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
	"strings"
)

type Signer struct {
	hash hash.Hash
}

func (p *Signer) Init(key string) *Signer {
	p.hash = hmac.New(sha256.New, []byte(key))
	return p
}

func (p *Signer) Sign(method string, host string, path string, parameters string) string {
	if method == "" || host == "" || path == "" || parameters == "" {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(method)
	sb.WriteString("\n")
	sb.WriteString(host)
	sb.WriteString("\n")
	sb.WriteString(path)
	sb.WriteString("\n")
	sb.WriteString(parameters)

	return p.sign(sb.String())
}

func (p *Signer) sign(payload string) string {
	p.hash.Reset()
	p.hash.Write([]byte(payload))
	result := base64.StdEncoding.EncodeToString(p.hash.Sum(nil))
	return result
}
