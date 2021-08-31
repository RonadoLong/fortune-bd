package model

import (
	"sync"

	"github.com/globalsign/mgo/bson"
	"github.com/zhufuyi/pkg/logger"
	"github.com/zhufuyi/pkg/mongo"
)

// GridTypeCollectionName 表名
const GridTypeCollectionName = "strategyType"

// StrategyType 策略类型数据
type StrategyType struct {
	mongo.PublicFields `bson:",inline"`

	Type     int      `json:"type" bson:"type"`         // 0:网格交易，1:杠杆网格，2:借贷网格，3:反向网格，4:反向杠杆网格，5:无线网格
	Name     string   `json:"name" bson:"name"`         // 网格名称
	Labels   []string `json:"labels" bson:"labels"`     // 标签
	Describe string   `json:"describe" bson:"describe"` // 说明
}

// Insert 插入一条新的记录
func (object *StrategyType) Insert() (err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()
	object.SetFieldsValue()

	return mconn.Insert(GridTypeCollectionName, object)
}

// FindStrategyType 获取单条记录
func FindStrategyType(selector bson.M, field bson.M) (*StrategyType, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	object := &StrategyType{}
	return object, mconn.FindOne(GridTypeCollectionName, object, selector, field)
}

// FindStrategyTypes 获取多条记录
func FindStrategyTypes(selector bson.M, field bson.M, page int, limit int, sort ...string) ([]*StrategyType, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	// 默认从第一页开始
	if page < 0 {
		page = 0
	}
	if len(sort) == 0 || sort[0] == "" {
		sort = []string{"-_id"}
	}

	objects := []*StrategyType{}
	return objects, mconn.FindAll(GridTypeCollectionName, &objects, selector, field, page, limit, sort...)
}

// UpdateStrategyType 更新单条记录
func UpdateStrategyType(selector, update bson.M) (err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateOne(GridTypeCollectionName, selector, mongo.UpdatedTime(update))
}

// UpdateStrategyTypes 更新多条记录
func UpdateStrategyTypes(selector, update bson.M) (n int, err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateAll(GridTypeCollectionName, selector, mongo.UpdatedTime(update))
}

// FindAndModifyStrategyType 更新并返回最新记录
func FindAndModifyStrategyType(selector bson.M, update bson.M) (*StrategyType, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	object := &StrategyType{}
	return object, mconn.FindAndModify(GridTypeCollectionName, object, selector, mongo.UpdatedTime(update))
}

// CountStrategyTypes 统计数量，不包括删除记录
func CountStrategyTypes(selector bson.M) (n int, err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.Count(GridTypeCollectionName, mongo.ExcludeDeleted(selector))
}

// DelStrategyType 删除记录
func DeleteStrategyType(selector bson.M) (int, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateAll(GridTypeCollectionName, selector, mongo.DeletedTime(bson.M{}))
}

// -------------------------------------------------------------------------------------------------

// 策略类型数据缓存，type值作为key，对象StrategyType作为值
var strategyTypeCache = new(sync.Map)

// GetStrategyTypeCache 读取
func GetStrategyTypeCache(key int) *StrategyType {
	if value, ok := strategyTypeCache.Load(key); ok {
		return value.(*StrategyType)
	}

	st, err := FindStrategyType(bson.M{"type": key}, bson.M{})
	if err != nil {
		logger.Warn("FindStrategyType error", logger.Err(err), logger.Int("type", key))
	} else {
		SetStrategyTypeCache(key, st)
		return st
	}

	return &StrategyType{}
}

// SetStrategyTypeCache 保存
func SetStrategyTypeCache(key int, value *StrategyType) {
	strategyTypeCache.Store(key, value)
}

// DelStrategyTypeCache 删除
func DelStrategyTypeCache(key int) {
	strategyTypeCache.Delete(key)
}

// InitStrategyTypeCache 初始化策略类型数据
func InitStrategyTypeCache() error {
	if strategyTypeCache == nil {
		strategyTypeCache = new(sync.Map)
	}

	sts, err := FindStrategyTypes(bson.M{}, bson.M{}, 0, 100)
	if err != nil {
		return err
	}
	count := 0
	for _, st := range sts {
		if st.Name != "" {
			count++
			SetStrategyTypeCache(st.Type, st)
		}
	}

	logger.Infof("InitStrategyTypeCache finish, success=%d, total=%d", count, len(sts))

	return nil
}
