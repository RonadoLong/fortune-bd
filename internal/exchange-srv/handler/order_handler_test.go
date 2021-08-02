package handler

import (
	jsoniter "github.com/json-iterator/go"
	"log"
	"testing"
	pb "wq-fotune-backend/internal/exchange-srv/proto"
)

func TestNewExOrderHandler(t *testing.T) {
	signal := &pb.TradeSignal{
		Uid:       "1",
		FileID:    "1",
		SharedID:  "1",
		TradeType: 10,
		Orders: &pb.OrdinaryOrder{
			Side:                  "1",
			OrdType:               "1",
			OrderQty:              60,
			OrderID:               "1",
			DelayerTime:           10,
			TryCount:              10,
			Exchange:              "1",
			Symbol:                "1",
			Price:                 10,
			SlipPrice:             10,
			TradingQtyCoefficient: "1",
		},
	}
	marshalToString, err := jsoniter.MarshalToString(signal)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(marshalToString)
}
