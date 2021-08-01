package handler

import (
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zhufuyi/pkg/render"
)

func InitEngine() *gin.Engine {
	engine := gin.Default()
	engine.Use(render.InOutLog())
	engine.Use(ginprom.PromMiddleware(nil))
	// register the `/metrics` route.
	engine.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))
	routerGroup := engine.Group("/api/v1")
	InitUserEngine(routerGroup)
	InitCommonEngine(routerGroup)
	InitExOrderEngine(routerGroup)
	InitQuoteEngine(routerGroup)
	InitWalletEngine(routerGroup)
	InitGridEngine(routerGroup)

	pprof.Register(engine)
	return engine
}
