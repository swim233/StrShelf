package db

import (
	"github.com/spf13/viper"
	"gopkg.ilharper.com/strshelf/api/lib"
	"gorm.io/gorm"
)

type DBInstance interface {
	GetShelfItem() ([]lib.ShelfItem, error)
	EditShelfItem(title string, link string, comment string, id uint64) error
	PostShelfItem(item lib.ShelfItem) error
	GetMatchUser(username string) ([]lib.UserInfo, error)
	PushNewUser(username, password string) error
	GetShelfItemByID(id uint64) ([]lib.ShelfItem, error)
	DeleteShelfItem(id uint64) error
}

var DB *gorm.DB

var ActivatedDB string

func InitDB() DBInstance {
	dbType := viper.GetString("db_type")
	dbType = "postgres"
	switch dbType {
	case "postgres":
		{
			pg := InitPostgresDB()
			return pg
		}
	case "sqlite":
		{
			InitSqliteDB()
		}

	}
	return nil
}
