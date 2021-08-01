package cron

import (
	"log"
	"testing"
	global2 "wq-fotune-backend/libs/global"
	"wq-fotune-backend/service/quote-srv/config"
)

func TestStoreBinanceTick(t *testing.T) {
	config.Init("../../quote-srv/config/conf.yaml")
	global2.InitRedis()
	StoreBinanceTick()
	select {}
}

func TestStringNum(t *testing.T) {
	a := "yutbtc"
	log.Println(len(a))
}
