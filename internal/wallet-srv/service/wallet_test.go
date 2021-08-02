package service

import (
	"testing"
	"wq-fotune-backend/internal/wallet-srv/dao"
)

func TestWalletService_AddIfcBalance(t *testing.T) {
	srv := &WalletService{
		dao:          dao.New(),
		cacheService: nil,
		binance:      nil,
		UserSrv:      nil,
	}
	if err := srv.AddIfcBalance("1273211817757249536", "hhhhhhhh", "register", "", 1.0); err != nil {
		t.Errorf("错误 %v", err.Error())
	}
	t.Log("okok")
}
