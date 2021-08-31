package dao

import (
	"testing"
	"wq-fotune-backend/libs/cache"
)

func TestDao_GetUserStrategy(t *testing.T) {
	mgoClient := cache.InitMongo("mongodb://wquant:wqabc123@47.57.161.25:38028/ifortune")
	d := &Dao{
		mongo: mgoClient,
	}
	got, err := d.GetUserStrategy("1266203624401276928", "5f115c44a630290001680358")
	if err != nil {
		t.Errorf("查询失败 %v", err)
	}
	t.Logf("%+v", got)
}

func TestDao_GetUserStrategyOfRun(t *testing.T) {
	mgoClient := cache.InitMongo("mongodb://wquant:wqabc123@47.57.161.25:38028/ifortune")
	d := &Dao{
		mongo: mgoClient,
	}
	got := d.GetUserStrategyOfRun(nil)
	t.Logf("kkkk%+v", got)
}

func TestDao_GetUserStrategyByApiKey(t *testing.T) {
	mgoClient := cache.InitMongo("mongodb://wquant:wqabc123@47.57.161.25:38028/ifortune")
	d := &Dao{
		mongo: mgoClient,
	}
	got := d.GetUserStrategyByApiKey("1266203624401276928", "20b55e71-e4b3ff5b-48d0c8f4-dab4c45e6f")
	t.Logf("kkkk%+v", got)
}

func TestDao_SetUserAllStrategyApi(t *testing.T) {
	mgoClient := cache.InitMongo("mongodb://wquant:wqabc123@47.57.161.25:38028/ifortune")
	d := &Dao{
		mongo: mgoClient,
	}
	if err := d.SetUserAllStrategyApi("1266203624401276928", "20b55e71-e4b3ff5b-48d0c8f4-dab4c45e6f", "huobi"); err != nil {
		t.Errorf("%v", err)
	}

}
