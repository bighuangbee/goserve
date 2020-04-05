package dbModel

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"goserve/modules/config"
	"goserve/modules/loger"

	_ "database/sql/driver"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func Setup(){

	var err error
	DB, err = gorm.Open(config.ConfigData.DbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.ConfigData.DbUser,
		config.ConfigData.DbPassword,
		config.ConfigData.DbHost,
		config.ConfigData.DbName))

	if err != nil {
		panic("DataBase Connected Failed!" + err.Error())
	}

	//DB.LogMode(true)
	DB.SingularTable(true)
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)

	loger.Info("MySql Client SetUp Success...", config.ConfigData.DbHost)
}