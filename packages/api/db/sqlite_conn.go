package db

import (
	"context"
	"os"

	"github.com/glebarez/sqlite"
	"github.com/spf13/viper"
	"gopkg.ilharper.com/strshelf/api/lib"
	"gopkg.ilharper.com/strshelf/api/logger"
	"gorm.io/gorm"
)

type SqliteDB struct {
	DB *gorm.DB
}

func InitSqliteDB() *SqliteDB {
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
			return nil
		}
		return nil
	}
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Suger.Errorf("error when connect sqlite db: %s", err.Error())
		return nil
	}
	return &SqliteDB{
		DB: db,
	}
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

func (db *SqliteDB) GetShelfItem() ([]lib.ShelfItem, error) {
	// 查询未删除的商品，deleted = 0 表示未删除
	return gorm.G[lib.ShelfItem](db.DB).Raw("SELECT * FROM shelf_item_v1 WHERE deleted != 1 ORDER BY id DESC").Find(context.Background())
}

func (db *SqliteDB) EditShelfItem(title string, link string, comment string, id uint64) error {
	// 更新时用 CURRENT_TIMESTAMP 替代 now()
	return gorm.G[any](db.DB).Exec(context.Background(),
		"UPDATE shelf_item_v1 SET title = ?, link = ?, comment = ?, gmt_modified = CURRENT_TIMESTAMP WHERE id = ?",
		title, link, comment, id)
}

func (db *SqliteDB) PostShelfItem(item lib.ShelfItem) error {
	// 插入时 deleted 传布尔或整数0
	return gorm.G[any](db.DB).Exec(context.Background(),
		"INSERT INTO shelf_item_v1 (title, link, comment, deleted) VALUES (?, ?, ?, ?)",
		item.Title, item.Link, item.Comment, 0) // 0 表示未删除
}

func (db *SqliteDB) GetMatchUser(username string) ([]lib.UserInfo, error) {
	return gorm.G[lib.UserInfo](db.DB).Raw("SELECT * FROM shelf_user_v1 WHERE username = ?", username).Find(context.Background())
}

func (db *SqliteDB) PushNewUser(username, password string) error {
	return gorm.G[any](db.DB).Exec(context.Background(),
		"INSERT INTO shelf_user_v1 (username, password) VALUES (?, ?)",
		username, password)
}

func (db *SqliteDB) GetShelfItemByID(id uint64) ([]lib.ShelfItem, error) {
	return gorm.G[lib.ShelfItem](db.DB).Raw("SELECT * FROM shelf_item_v1 WHERE deleted != 1 AND id = ?", id).Find(context.Background())
}

func (db *SqliteDB) DeleteShelfItem(id uint64) error {
	// 逻辑删除，deleted 设为1，gmt_deleted 用 CURRENT_TIMESTAMP
	return gorm.G[any](db.DB).Exec(context.Background(),
		"UPDATE shelf_item_v1 SET deleted = 1, gmt_deleted = CURRENT_TIMESTAMP WHERE id = ?",
		id)
}
