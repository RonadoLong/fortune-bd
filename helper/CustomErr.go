package helper

import "errors"

var (
	// ErrMissingRealm indicates Realm name is required
	ErrStockNumberNotEnough = errors.New("Stock Not Enough or not exites")

	ErrSkuIdAndGoodsIdWrong = errors.New("SkuId Or GoodsId Wrong")

	ErrCreateStockFlow = errors.New("Err Create StockFlow")

)