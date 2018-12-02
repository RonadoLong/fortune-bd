package helper


var (
	// 状态：1未确认 2已确认 3退款 4交易成功(已收货) 5交易关闭 6无效
	// order status
	ORDERSTATUS_NO_CONFIRM = 1
	ORDERSTATUS_CONFIRMED = 2
	ORDERSTATUS_RETURN = 3
	ORDERSTATUS_SUCCESSFULL = 4
	ORDERSTATUS_ClOSE = 5
	ORDERSTATUS_INVAILD = 6


	//pay status 0未支付 1支付中 2已支付
	PAYSTATUS_NO_PAY = 1
	PAYSTATUS_PAIED = 2


	//ShippingStatus 发货状态 0未发货 1已发货 2已收货
	SHIPPINGSTATUS_NO_SHIP = 1
	SHIPPINGSTATUS_SHIPED = 2
	SHIPPINGSTATUS_DONE = 3

	//支付类型 1 pay pal  2 信用卡

)