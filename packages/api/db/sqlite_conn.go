package db

import (
	"os"

	"github.com/glebarez/sqlite"
	"github.com/spf13/viper"
	"gopkg.ilharper.com/strshelf/api/logger"
	"gorm.io/gorm"
)

type SqliteDB struct {
}

func InitSqliteDB() {
	dsn := func() string {
		if dsn := viper.GetString("sqlite_dsn"); dsn != "" {
			return dsn
		} else {
			logger.Suger.Warnf("fail to open sqlite file: %s, trying to make db/data.db", dsn)
			return "db/data.db"

		}
	}()
	if !checkSqliteDB(dsn) {
		err := os.MkdirAll("db", 0700)
		if err != nil {
			logger.Suger.Errorf("fail to make db dir: %s", err.Error())
		}
		_, err = os.Create("db/data.db")
		if err != nil {
			logger.Suger.Errorf("fail to create db file: %s", err.Error())
			return
		}
		return
	}
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Suger.Errorf("error when connect sqlite db: %s", err.Error())
		return
	}
	DB = db
}

func checkSqliteDB(dbPath string) bool {
	fInfo, err := os.Stat(dbPath)
	logger.Suger.Debugf("checking %s", dbPath)
	if err != nil {
		logger.Suger.Errorf("fail to stat db file: %s", err.Error())
		return false
	} else {
		logger.Suger.Debugln(fInfo)
		return true
	}

}

// func (db *SqliteDB) GetShelfItem() lib.ShelfItem {

// }
