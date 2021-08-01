package cron

import (
	"time"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/pkg/utils"
)

func (s *serviceCron) evaluationDaily() {
	strategys := s.exOrderSrv.GetUserStrategyOfRun()
	for _, strategy := range strategys {
		wqProfit, err := s.exOrderSrv.GetProfitByStrID(strategy.UID, strategy.ID)
		if err != nil {
			logger.Warnf("日线统计错误，没有找到实时统计信息 策略id %s", strategy.ID)
			continue
		}
		runDay := int(time.Now().Sub(strategy.CreatedAt).Hours() / 24)
		runDay += 1

		if err := s.exOrderSrv.CreateWqProfitDaily(wqProfit); err != nil {
			logger.Warnf("创建日线统计失败 策略id %s err %v", strategy.ID, err)
			continue
		}
		//更新一下实时统计的年化率
		rateReturnYear := wqProfit.RateReturn / (float64(runDay) / 365)
		wqProfit.RateReturnYear = utils.Keep2Decimal(rateReturnYear)
		wqProfit.UpdatedAt = time.Now()
		if err := s.exOrderSrv.UpdateProfit(strategy.ID, wqProfit); err != nil {
			logger.Warnf("更新实时年化率失败 策略id%s err %v", strategy.ID, err)
		}
		logger.Infof("策略id %s 日线统计成功", strategy.ID)
	}
}
