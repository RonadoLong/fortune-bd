package handler

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"math/rand"
	"testing"
	"time"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/pkg/snowflake"
	"wq-fotune-backend/app/exchange-srv/client"
	fotune_srv_exchange "wq-fotune-backend/app/exchange-srv/proto"
)

func TestExOrderHandler_Evaluation(t *testing.T) {
	orderClient := client.NewExOrderClient("127.0.0.0.1:2379")
	tradeAt, _ := ptypes.TimestampProto(time.Now())

	rand.Seed(time.Now().Unix())
	reasons := []string{
		//"sell",
		//"sell",
		"buy",
		//"sell",
	}
	prisetest := []string{
		//"200",
		"210.5",
		//"230",
		//"230",
	}
	for i := 1; i <= 4; i++ {
		orderID := snowflake.SNode.Generate().String()
		orderID = orderID + "trade"

		n := rand.Int() % len(reasons)

		req := &fotune_srv_exchange.TradeReq{
			TradeId:    orderID,
			UserId:     "1273211817757249536",
			ApiKey:     "62dd3fc3-03d6-4e7e-a4fc-ce238df6aa40",
			OrderId:    orderID,
			StrategyId: "6672302101017661440",
			Direction:  reasons[n],
			Volume:     "10",
			Commission: "0.5",
			Unit:       "eth",
			Symbol:     "ETH-USDT-SWAP",
			Price:      prisetest[n],
			TradeAt:    tradeAt,
			Exchange:   "okex",
		}
		_, err := orderClient.Evaluation(context.Background(), req)
		if err != nil {
			logger.Warnf("%+v", err)
		}
	}

}
