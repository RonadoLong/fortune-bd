package handler

import (
	"wq-fotune-backend/service/exchange-srv/service"
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

type ExOrderHandler struct {
	exOrderSrv *service.ExOrderService
}

func NewExOrderHandler() *ExOrderHandler {
	handler := &ExOrderHandler{
		exOrderSrv: service.NewExOrderService(),
	}
	return handler
}
