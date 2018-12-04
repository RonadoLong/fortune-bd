package cost


var (

	//支付类型 1 paypal 2 信用卡
	order_pay_type_paypal  = "1"
	order_pay_type_Cen = "2"

	//订单状态：1未确认 2已确认 3退款 4交易成功(已收货) 5交易关闭 6无效
	order_status_unconfirmed = "1"
	order_status_confirmed = "2"
	order_status_refund = "3"
	order_status_successful = "4"
	order_status_close = "5"
	order_status_invalid = "6"

	//发货状态 1未发货 2已发货 3已收货
	shipping_status_not = "1"
	shipping_status_doing = "2"
	shipping_status_done = "3"

	//支付状态 1未支付 2支付中 3已支付
	pay_status_no = "1"
	pay_status_doing = "2"
	pay_status_done = "3"
)
