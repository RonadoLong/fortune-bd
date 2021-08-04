package main

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/libs/micro_service"
	"wq-fotune-backend/app/exchange-srv/client"
	"wq-fotune-backend/app/exchange-srv/cron"
	"wq-fotune-backend/app/exchange-srv/handler"
	pd "wq-fotune-backend/app/exchange-srv/proto"
)


func main() {
	service := micro_service.InitBase(
		micro.Name(env.EXCHANGE_ORDER_SRV_NAME),
		micro.Version("latest"),
		micro.Action(func(c *cli.Context) error {
			client.InitService()
			return nil
		}),
	)

	micro_service.RegisterEtcd(service, env.EtcdAddr)
	if err := pd.RegisterExOrderHandler(service.Server(), handler.NewExOrderHandler()); err != nil {
		logger.Errorf("注册服务错误 %v", err)
		panic(-1)
	}
	if err := pd.RegisterForwardOfferHandler(service.Server(), handler.NewForwardOfferHandle()); err != nil {
		logger.Errorf("注册服务错误 %v", err)
		panic(-1)
	}
	logger.Infof("%s app run successful", env.EXCHANGE_ORDER_SRV_NAME)
	go cron.RunCron() //日线统计
	if err := service.Run(); err != nil {
		logger.Errorf("%v", err)
		panic(-1)
	}
}
