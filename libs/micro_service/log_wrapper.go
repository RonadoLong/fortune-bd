package micro_service

import (
	"context"
	"github.com/micro/go-micro/v2/server"
	"wq-fotune-backend/libs/logger"
)

func LogWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		logger.Infof("[wrapper] server request method: %v req: %+v", req.Method(), req.Body())
		err := fn(ctx, req, rsp)
		//logger.Infof("[wrapper] server response: %+v", rsp)
		return err
	}
}
