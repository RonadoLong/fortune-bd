package router

import "github.com/gin-gonic/gin"

// Init routers
func Init(engine *gin.Engine) {
	v1api(engine.Group("/quote/v1"))
}
