package micro_client

import (
	"context"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/logger"
	"time"
	"wq-fotune-backend/libs/registry"
)

func InitBase(etcdAddr string, opts ...micro.Option) micro.Service {
	s := selector.NewSelector(selector.Registry(registry.GetIns(etcdAddr)))
	opts = append(opts,
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		micro.WrapClient(logWrap),
		micro.Selector(s),
	)
	service := micro.NewService(opts...)
	return service
}

// log wrapper logs every time a request is made
type logWrapper struct {
	client.Client
}

func (l *logWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	logger.Infof("[wrapper] client request app: %s method: %s\n", req.Service(), req.Method())
	return l.Client.Call(ctx, req, rsp)
}

// Implements client.Wrapper as logWrapper
func logWrap(c client.Client) client.Client {
	return &logWrapper{c}
}
