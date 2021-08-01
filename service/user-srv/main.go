package main

import (

	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/libs/micro_service"
	"wq-fotune-backend/pkg/redis"
	"wq-fotune-backend/service/user-srv/handler"
	pd "wq-fotune-backend/service/user-srv/proto"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
)


func main() {
	service := micro_service.InitBase(
		micro.Name(env.USER_SRV_NAME),
		micro.Version("latest"),
		micro.Action(func(c *cli.Context) error {
			redis.InitRedis(env.RedisAddr, env.RedisPWD)
			return nil
		}),
	)
	micro_service.RegisterEtcd(service, env.EtcdAddr)
	if err := pd.RegisterUserHandler(service.Server(), handler.NewUserHandler()); err != nil {
		logger.Errorf("注册服务错误 %v", err)
		panic(-1)
	}
	if err := service.Run(); err != nil {
		logger.Errorf("%v", err)
		panic(-1)
	}
}
