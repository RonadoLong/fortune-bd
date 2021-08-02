package dao

import (
	"testing"
	"wq-fotune-backend/libs/mongoClient"
	"wq-fotune-backend/pkg/dbclient"
)

func TestDao_GetOrderRecord(t *testing.T) {
	mgoClient, err := mongoClient.InitMongo("mongodb://wquant:wqabc123@47.57.161.25:38028/ifortune")
	if err != nil {
		t.Errorf("")
		return
	}
	d := &Dao{
		mongo: mgoClient,
	}
	got, err := d.GetOrderRecord("63480760852636")
	if err != nil {
		t.Errorf("查询失败 %v", err)
	}
	t.Logf("%+v", got)
}

func TestDao_GetWqProfitRank(t *testing.T) {
	db := dbclient.NewDB("root:WQabc123@tcp(47.57.169.103:13306)/wq_fotune?charset=utf8mb4&parseTime=True&loc=Local")
	d := &Dao{
		db:    db,
		mongo: nil,
	}
	rank := d.GetWqProfitRank()
	for _, v := range rank {
		t.Logf("%+v", v)
	}
}
