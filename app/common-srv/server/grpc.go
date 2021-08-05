package server

import (
	"fmt"
	pb "wq-fotune-backend/api/common"
	"wq-fotune-backend/app/common-srv/internal/service"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/global"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/libs/micro_service"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
)

func RunGrpc() {
	ms := micro_service.InitBase(
		micro.Name(env.COMMON_SRV_NAME),
		micro.Version("latest"),
		micro.Action(func(c *cli.Context) error {
			global.InitRedis()
			return nil
		}),
	)
	micro_service.RegisterEtcd(ms, env.EtcdAddr)
	if err := pb.RegisterCommonHandler(ms.Server(), service.NewCommonService()); err != nil {
		logger.Panic(fmt.Sprintf("注册服务错误 %v", err))
	}
	if err := ms.Run(); err != nil {
		logger.Panic(fmt.Sprintf("启动grpc服务失败: %+v", err))
	}
}
