package client

import (
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/micro_client"
	pb "wq-fotune-backend/app/quote-srv/proto"
)

func NewQuoteClient(etcdAddr string) pb.QuoteService {
	service := micro_client.InitBase(
		etcdAddr,
	)
	quoteService := pb.NewQuoteService(env.QUOTE_SRV_NAME, service.Client())
	return quoteService
}
