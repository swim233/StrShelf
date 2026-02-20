package db

import (
	"github.com/spf13/viper"
	"gopkg.ilharper.com/strshelf/api/lib"
	"gopkg.ilharper.com/strshelf/api/logger"
	"gorm.io/gorm"
)

type ShelfDB interface {
	GetShelfItem() ([]lib.ShelfItem, error)
	EditShelfItem(title string, link string, comment string, id uint64) error
	PostShelfItem(item lib.ShelfItem) error
	GetMatchUser(username string) ([]lib.UserInfo, error)
	PushNewUser(username, password string) error
	GetShelfItemByID(id uint64) ([]lib.ShelfItem, error)
	DeleteShelfItem(id uint64) error
}

var DB *gorm.DB

func InitDB() ShelfDB {
	dbType := viper.GetString("db_type")
	switch dbType {
	case "postgres":
		{
			pg := InitPostgresDB()
			return pg
		}
	case "sqlite":
		{
			sqlite := InitSqliteDB()
			return sqlite

		}
	default:
		{
			logger.Suger.Panicf("fail to get database type,please check config,current type: %s", dbType)
		}
	}
	return nil
}
