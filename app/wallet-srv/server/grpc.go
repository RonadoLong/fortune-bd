package server

import (
	"fmt"
	pb "wq-fotune-backend/api/wallet"
	"wq-fotune-backend/app/wallet-srv/internal/service"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/libs/micro_service"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
)

func RunGrpc() {
	app := micro_service.InitBase(
		micro.Name(env.WALLET_SRV_NAME),
		micro.Version("latest"),
		micro.Action(func(c *cli.Context) error {
			return nil
		}),
	)
	micro_service.RegisterEtcd(app, env.EtcdAddr)
	if err := pb.RegisterWalletServiceHandler(app.Server(), service.NewWalletService()); err != nil {
		logger.Panic(fmt.Sprintf("注册服务错误 %v", err))
	}
	if err := app.Run(); err != nil {
		logger.Panic(fmt.Sprintf("启动grpc服务失败: %+v", err))
	}
}
