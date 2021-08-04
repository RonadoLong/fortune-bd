package client

import (
	pb "wq-fotune-backend/api/quote"
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/micro_client"
)

func NewQuoteClient(etcdAddr string) pb.QuoteService {
	service := micro_client.InitBase(
		etcdAddr,
	)
	quoteService := pb.NewQuoteService(env.QUOTE_SRV_NAME, service.Client())
	return quoteService
}
