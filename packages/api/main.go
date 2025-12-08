package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ShelfItem struct {
	Id          uint64     `json:"id"`
	Title       string     `json:"title"`
	Link        string     `json:"link"`
	Comment     string     `json:"comment"`
	GMTCreated  CustomTime `json:"gmt_created"`
	GMTModified CustomTime `json:"gmt_modified"`
	GMTDeleted  CustomTime `json:"gmt_deleted"`
	Deleted     bool       `json:"deleted"`
}

type StrShelfResponse[T any] struct {
	Code uint16 `json:"code"`
	Data T      `json:"data"`
	Msg  string `json:"msg"`
}

type CustomTime time.Time

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	t := time.Time(ct)
	return json.Marshal(t.UnixMilli())
}

func main() {
	dsn := "host=localhost user=postgres dbname=strshelf port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		println(err.Error())
	}

	r := gin.Default()

	r.Use(cors.Default())

	r.POST("/v1/item.get", func(ctx *gin.Context) {
		shelfItems, err := gorm.G[ShelfItem](db).Raw("SELECT * FROM public.shelf_item_v1 WHERE deleted IS NOT TRUE ORDER BY id DESC").Find(context.Background())
		if err != nil {
			ctx.JSON(500, "error in getting items")
			return
		}

		ctx.JSON(200, StrShelfResponse[[]ShelfItem]{
			Code: 200,
			Data: shelfItems,
			Msg:  "",
		})
	})

	r.Run(":1111")

}
