package cron

import (
	"github.com/robfig/cron"
	"log"
	"time"
	"wq-fotune-backend/app/exchange-srv/service"
)

var SrvCron *serviceCron

type serviceCron struct {
	exOrderSrv *service.ExOrderService
}

func Init() {
	SrvCron = &serviceCron{
		exOrderSrv: service.NewExOrderService(),
	}
}

func RunCron() {
	Init()
	cronTab := "0 0 0 * * *"
	timeZone := time.FixedZone("CST", 8*3600)
	c := cron.NewWithLocation(timeZone)
	c.AddFunc(cronTab, func() {
		SrvCron.evaluationDaily() //日线统计
	})
	c.Start()
	log.Println("进入定时任务")
	SrvCron.saveRateReturnRank() //收益率排名
}
