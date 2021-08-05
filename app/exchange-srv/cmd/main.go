package main

import (
	"os"
	"os/signal"
	"syscall"
	"wq-fotune-backend/app/exchange-srv/cron"
	"wq-fotune-backend/app/exchange-srv/server"
	"wq-fotune-backend/libs/logger"
)

const (
	port = "0.0.0.0:9530"
)

func main() {
	go server.RunGrpc()
	go server.RunHttp(port)
	//日线统计
	go cron.RunCron()
	wait()
}

func wait() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt,syscall.SIGTERM)
	<-c
	logger.Info("服务已关闭")
}