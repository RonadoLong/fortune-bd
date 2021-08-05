package main

import (
	"os"
	"os/signal"
	"syscall"
	"wq-fotune-backend/app/wallet-srv/internal/biz"
	"wq-fotune-backend/app/wallet-srv/server"
	"wq-fotune-backend/libs/logger"
)

func init() {
	biz.NewWalletRepo().CreateWalletAtRunning()
}

const (
	port = "0.0.0.0:9530"
)

func main() {
	go server.RunGrpc()
	go server.RunHttp(port)
	wait()
}

func wait() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt,syscall.SIGTERM)
	<-c
	logger.Info("服务已关闭")
}