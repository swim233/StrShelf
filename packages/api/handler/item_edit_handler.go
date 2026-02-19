package handler

import (
	"github.com/gin-gonic/gin"
	"gopkg.ilharper.com/strshelf/api/db"
	"gopkg.ilharper.com/strshelf/api/lib"
	"gopkg.ilharper.com/strshelf/api/logger"
	"gopkg.ilharper.com/strshelf/api/middleware"
	"gorm.io/gorm"
)

func ItemEditHandler(r *gin.Engine, DB db.DBInstance) {
	r.POST("/v1/item.edit", middleware.JWTAuthMiddleWare(), func(ctx *gin.Context) {
		editShelfItem := lib.ShelfEditItem{}
		result := gorm.WithResult()
		err := ctx.ShouldBind(&editShelfItem)
		if err != nil {
			logger.Suger.Errorf("error in building editShelf: %s", err.Error())
			ctx.JSON(500, gin.H{"msg": "internal error"})
			return
		}

		shelfItems, err := DB.GetShelfItemByID(editShelfItem.Id)
		logger.Suger.Debugf("edit item: %v", shelfItems)
		if err != nil {
			logger.Suger.Errorf("error in checking item: %s", err.Error())
			ctx.JSON(400, gin.H{"msg": "internal error"})
			return
		}
		if len(shelfItems) != 1 {
			ctx.JSON(200, gin.H{"msg": "origin item not found", "code": "404"})
			return
		}
		err = DB.EditShelfItem(editShelfItem.NewTitle, editShelfItem.NewLink, editShelfItem.NewComment, editShelfItem.Id)

		if err != nil {
			logger.Suger.Errorf("error in updating database: %s", err.Error())
			ctx.JSON(500, gin.H{"msg": "internal error"})
			return
		}
		ctx.JSON(200, gin.H{"msg": "ok", "result": result})
	})
}
