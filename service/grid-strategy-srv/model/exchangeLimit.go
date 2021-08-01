package model

import (
	"strings"
	"sync"
	"time"
	"wq-fotune-backend/service/grid-strategy-srv/util/goex"

	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/zhufuyi/pkg/logger"
	"github.com/zhufuyi/pkg/mongo"
)

// ExchangeLimitCollectionName 表名
const ExchangeLimitCollectionName = "exchangeLimit"

// ExchangeLimit 交易所限制
type ExchangeLimit struct {
	mongo.PublicFields `bson:",inline"`

	Exchange string `json:"exchange" bson:"exchange"` // 交易所
	Symbol   string `json:"symbol" bson:"symbol"`     // 品种名称
	ES       string `json:"es" bson:"es"`             // 唯一标识，由交易所和品种组成

	Currency       string `json:"currency" bson:"currency"`             // 交易品种
	AnchorCurrency string `json:"anchorCurrency" bson:"anchorCurrency"` // 锚定币

	QuantityLimit     float64 `json:"quantityLimit" bson:"quantityLimit"`         // 买卖最小值数量
	VolumeLimit       float64 `json:"volumeLimit" bson:"volumeLimit"`             // 买卖最小金额
	PricePrecision    int     `json:"pricePrecision" bson:"pricePrecision"`       // 价格精度
	QuantityPrecision int     `json:"quantityPrecision" bson:"quantityPrecision"` // 数量精度
	LeverageRatio     float64 `json:"leverageRatio" bson:"leverageRatio"`         // 最大杠杆比例
}

// Insert 插入一条新的记录
func (object *ExchangeLimit) Insert() (err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()
	object.SetFieldsValue()

	return mconn.Insert(ExchangeLimitCollectionName, object)
}

// FindExchangeLimit 获取单条记录
func FindExchangeLimit(selector bson.M, field bson.M) (*ExchangeLimit, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	object := &ExchangeLimit{}
	return object, mconn.FindOne(ExchangeLimitCollectionName, object, selector, field)
}

// FindExchangeLimits 获取多条记录
func FindExchangeLimits(selector bson.M, field bson.M, page int, limit int, sort ...string) ([]*ExchangeLimit, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	// 默认从第一页开始
	if page < 0 {
		page = 0
	}
	if len(sort) == 0 || sort[0] == "" {
		sort = []string{"-_id"}
	}

	objects := []*ExchangeLimit{}
	return objects, mconn.FindAll(ExchangeLimitCollectionName, &objects, selector, field, page, limit, sort...)
}

// UpdateExchangeLimit 更新单条记录
func UpdateExchangeLimit(selector, update bson.M) (err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateOne(ExchangeLimitCollectionName, selector, mongo.UpdatedTime(update))
}

// UpdateExchangeLimits 更新多条记录
func UpdateExchangeLimits(selector, update bson.M) (n int, err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateAll(ExchangeLimitCollectionName, selector, mongo.UpdatedTime(update))
}

// FindAndModifyExchangeLimit 更新并返回最新记录
func FindAndModifyExchangeLimit(selector bson.M, update bson.M) (*ExchangeLimit, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	object := &ExchangeLimit{}
	return object, mconn.FindAndModify(ExchangeLimitCollectionName, object, selector, mongo.UpdatedTime(update))
}

// CountExchangeLimits 统计数量，不包括删除记录
func CountExchangeLimits(selector bson.M) (n int, err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.Count(ExchangeLimitCollectionName, mongo.ExcludeDeleted(selector))
}

// DeleteExchangeLimit 删除记录
func DeleteExchangeLimit(selector bson.M) (int, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateAll(ExchangeLimitCollectionName, selector, mongo.DeletedTime(bson.M{}))
}

func (e *ExchangeLimit) ToLimitValues() *LimitValues {
	return &LimitValues{
		Key:               e.ES,
		QuantityLimit:     e.QuantityLimit,
		VolumeLimit:       e.VolumeLimit,
		PricePrecision:    e.PricePrecision,
		QuantityPrecision: e.QuantityPrecision,
		LeverageRatio:     e.LeverageRatio,
	}
}

// -------------------------------------------------------------------------------------------------

// ExchangeLimitsMap 交易所限制值字典，交易所.品种作为key，对象LimitValues作为值
var ExchangeLimitsMap = new(sync.Map)

// LimitValues 限制值
type LimitValues struct {
	Key               string  `json:"key"`               // 字典的key
	QuantityLimit     float64 `json:"buyLimit"`          // 买卖最小值数量
	VolumeLimit       float64 `json:"volumeLimit"`       // 买卖最小金额
	PricePrecision    int     `json:"pricePrecision"`    // 价格精度
	QuantityPrecision int     `json:"quantityPrecision"` // 数量精度
	LeverageRatio     float64 `json:"leverageRatio"`     // 最大杠杆比例
}

// GetVolumePrecision 获取最小交易金额小数点位数(只针对币本位)
func (l *LimitValues) GetVolumePrecision() int {
	ss := strings.Split(fmt.Sprintf("%v", l.VolumeLimit), ".")
	if len(ss) == 2 {
		return len(ss[1])
	}

	return 0
}

// GetKey 由交易所和品种组成key
func GetKey(exchange string, symbol string) string {
	return exchange + "." + symbol
}

// GetGetExchangeLimit 获取限制值
func GetExchangeLimitCache(key string) *LimitValues {
	if value, ok := ExchangeLimitsMap.Load(key); ok {
		return value.(*LimitValues)
	}

	elv, err := FindExchangeLimit(bson.M{"es": key}, bson.M{})
	if err != nil {
		logger.Warn("model.FindExchangeLimit error", logger.Err(err), logger.String("params", key))
	} else {
		el := elv.ToLimitValues()
		SetExchangeLimitCache(key, el)
		return el
	}

	return &LimitValues{}
}

// GetExchangeLimitAllCache 获取所有限制值
func GetExchangeLimitAllCache() []*LimitValues {
	lvs := []*LimitValues{}
	ExchangeLimitsMap.Range(func(key, value interface{}) bool {
		lvs = append(lvs, value.(*LimitValues))
		return true
	})
	return lvs
}

// SetExchangeLimitCache 设置限制值
func SetExchangeLimitCache(key string, value *LimitValues) {
	ExchangeLimitsMap.Store(key, value)
}

// InitExchangeLimitCache 初始化交易所限制值
func InitExchangeLimitCache() error {
	if ExchangeLimitsMap == nil {
		ExchangeLimitsMap = new(sync.Map)
	}

	total, _ := CountExchangeLimits(bson.M{})
	limit := 100
	page := total / 100
	count := 0
	count2 := 0

	for i := 0; i <= page; i++ {
		els, err := FindExchangeLimits(bson.M{}, bson.M{}, i, limit)
		if err != nil {
			return err
		}

		for _, v := range els {
			if v.ES != "" {
				count++
				SetExchangeLimitCache(v.ES, v.ToLimitValues())

				if v.Currency == "" || v.AnchorCurrency == "" {
					continue
				}
				if _, ok := goex.CurrencyPairSyncMap.Load(v.Symbol); ok {
					continue
				}
				count2++
				SetCurrencyPairCache(v.Symbol, v.Currency, v.AnchorCurrency)
			}
		}
		time.Sleep(time.Millisecond * 50)
	}

	logger.Infof("InitExchangeLimitCache finish, success=%d, total=%d", count, total)
	logger.Infof("InitCurrencyPairCache finish, success=%d", count2)
	return nil
}

// -------------------------------------------------------------------------------------------------

// GetCurrencyPairCache 交易对
func GetCurrencyPairCache(key string) *goex.CurrencyPair {
	key = strings.ToLower(key)
	if value, ok := goex.CurrencyPairSyncMap.Load(key); ok {
		return value.(*goex.CurrencyPair)
	}

	elv, err := FindExchangeLimit(bson.M{"symbol": key}, bson.M{})
	if err != nil {
		logger.Warn("model.FindExchangeLimit error", logger.Err(err), logger.String("params", key))
		return &goex.CurrencyPair{}
	}

	SetCurrencyPairCache(key, elv.Currency, elv.AnchorCurrency)
	return &goex.CurrencyPair{
		CurrencyA: goex.Currency{Symbol: strings.ToUpper(elv.Currency)},
		CurrencyB: goex.Currency{Symbol: strings.ToUpper(elv.AnchorCurrency)},
	}
}

// SetCurrencyPairCache 设置交易对
func SetCurrencyPairCache(key string, currency string, anchorCurrency string) {
	key = strings.ToLower(key)
	goex.CurrencyPairSyncMap.Store(key, &goex.CurrencyPair{
		CurrencyA: goex.Currency{Symbol: strings.ToUpper(currency)},
		CurrencyB: goex.Currency{Symbol: strings.ToUpper(anchorCurrency)},
	})
}
