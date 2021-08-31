package cron
//
//import (
//	"github.com/robfig/cron"
//	"log"
//	"time"
//	"wq-fotune-backend/app/exchange-srv/internal/biz"
//)
//
//var SrvCron *serviceCron
//
//type serviceCron struct {
//	exOrderSrv *biz.ExOrderRepo
//}
//
//func Init() {
//	SrvCron = &serviceCron{
//		exOrderSrv: biz.NewExOrderRepo(),
//	}
//}
//
//func RunCron() {
//	//Init()
//	//cronTab := "0 0 0 * * *"
//	//timeZone := time.FixedZone("CST", 8*3600)
//	//c := cron.NewWithLocation(timeZone)
//	//err := c.AddFunc(cronTab, func() {
//	//	SrvCron.evaluationDaily() //日线统计
//	//})
//	//if err != nil {
//	//	return
//	//}
//	//c.Start()
//	//log.Println("进入定时任务")
//	//SrvCron.saveRateReturnRank() //收益率排名
//}
