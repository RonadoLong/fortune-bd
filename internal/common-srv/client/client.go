package client

import (
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/micro_client"
	pb "wq-fotune-backend/internal/common-srv/proto"
)

func NewCommonClient(etcdAddr string) pb.CommonService {
	service := micro_client.InitBase(
		etcdAddr,
	)
	commonService := pb.NewCommonService(env.COMMON_SRV_NAME, service.Client())
	return commonService
}
