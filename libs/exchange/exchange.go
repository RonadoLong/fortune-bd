package exchange

type ServiceI interface {
	Buy()
	Sell()
	GetAllOrder()
	GetUnfinishOrder()
	SubOrderCallBack()
	SubTraderCallBack()
	GetExchangeName() string
}


func New(exchange string, info Info) ServiceI {
	switch exchange {
	case OKEX:
		return InitOKEX(info)
	}
	return nil
}