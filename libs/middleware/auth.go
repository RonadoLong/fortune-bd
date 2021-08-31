package middleware

import (
	"fortune-bd/libs/jwt"
	"time"
)

const (
	expireTime = 720 * time.Hour
	secret     = "bitquant"
	RoleUser   = 1
	RoleAdmin  = 2
)

// NewToken role 1 平台用户  2  管理员用户
func NewToken(uid string, role int) (string, error) {
	token := jwt.NewJWT(uid, role, time.Now().Add(expireTime).Unix())
	s, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return s, nil
}
