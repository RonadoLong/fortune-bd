package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/service/forward-offer-srv/config"
	"wq-fotune-backend/service/forward-offer-srv/srv"
)

var configPath = flag.String("configPath", "config/conf.yaml", "配置文件路径")

func init() {
	flag.Parse()
	config.Init(*configPath)
}

func main() {
	srv.ListeningOrderService()
	handleSigterm(func() {
		logger.Warn("server was exit")
	})
	_ = http.ListenAndServe("localhost:6070", nil)
}

func handleSigterm(handleFunc func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		s := <-c
		handleFunc()
		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}
	}()
}
