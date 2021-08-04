package router

import (

	"github.com/gin-gonic/gin"
	v1 "wq-fotune-backend/app/grid-strategy-srv/router/v1"
)

func v1api(rg *gin.RouterGroup) {
	rg.POST("/strategy/grid/startup", v1.StartupGridStrategy) // 启动网格交易策略
	rg.POST("/strategy/grid/stop", v1.StopGridStrategy)       // 停止网格交易策略
	rg.POST("/strategy/grid/update", v1.UpdateGridStrategy)   // 更新网格交易策略

	rg.POST("/strategy/grid/calculateMoney", v1.CalculateMoney)      // 计算整个网格交易所需总资金
	rg.GET("/strategy/grid/calculateParams", v1.CalculateGridParams) // 通过参数计算网格参数
	rg.GET("/strategy/grid/auto/minMoney", v1.GetMinMoney)           // 根据交易所和品种获取最小投入资金
	rg.GET("/strategy/grid/auto/gridParams", v1.GetGridParams)       // 根据投入资金自动生成网格参数

	rg.GET("/strategy/bigGrid/auto/gridParams", v1.GetBigGridParams) // 获取大网格参数

	rg.GET("/strategy/reverseGrid/minMoney", v1.GetReverseMinMoney)                // 根据交易所和品种获取反向网格最小投入资金
	rg.GET("/strategy/reverseGrid/calculateParams", v1.CalculateReverseGridParams) // 根据投入资金自动生成反向网格参数
	rg.POST("/strategy/reverseGrid/startup", v1.StartupReverseGridStrategy)        // 获取反向网格参数
	rg.POST("/strategy/reverseGrid/stop", v1.StopReverseGridStrategy)              // 停止反向网格交易策略

	rg.GET("/strategy/list", v1.ListRunningStrategies)   // 获取正在运行的策略列表
	rg.GET("/strategy/detail/:id", v1.GetStrategyDetail) // 获取策略详情
	rg.GET("/strategy/simple/:id", v1.GetStrategySimple) // 获取策略参数和利润

	rg.GET("/strategy/total", v1.GetStrategyTotal) // 获取一共有多少个网格策略在运行

	rg.POST("/exchange/limit", v1.SyncExchangeLimit)                  // 同步交易所品种数量和交易额限制信息
	rg.GET("/exchange/limits", v1.GetExchangeLimits)                  // 获取交易所品种限制列表
	rg.GET("/strategy/currencyPair", v1.GetCurrencyPair)              // 根据品种名称获取品种的交易对
	rg.PATCH("/exchange/limit/currencyPairs", v1.UpdateCurrencyPairs) // 根据锚定币更新交易对

	rg.POST("/strategy/type", v1.CreateStrategyType)       // 添加策略类型
	rg.DELETE("/strategy/type/:id", v1.DeleteStrategyType) // 删除策略类型
	rg.GET("/strategy/types", v1.ListStrategyTypes)        // 获取策略类型列表

	rg.GET("/user/strategyCount/:uid", v1.GetStrategyCount) // 查看用户当前正在运行策略数量

}
