package client

import (
	pb "wq-fotune-backend/api/exchange"
	pbQuote "wq-fotune-backend/api/quote"
	pbUser "wq-fotune-backend/api/usercenter"
	quoteCli "wq-fotune-backend/app/quote-srv/client"
	userCli "wq-fotune-backend/app/usercenter-srv/client"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/micro_client"
)

func NewExOrderClient(etcdAddr string) pb.ExOrderService {
	service := micro_client.InitBase(
		etcdAddr,
	)
	exOrderService := pb.NewExOrderService(env.EXCHANGE_SRV_NAME, service.Client())
	return exOrderService
}

func NewForwardOfferClient(etcdAddr string) pb.ForwardOfferService {
	service := micro_client.InitBase(
		etcdAddr,
	)
	exOrderService := pb.NewForwardOfferService(env.EXCHANGE_SRV_NAME, service.Client())
	return exOrderService
}

func GetQuoteService() pbQuote.QuoteService {
	return quoteCli.NewQuoteClient(env.EtcdAddr)
}

func GetUserService() pbUser.UserService {
	return userCli.NewUserClient(env.EtcdAddr)
}

