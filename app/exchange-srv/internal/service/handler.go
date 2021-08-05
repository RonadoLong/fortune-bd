package service

import (
	"wq-fotune-backend/app/exchange-srv/internal/biz"
)

const (
	ErrID = "exchangeOrder"
)

//var (
//	quoteService pbQuote.QuoteService
//)

//func InitQuoteCli() {
//	quoteService = quoteCli.NewQuoteClient(config.Config.EtcdAddr)
//}

type ExOrderService struct {
	exOrderSrv *biz.ExOrderRepo
}

func NewExOrderHandler() *ExOrderService {
	handler := &ExOrderService{
		exOrderSrv: biz.NewExOrderRepo(),
	}
	return handler
}
