package cron

import (
	"time"
)

func (s *serviceCron) saveRateReturnRank() {
	d := time.Second * 120
	t := time.NewTicker(d)
	defer t.Stop()

	for {
		<-t.C
		//定时保存排名数据
		s.exOrderSrv.CacheRateReturn()
		s.exOrderSrv.CacheRateReturnYear()
	}
}
