package cron

import (
	"testing"
	"wq-fotune-backend/app/exchange-srv/config"
)

func Test_serviceCron_evaluationDaily(t *testing.T) {
	config.Init("../../exchange-srv/config/conf.yaml")
	Init()
	SrvCron.evaluationDaily()
}

//func Test_serviceCron_saveRateReturnRank(t *testing.T) {
//
//	config.Init("../../exchange-srv/config/conf.yaml")
//	Init()
//	SrvCron.saveRateReturnRank()
//}
