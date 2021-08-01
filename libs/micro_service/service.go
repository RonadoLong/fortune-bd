package micro_service

import (
	"github.com/micro/go-micro/v2"
	"time"
	"wq-fotune-backend/libs/registry"
)

func InitBase(opts ...micro.Option) micro.Service {
	service := micro.NewService(
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		micro.WrapHandler(LogWrapper),
		// TODO 监控完善
		//micro.WrapHandler(prometheus.NewHandlerWrapperer()),
	)
	service.Init(opts...)
	return service
}

func RegisterEtcd(service micro.Service, etcdAddr string) {
	service.Init(micro.Registry(registry.GetIns(etcdAddr)))
}
