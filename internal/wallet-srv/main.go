package main

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/libs/micro_service"
	"wq-fotune-backend/pkg/redis"
	"wq-fotune-backend/internal/wallet-srv/handler"
	"wq-fotune-backend/internal/wallet-srv/job"
	pd "wq-fotune-backend/internal/wallet-srv/proto"
)

func main() {
	service := micro_service.InitBase(
		micro.Name(env.WALLET_SRV_NAME),
		micro.Version("latest"),
		micro.Action(func(c *cli.Context) error {
			redis.InitRedis(env.RedisAddr, env.RedisPWD)
			return nil
		}),
	)
	micro_service.RegisterEtcd(service, env.EtcdAddr)
	if err := pd.RegisterWalletServiceHandler(service.Server(), handler.NewWalletHandler()); err != nil {
		logger.Errorf("注册服务错误 %v", err)
		panic(-1)
	}
	//计划任务
	go job.CreateUserWallet()
	if err := service.Run(); err != nil {
		logger.Errorf("%v", err)
		panic(-1)
	}
}
