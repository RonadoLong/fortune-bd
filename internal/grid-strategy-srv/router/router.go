package router

import (

	"github.com/gin-gonic/gin"
)

// Init routers
func Init(engine *gin.RouterGroup) {
	v1api(engine.Group("/grid"))
}
