package handler

import (
	"github.com/gin-gonic/gin"
	"gopkg.ilharper.com/strshelf/api/db"
	"gopkg.ilharper.com/strshelf/api/lib"
)

func ItemGetHandler(r *gin.Engine, DB db.DBInstance) {
	r.POST("/v1/item.get", func(ctx *gin.Context) {
		shelfItems, err := DB.GetShelfItem()
		if err != nil {
			ctx.JSON(500, gin.H{"msg": "internal error"})
			return
		}

		ctx.JSON(200, lib.StrShelfResponse[[]lib.ShelfItem]{
			Code: 200,
			Data: shelfItems,
			Msg:  "",
		})
	})

}
