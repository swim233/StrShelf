package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"gopkg.ilharper.com/strshelf/api/db"
	"gorm.io/gorm"
)

func ItemGetHandler(r *gin.Engine) {
	r.POST("/v1/item.get", func(ctx *gin.Context) {
		shelfItems, err := gorm.G[ShelfItem](db.DB).Raw("SELECT * FROM public.shelf_item_v1 WHERE deleted IS NOT TRUE ORDER BY id DESC").Find(context.Background())
		if err != nil {
			ctx.JSON(500, gin.H{"msg": "internal error"})
			return
		}

		ctx.JSON(200, StrShelfResponse[[]ShelfItem]{
			Code: 200,
			Data: shelfItems,
			Msg:  "",
		})
	})

}
