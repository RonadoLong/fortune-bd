package job

import "wq-fotune-backend/internal/wallet-srv/service"

func CreateUserWallet() {
	srv := service.NewWalletService()
	srv.CreateWalletAtRunning()
}
