package exchange_clientI

import "wq-fotune-backend/pkg/goex"

type ClientI interface {
	GetAccountSpot() (*goex.Account, error)
	GetAccountSwap() (*goex.Account, error)
	CheckIfApiValid() error
}
