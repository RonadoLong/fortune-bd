package client

import (
	pb "wq-fotune-backend/api/usercenter"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/micro_client"
)

func NewUserClient(etcdAddr string) pb.UserService {
	service := micro_client.InitBase(
		etcdAddr,
	)
	userService := pb.NewUserService(env.USER_SRV_NAME, service.Client())
	return userService
}
