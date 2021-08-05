package service

import (
	"wq-fotune-backend/app/exchange-srv/internal/biz"
)

const (
	ErrID = "exchangeOrder"
)

type ExOrderService struct {
	exOrderSrv *biz.ExOrderRepo
}

func NewExOrderService() *ExOrderService {
	handler := &ExOrderService{
		exOrderSrv: biz.NewExOrderRepo(),
	}
	return handler
}
