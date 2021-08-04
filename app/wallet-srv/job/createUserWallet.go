package job

import "wq-fotune-backend/app/wallet-srv/service"

func CreateUserWallet() {
	srv := service.NewWalletService()
	srv.CreateWalletAtRunning()
}
