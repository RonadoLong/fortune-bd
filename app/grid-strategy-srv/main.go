package grid_strategy_srv

import (
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/app/grid-strategy-srv/model"
	"wq-fotune-backend/app/grid-strategy-srv/router"

	"github.com/zhufuyi/pkg/logger"
	"github.com/zhufuyi/pkg/mongo"
)

func InitMain(engine *gin.RouterGroup) {
	initServer()
	router.Init(engine)
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

//func runWebServer() {
//	if config.IsProd() {
//		gin.SetMode(gin.ReleaseMode)
//	}
//
//	engine := gin.Default()
//	engine.Use(render.InOutLog())
//	router.Init(engine)
//
//	if config.IsEnableProfile() {
//		pprof.Register(engine, "/goprofile/"+config.GetAppName())
//	}
//
//	logger.Infof("启动服务，监听端口：%d", config.GetAppPort())
//	engine.Run(fmt.Sprintf(":%d", config.GetAppPort()))
//}

func wait() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("服务已关闭")
}
