package server

import (
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zhufuyi/pkg/render"
	"log"
	"net/http"
	"time"
	"wq-fotune-backend/app/quote-srv/router"
	"wq-fotune-backend/libs/logger"
)

func RunHttp(port string) {
	engine := gin.Default()
	engine.Use(render.InOutLog(), gin.Recovery())
	engine.Use(ginprom.PromMiddleware(nil))
	engine.GET("/quote/metrics", ginprom.PromHandler(promhttp.Handler()))
	router.Init(engine)

	s := &http.Server{
		Addr:           port,
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	pprof.Register(engine, "/quote/debug")
	logger.Infof("启动服务，监听端口：%v", port)
	if err := s.ListenAndServe(); err != nil {
		log.Println("启动服务失败 ", port)
	}
}
