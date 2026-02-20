package handler

import (
	"github.com/gin-gonic/gin"
	"gopkg.ilharper.com/strshelf/api/db"
	"gopkg.ilharper.com/strshelf/api/lib"
	"gopkg.ilharper.com/strshelf/api/logger"
	"gopkg.ilharper.com/strshelf/api/middleware"
	"gorm.io/gorm"
)

func ItemPostHandler(r *gin.Engine, DB db.ShelfDB) {

	r.POST("/v1/item.post", middleware.JWTAuthMiddleWare(), func(ctx *gin.Context) {
		newShelfItem := lib.ShelfItem{}
		result := gorm.WithResult()
		err := ctx.ShouldBindJSON(&newShelfItem)
		if err != nil {
			ctx.JSON(500, gin.H{"msg": "error in binding json: " + err.Error()})
			logger.Suger.Errorf("error in getting items: %s", err.Error())
			return
		}

		err = DB.PostShelfItem(newShelfItem)

		if err != nil {
			ctx.JSON(500, gin.H{"msg": "internal error"})
			logger.Suger.Errorf("error in insert database: %s", err.Error())
			return
		}
		ctx.JSON(200, lib.PostRequestResponse{
			Code:   200,
			Result: result,
		})
	})
}
