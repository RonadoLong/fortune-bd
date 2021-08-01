package middleware

import (
	"github.com/gin-gonic/gin"
	"wq-fotune-backend/libs/jwt"
	"wq-fotune-backend/pkg/response"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			response.NewErrorParam(c, "请登录", nil)
			c.Abort()
			return
		}
		//log.Print("get token: ", token)
		// parseToken 解析token包含的信息
		j, err := parseToken(token)
		if err != nil || j == nil {
			response.NewErrorParam(c, "解析token错误", nil)
			c.Abort()
			return
		}
		if j.Expired() {
			response.NewErrorParam(c, "token已经过期", nil)
			c.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("claims", j.Payload)
	}
}

func parseToken(tokenString string) (*jwt.JWT, error) {
	j := &jwt.JWT{}
	if err := j.JWTParse(tokenString, secret); err != nil {
		return nil, err
	}
	return j, nil
}
