package model

import (
	"fmt"
	"sync"
	"testing"
	"time"
	"wq-fotune-backend/app/grid-strategy-srv/util/huobi"

	"github.com/globalsign/mgo/bson"
	"github.com/k0kubun/pp"
	"github.com/zhufuyi/pkg/logger"
)

func TestUpdateNormalGrid(t *testing.T) {
	// 测试参数
	var (
		gsid          = bson.ObjectIdHex("5f0ff0f0c7ef825cd05eb0c9") // 策略id
		orderID       = "47327486515982"                             // 订单id
		tradeTime     = time.Now().UnixNano() / 1e6
		orderStatus   = huobi.OrderStateFilled
		clientOrderID = ""

		exchange     = "huobi"
		symbol       = "btcusdt"
		anchorSymbol = "usdt"
	)

	query := bson.M{"orderID": orderID}
	gtr, err := FindGridTradeRecord(query, bson.M{})
	if err != nil {
		t.Error(err)
		return
	}
	clientOrderID = gtr.ClientOrderID

	query = bson.M{"gsid": gsid}
	gpo, err := FindGridPendingOrder(query, bson.M{})
	if err != nil {
		t.Error(err)
		return
	}

	gp := &GridProcess{
		Exchange:        exchange,
		Symbol:          symbol,
		AnchorSymbol:    anchorSymbol,
		Grids:           gpo.Grids,
		BasisGridNO:     gpo.BasisGridNO,
		ExchangeAccount: &huobi.Account{},
	}

	err = gp.UpdateGridOrder(tradeTime, gtr.Price*gtr.Quantity, orderID, clientOrderID, orderStatus)
	if err != nil {
		t.Error(err)
		return
	}

	pp.Println(gp)
}

func TestGetExchangeKey(t *testing.T) {
	uid := "1266203624401276928"
	apiKey := "26324e15-bgrveg5tmn-6e3a6eb8-09a97"
	key, secret, err := GetExchangeKey(uid, apiKey)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(key, secret)
}

func TestSyncGridStrategy(t *testing.T) {
	if strategyRunInfo == nil {
		strategyRunInfo = new(sync.Map)
	}

	query := bson.M{"isRun": true}
	gss, err := FindGridStrategies(query, bson.M{}, 0, 1000000)
	if err != nil {
		t.Error(err)
		return
	}

	var key string
	count := 0
	for _, v := range gss {
		accessKey, secretKey, account, err := InitExchangeAccount(v.UID, v.Exchange, v.ApiKey)
		if err != nil {
			logger.Error("InitExchangeAccount error", logger.Err(err), logger.String("uid", v.UID), logger.String("exchange", v.Exchange))
			continue
		}

		query = bson.M{"gsid": v.ID}
		gpo, err := FindGridPendingOrder(query, bson.M{})
		if err != nil {
			logger.Error("FindGridPendingOrder error", logger.Err(err), logger.Any("query", query))
			continue
		}

		gp := &GridProcess{
			Gsid:            v.ID.Hex(),
			UID:             v.UID,
			Exchange:        v.Exchange,
			Symbol:          v.Symbol,
			AccessKey:       accessKey,
			SecretKey:       secretKey,
			AnchorSymbol:    v.AnchorSymbol,
			Grids:           gpo.Grids,
			BasisGridNO:     gpo.BasisGridNO,
			ExchangeAccount: account,
		}

		key = StrategyCacheKey(v.UID, v.Exchange, v.Symbol, gp.Gsid)
		SetStrategyCache(key, gp)
		count++

		time.Sleep(40 * time.Millisecond) // 限制初始化速度
	}

	logger.Infof("InitStrategyCache finish, success=%d, total=%d", count, len(gss))

	var (
		uid      = "1593572486"
		exchange = "huobi"
		symbol   = "btcusdt"
		gsid     = "5f1e95efc7ef826e9c8df40e"
	)

	cacheKey := StrategyCacheKey(uid, exchange, symbol, gsid)
	gp, ok := GetStrategyCache(cacheKey)
	if !ok {
		fmt.Println("not found cache", cacheKey)
		return
	}

	err = gp.syncGridStrategy()
	if err != nil {
		t.Error(err)
		return
	}

	select {}
}

func TestGenerateClientOrderID(t *testing.T) {
	exchange := ExchangeBinance
	side := "buy"
	gsid := "5f1fd365c7ef823de44919c1"
	gridNO := 3
	fmt.Println(GenerateClientOrderID(exchange, side, gsid))
	fmt.Println(NewGridClientOrderID(exchange, side, gsid, gridNO))

	// ----------------------------------------------------------------------------------------------
	gp := GridProcess{}
	gp.ListenOrder()
}
