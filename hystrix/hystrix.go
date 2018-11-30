package hystrix

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/client"
	"log"
	"net"
	"net/http"

	"context"
)

type clientWrapper struct {
	client.Client
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return hystrix.Do(req.Service()+"."+req.Method(), func() error {
		return c.Client.Call(ctx, req, rsp, opts...)
	}, func(err error) error {
		// write your fallback logic
		log.Printf("fallback error!!!!!  %v", err)
		return err
	})
}

func NewClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &clientWrapper{c}
	}
}

func Configure(names []string) {
	// overwrite default configure
	config := hystrix.CommandConfig{
		Timeout:               2000,
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 25,
	}
	configs := make(map[string]hystrix.CommandConfig)
	for _, name := range names {
		configs[name] = config
	}
	hystrix.Configure(configs)

	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "8181"), hystrixStreamHandler)
	log.Println("Launched hystrixStreamHandler at 8181")
}
