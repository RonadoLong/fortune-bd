// Pipe - A small and beautiful blogging platform written in golang.
// Copyright (C) 2017-2018, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package dbutil

import (
	"os"
	"shop-web/common/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Logger
var logger = log.NewLogger(os.Stdout)

var DB *gorm.DB

func ConnectDB(mysql string, logLevel string) {
	var err error
	DB, err = gorm.Open("mysql", mysql)

	//if err = DB.AutoMigrate(commonUtils.Models...).Error; nil != err {
	//	logger.Fatal("auto migrate tables failed: " + err.Error())
	//}
	if nil != err {
		logger.Fatalf("opens database failed: " + err.Error())
	}
	if logLevel == "dev" {
		DB.LogMode(true)
	}else {
		DB.LogMode(false)
	}
	DB.SingularTable(true)
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
}

func DisconnectDB() {
	if err := DB.Close(); nil != err {
		logger.Errorf("Disconnect from database failed: " + err.Error())
	}
}
