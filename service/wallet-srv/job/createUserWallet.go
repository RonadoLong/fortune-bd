package job

import "wq-fotune-backend/service/wallet-srv/service"

func CreateUserWallet() {
	srv := service.NewWalletService()
	srv.CreateWalletAtRunning()
}
