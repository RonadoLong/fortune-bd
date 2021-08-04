package router

import "github.com/gin-gonic/gin"

func Init(eg *gin.Engine) {
	apiV1(eg.Group("/v1"))
}