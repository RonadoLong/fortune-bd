package client

import (
	pb "wq-fotune-backend/api/exchange"
	pbQuote "wq-fotune-backend/api/quote"
	pbUser "wq-fotune-backend/api/usercenter"
	pbWallet "wq-fotune-backend/api/wallet"
	quoteCli "wq-fotune-backend/app/quote-srv/client"
	userCli "wq-fotune-backend/app/usercenter-srv/client"
	walletCli "wq-fotune-backend/app/wallet-srv/client"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/micro_client"
)

func NewExOrderClient(etcdAddr string) pb.ExOrderService {
	service := micro_client.InitBase(
		etcdAddr,
		//micro.Name("exchange-srv.client"),
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

var (
	QuoteService  pbQuote.QuoteService
	UserService   pbUser.UserService
	WalletService pbWallet.WalletService
)

func InitService() {
	QuoteService = quoteCli.NewQuoteClient(env.EtcdAddr)
	UserService = userCli.NewUserClient(env.EtcdAddr)
	WalletService = walletCli.NewWalletClient(env.EtcdAddr)
}
