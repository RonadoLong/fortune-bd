package main

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/global"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/libs/micro_service"
	"wq-fotune-backend/app/common-srv/handler"
	pd "wq-fotune-backend/app/common-srv/proto"
)


func main() {
	service := micro_service.InitBase(
		micro.Name(env.COMMON_SRV_NAME),
		micro.Version("latest"),
		micro.Action(func(c *cli.Context) error {
			global.InitRedis()
			return nil
		}),
	)
	micro_service.RegisterEtcd(service, env.EtcdAddr)
	if err := pd.RegisterCommonHandler(service.Server(), handler.NewCommonHandler()); err != nil {
		logger.Errorf("注册服务错误 %v", err)
		panic(-1)
	}
	if err := service.Run(); err != nil {
		logger.Errorf("%v", err)
		panic(-1)
	}
}
