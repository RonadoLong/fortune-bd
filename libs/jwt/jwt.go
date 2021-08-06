package jwt

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
	"wq-fotune-backend/libs/encoding"
)

var (
	jwtHeader = JWTHeader{"HS256", "JWT"}
)

// JWT JWT
type JWT struct {
	Header  *JWTHeader
	Payload *JWTPayload
}

// JWTHeader token的header部分
type JWTHeader struct {
	Algorithm string `json:"alg"`
	TokenType string `json:"typ"`
}

// JWTPayload token的payload部分,由用户id和过期时间组成
type JWTPayload struct {
	UserID     string `json:"uid"`
	ExpireTime int64  `json:"exp"`
	Role       int    `json:"role"`
}

// NewJWT NewJWT
func NewJWT(uid string, role int, exp int64) *JWT {
	return &JWT{
		Header: &jwtHeader,
		Payload: &JWTPayload{
			UserID:     uid,
			ExpireTime: exp,
			Role:       role,
		},
	}
}

// SignedString 使用HS256加密标头返回token字符串
func (t *JWT) SignedString(secret string) (string, error) {
	header, err := encoding.Base64EncodeJSON(base64.RawURLEncoding, t.Header)
	if err != nil {
		return "", err
	}
	payload, err := encoding.Base64EncodeJSON(base64.RawURLEncoding, t.Payload)
	if err != nil {
		return "", err
	}
	signature := encoding.Base64Encode(base64.RawURLEncoding,
		encoding.HMAC(sha256.New, []byte(string(header)+"."+string(payload)), []byte(secret)))

	return strings.Join([]string{string(header), string(payload), string(signature)}, "."), nil
}

//JWTParse 解析token字符串并验证签名，如果一切正常，将设置有效负载
func (t *JWT) JWTParse(token, secret string) error {
	s := strings.Split(token, ".")
	if len(s) != 3 {
		return fmt.Errorf("token format error")
	}
	header, payload, signature := s[0], s[1], s[2]
	// 验证签名
	if !encoding.ConstTimeEqual([]byte(signature), encoding.Base64Encode(base64.RawURLEncoding,
		encoding.HMAC(sha256.New, []byte(header+"."+payload), []byte(secret)))) {
		return fmt.Errorf("token signature error")
	}
	// decode header
	if err := encoding.Base64DecodeJSON(base64.RawURLEncoding, []byte(header), &t.Header); err != nil {
		return err
	}
	// decode payload
	return encoding.Base64DecodeJSON(base64.RawURLEncoding, []byte(payload), &t.Payload)
}

//Expired checks if access token is expired
func (t *JWT) Expired() bool {
	return time.Now().After(time.Unix(t.Payload.ExpireTime, 0))
}
