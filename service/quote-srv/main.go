package main

import (
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/global"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/libs/micro_service"
	"wq-fotune-backend/service/quote-srv/cron"
	"wq-fotune-backend/service/quote-srv/handler"
	pb "wq-fotune-backend/service/quote-srv/proto"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
)

func init() {
	go cron.StoreOkexTick()
	go cron.StoreBinanceTick()
	go cron.StoreHuobiTick()
}

func main() {
	service := micro_service.InitBase(
		micro.Name(env.QUOTE_SRV_NAME),
		micro.Version("latest"),
		micro.Action(func(c *cli.Context) error {
			global.InitRedis()
			return nil
		}),
	)
	micro_service.RegisterEtcd(service, env.EtcdAddr)
	if err := pb.RegisterQuoteServiceHandler(service.Server(), handler.NewQuoteHandler()); err != nil {
		logger.Errorf("注册服务错误 %v", err)
		panic(-1)
	}
	if err := service.Run(); err != nil {
		logger.Errorf("%v", err)
		panic(-1)
	}
}
