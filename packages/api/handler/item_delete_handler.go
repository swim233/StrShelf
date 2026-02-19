package handler

import (
	"github.com/gin-gonic/gin"
	"gopkg.ilharper.com/strshelf/api/db"
	"gopkg.ilharper.com/strshelf/api/lib"
	"gopkg.ilharper.com/strshelf/api/logger"
	"gopkg.ilharper.com/strshelf/api/middleware"
)

func ItemDeleteHandler(r *gin.Engine, DB db.DBInstance) {
	r.POST("/v1/item.delete", middleware.JWTAuthMiddleWare(), func(ctx *gin.Context) {
		deleteItem := lib.ShelfDeleteItem{}
		err := ctx.BindJSON(&deleteItem)
		if err != nil {
			logger.Suger.Errorf("error when binding json: %s", err.Error())
			ctx.JSON(400, gin.H{"msg": "internal error"})
			return
		}
		shelfItems, err := DB.GetShelfItemByID(deleteItem.Id)

		if err != nil {
			logger.Suger.Errorf("error in checking item: %s", err.Error())
			ctx.JSON(400, gin.H{"msg": "internal error"})
			return
		}
		if len(shelfItems) != 1 {
			ctx.JSON(200, gin.H{"msg": "origin item not found", "code": "404"})
			return
		}
		err = DB.DeleteShelfItem(deleteItem.Id)

		if err != nil {
			logger.Suger.Errorf("error in updating database: %s", err.Error())
			ctx.JSON(500, gin.H{"msg": "internal error"})
			return
		}
		ctx.JSON(200, gin.H{"msg": "ok"})

	})
}
