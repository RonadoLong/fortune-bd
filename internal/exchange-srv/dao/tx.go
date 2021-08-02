package dao

import "github.com/jinzhu/gorm"

func (d *Dao) BeginTran() *gorm.DB {
	return d.db.Begin()
}
