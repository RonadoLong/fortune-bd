package service

import (
	"wq-fotune-backend/app/exchange-srv/internal/biz"
)

const (
	ErrID = "exchangeOrder"
)

type ExOrderService struct {
	ExOrderSrv *biz.ExOrderRepo
}

func NewExOrderService() *ExOrderService {
	handler := &ExOrderService{
		ExOrderSrv: biz.NewExOrderRepo(),
	}
	return handler
}
