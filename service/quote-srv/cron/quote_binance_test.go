package cron

import (
	"log"
	"testing"
	"wq-fotune-backend/libs/global"
)

func TestStoreBinanceTick(t *testing.T) {
	global.InitRedis()
	StoreBinanceTick()
	select {}
}

func TestStringNum(t *testing.T) {
	a := "yutbtc"
	log.Println(len(a))
}
