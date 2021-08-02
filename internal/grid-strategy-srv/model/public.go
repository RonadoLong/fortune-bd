package model

import (
	"errors"
	"time"

	"github.com/globalsign/mgo/bson"
)

// PublicFields 公共字段
type PublicFields struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`                                  // 唯一ID
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`                     // 创建时间
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`                     // 修改时间
	DeletedAt *time.Time    `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"` // 删除时间
}

// SetFieldsValue 设置公共字段值，在插入数据时使用
func (p *PublicFields) SetFieldsValue() {
	now := time.Now()
	if !p.ID.Valid() {
		p.ID = bson.NewObjectId()
	}

	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
}

// PublicFieldsInt 设置公共字段值，自定id，在插入数据时使用
type PublicFieldsInt struct {
	ID        int64      `bson:"_id" json:"id,string"`                           // 唯一ID
	CreatedAt time.Time  `bson:"createdAt" json:"createdAt"`                     // 创建时间
	UpdatedAt time.Time  `bson:"updatedAt" json:"updatedAt"`                     // 修改时间
	DeletedAt *time.Time `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"` // 删除时间
}

// SetFieldsValue 初始化
func (m *PublicFieldsInt) SetFieldsValue(newID int64) {
	t := time.Now()

	if m.ID == 0 {
		m.ID = newID
	}

	if m.CreatedAt.IsZero() {
		m.CreatedAt = t
	}
}

// ----------------------------------- 公共函数 ----------------------------------------

// CheckUpdateContent 执行更新操作之前先判断有没有$操作符
func CheckUpdateContent(update bson.M) error {
	for k := range update {
		if k[0] != '$' {
			return errors.New("update content must start with '$'")
		}
	}
	return nil
}

// ExcludeDeleted 不包含已删除的
func ExcludeDeleted(selector bson.M) bson.M {
	selector["deletedAt"] = bson.M{"$exists": false}
	return selector
}

// UpdatedTime 更新updatedAt时间
func UpdatedTime(update bson.M) bson.M {
	if v, ok := update["$set"]; ok {
		v.(bson.M)["updatedAt"] = time.Now()
	} else {
		update["$set"] = bson.M{"updatedAt": time.Now()}
	}
	return update
}

// DeletedTime 更新deletedAt时间
func DeletedTime(update bson.M) bson.M {
	if v, ok := update["$set"]; ok {
		v.(bson.M)["deletedAt"] = time.Now()
	} else {
		update["$set"] = bson.M{"deletedAt": time.Now()}
	}
	return update
}
