package client

import (
	"wq-fotune-backend/libs/env"
	"wq-fotune-backend/libs/micro_client"
	pb "wq-fotune-backend/service/wallet-srv/proto"
)

func NewWalletClient(etcdAddr string) pb.WalletService {
	service := micro_client.InitBase(
		etcdAddr,
	)
	walletService := pb.NewWalletService(env.WALLET_SRV_NAME, service.Client())
	return walletService
}
