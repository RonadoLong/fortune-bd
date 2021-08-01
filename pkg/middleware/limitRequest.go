package middleware

import (
	"github.com/gin-gonic/gin"
	"wq-fotune-backend/libs/jwt"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/pkg/limitReq"
	"wq-fotune-backend/pkg/response"
)

//ipCount := limitReq.GetReqCount(req.Phone)
//if ipCount > 0 {
//return response.NewLoginReqFreqErrMsg(errID)
//}
//if err := limitReq.SetReqCount(req.Phone, ipCount+1); err != nil {
//logger.Errorf("limitReq.SetReqCount %v %s", err, req.Phone)
//return response.NewInternalServerErrMsg(errID)
//}

func LimitRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		jwtP := claims.(*jwt.JWTPayload)
		ipCount := limitReq.GetReqCount(jwtP.UserID)
		if ipCount > 0 {
			response.NewErrorMaxReq(c, "请求过于频繁", nil)
			c.Abort()
			return
		}
		if err := limitReq.SetReqCountWithTimeOut(jwtP.UserID, ipCount+1, 5); err != nil {
			logger.Errorf("limitReq.SetReqCount %v %s", err, jwtP.UserID)
			response.NewInternalServerErr(c, nil)
			c.Abort()
			return
		}
	}
}
