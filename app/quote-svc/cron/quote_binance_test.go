package cron

import (
	"log"
	"testing"
)

func TestStoreBinanceTick(t *testing.T) {
	StoreBinanceTick()
	select {}
}

func TestStringNum(t *testing.T) {
	a := "yutbtc"
	log.Println(len(a))
}

func Test_storeBinanceTick(t *testing.T) {
	storeBinanceTick()
}