package main

import (
	"fortune-bd/app/grid-strategy-svc/model"
	"fortune-bd/app/grid-strategy-svc/server"
	"fortune-bd/libs/env"
	"os"
	"os/signal"


	"github.com/zhufuyi/pkg/logger"
	"github.com/zhufuyi/pkg/mongo"
)

const (
	port = "0.0.0.0:9530"
)


func main() {
	initServer()
	go server.RunHttp(port)
	wait()
}

func initServer() {
	// 配置初始化
	logger.Infof("init mongodb ...... %s", env.MongoAddr)
	err := mongo.InitializeMongodb(env.MongoAddr)
	if err != nil {
		logger.Fatal("init mongodb failed.", logger.Err(err), logger.Any("mongodb", env.MongoAddr))
	}
	initCache()
}

func initCache() {
	// 初始化限制值
	err := model.InitExchangeLimitCache()
	if err != nil {
		logger.Error("InitExchangeLimitCache failed", logger.Err(err))
	}
	// 加载策略类型
	err = model.InitStrategyTypeCache()
	if err != nil {
		logger.Error("InitStrategyTypeCache failed", logger.Err(err))
	}

	// 加载策略运行信息
	logger.Info("init strategy run info ......")
	err = model.InitStrategyCache()
	if err != nil {
		logger.Error("InitStrategyCache failed", logger.Err(err))
	}
}

func wait() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("服务已关闭")
}
