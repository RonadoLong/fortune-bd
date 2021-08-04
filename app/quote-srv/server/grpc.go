package server

import (
	"fmt"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	pb "wq-fotune-backend/api/quote"
	"wq-fotune-backend/app/quote-srv/internal/service"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/global"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/libs/micro_service"
)

func RunGrpc() {
	ms := micro_service.InitBase(
		micro.Name(env.QUOTE_SRV_NAME),
		micro.Version("latest"),
		micro.Action(func(c *cli.Context) error {
			global.InitRedis()
			return nil
		}),
	)
	micro_service.RegisterEtcd(ms, env.EtcdAddr)
	if err := pb.RegisterQuoteServiceHandler(ms.Server(), service.NewQuoteService()); err != nil {
		logger.Panic(fmt.Sprintf("注册服务错误 %v", err))
	}
	if err := ms.Run(); err != nil {
		logger.Panic(fmt.Sprintf("启动grpc服务失败: %+v", err))
	}
}