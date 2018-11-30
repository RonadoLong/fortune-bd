package main

import (
	"context"
	"github.com/micro/go-micro/server"
	"log"
)

// logWrapper is a handler wrapper
func logWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		log.Printf("[wrapper] server request: %v", req.Method())
		err := fn(ctx, req, rsp)
		return err
	}
}
