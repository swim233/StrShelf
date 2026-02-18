package db

import (
	"context"

	"github.com/spf13/viper"
	"gopkg.ilharper.com/strshelf/api/lib"
	"gopkg.ilharper.com/strshelf/api/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	DB *gorm.DB
}

func InitPostgresDB() *PostgresDB {
	dsn := func() string {
		if dsn := viper.GetString("dsn"); dsn != "" {
			return dsn
		} else {
			return "host=localhost user=postgres dbname=strshelf port=5432 sslmode=disable TimeZone=Asia/Shanghai"

		}
	}()
	logger.Suger.Debugf("postgres dsn: %s", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Suger.Errorf("error when connect postgres db: %s", err.Error())
		return &PostgresDB{}
	}
	return &PostgresDB{
		DB: db,
	}

}

func (db *PostgresDB) GetShelfItem() ([]lib.ShelfItem, error) {
	return gorm.G[lib.ShelfItem](db.DB).Raw("SELECT * FROM public.shelf_item_v1 WHERE deleted IS NOT TRUE ORDER BY id DESC").Find(context.Background())
}

func (db *PostgresDB) EditShelfItem(title string, link string, comment string, id uint64) error {
	return gorm.G[any](db.DB).Exec(context.Background(), "UPDATE public.shelf_item_v1 SET title = ? ,link = ? , comment = ? ,gmt_modified = now() WHERE id = ?", title, link, comment, id)
}

func (db *PostgresDB) PostShelfItem(item lib.ShelfItem) error {
	return gorm.G[any](db.DB).Exec(context.Background(), "INSERT INTO public.shelf_item_v1 (title,link,comment,deleted) VALUES(?,?,?,?)", item.Title, item.Link, item.Comment, false)
}

func (db *PostgresDB) GetMatchUser(username string) ([]lib.UserInfo, error) {
	return gorm.G[lib.UserInfo](db.DB).Raw("SELECT * FROM public.shelf_user_v1 WHERE username = ?", username).Find(context.Background())
}

func (db *PostgresDB) PushNewUser(username, password string) error {
	return gorm.G[any](db.DB).Exec(context.Background(), "INSERT INTO public.shelf_user_v1 (username,password) VALUES(?,?)", username, password)
}

func (db *PostgresDB) GetShelfItemByID(id uint64) ([]lib.ShelfItem, error) {
	return gorm.G[lib.ShelfItem](db.DB).Raw("SELECT * FROM public.shelf_item_v1 WHERE deleted IS NOT TRUE AND id = ?", id).Find(context.Background())

}
func (db *PostgresDB) DeleteShelfItem(id uint64) error {
	return gorm.G[any](db.DB).Exec(context.Background(), "UPDATE public.shelf_item_v1 SET deleted = true ,gmt_deleted = now() WHERE id = ?", id)
}
