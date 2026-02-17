package db

import (
	"github.com/spf13/viper"
	"gopkg.ilharper.com/strshelf/api/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TODO:增加更多数据库驱动
func InitPostgresDB() {
	dsn := func() string {
		if dsn := viper.GetString("dsn"); dsn != "" {
			return dsn
		} else {
			return "host=localhost user=postgres dbname=strshelf port=5432 sslmode=disable TimeZone=Asia/Shanghai"

		}
	}()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Suger.Errorf("error when connect db: %s", err.Error())
		return
	}
	DB = db
}
