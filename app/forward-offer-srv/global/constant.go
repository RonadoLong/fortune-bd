package global

const (
	Exchange = "okex"
	ErrorStr = "error"

	SellType = "sell"
	BuyType  = "buy"

	PositionType   = "position"
	OrderType      = "order"
	ExecutionType  = "execution"
	InstrumentType = "instrument"

	CreateType  = "create"  // 普通退单类型
	CancelType  = "cancel"  // 取消订单类型
	AutoAddType = "autoAdd" // 自动追加订单

	MarketType = "Market"
	LimitType  = "Limit"

	RespSuccess   = 200
	RespError     = 1500 // 普通错误
	ReqOrderError = 1400 // 发单失败
	ReqTryError   = 1503 // 重试发单中
	ReqMaxError   = 1529 // 重试发单最大次数

	PositionKey = "POSITION:"

	Topic     = "trade:okex:reliable" // 自动追单
	Queue     = "trade:okex:request"  // 请求订单的队列
	RespQueue = "trade:okex:response" // 发送回调信息的队列

	InstrumentCacheKey  = "okex:instrument"  // 合约KEY
	StrategyIDOfApiKey  = "okex:apiKey:"     // 缓存 strategyID对应的账户信息
	OrderIDOfStrategyID = "okex:strategyID:" // 缓存 OrderID对应的strategyID

	StrategyEventKey = "StrategyEvent:okex" // 关闭策略事件

	// 交易状态
	FilledStatus         = "filled"         // 成交
	CanceledStatus       = "cancelled"      // 取消
	NewStatus            = "new"            // 委托成功
	RejectedStatus       = "rejected"       // 拒绝类型
	MaxTryCanceledStatus = "maxTryCanceled" // 取消

	DealStatusPending = "pending"
	DealStatusFinish  = "finish"

	//2001 交易账号校验失败
	CollectAccountErrorCode = 2001
	//2002 资金不足
	CollectPriceErrorCode = 2002
	//2003 平仓数量不足
	CollectOrderQtyErrorCode = 2003
	//2005 交易所繁忙
	CollectExchangeBusyErrorCode = 2005
	//2006 请求参数有误
	CollectOrderParamsErrorCode = 2006
	//2007 订单超过限制提交次数
	CollectOrderLimiterErrorCode = 2007
)
