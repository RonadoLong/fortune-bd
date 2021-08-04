package server

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/pkg/render"
	"log"
	"net/http"
	"time"
	"wq-fotune-backend/app/usercenter-srv/router"
	"wq-fotune-backend/libs/logger"
)

func RunHttp(port string) {
	engine := gin.Default()
	engine.Use(render.InOutLog(), gin.Recovery())
	router.Init(engine)

	s := &http.Server{
		Addr:           port,
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	pprof.Register(engine, "/debug")
	logger.Infof("启动服务，监听端口：%d", port)
	if err := s.ListenAndServe(); err != nil {
		log.Println("启动服务失败 ", port)
	}
}
