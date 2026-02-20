package handler

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.ilharper.com/strshelf/api/db"
	"gopkg.ilharper.com/strshelf/api/lib"
	"gopkg.ilharper.com/strshelf/api/logger"
	"gopkg.ilharper.com/strshelf/api/token"
)

func UserLoginHandler(r *gin.Engine, DB db.ShelfDB) {
	r.POST("/v1/user.login", func(ctx *gin.Context) {
		user := lib.UserInfo{}
		err := ctx.BindJSON(&user)
		if err != nil {
			ctx.JSON(500, gin.H{"msg": "internal error"})
			return
		}

		matchUsers, err := DB.GetMatchUser(user.Username)
		if err != nil {
			ctx.JSON(500, gin.H{"msg": "error in insert database: " + err.Error()})
			logger.Suger.Errorf("error in matching login user: %s", err.Error())
			return
		}
		if len(matchUsers) == 1 {
			matchUser := matchUsers[0]
			err := bcrypt.CompareHashAndPassword([]byte(matchUser.Password), []byte(user.Password))
			if err != nil {
				ctx.JSON(401, gin.H{"msg": "password is incorrect"})
				logger.Suger.Warnf("a failure logging request: %s", err.Error())
				return
			} else {
				token := token.CreateJWT(user.Username)
				ctx.JSON(200, gin.H{"token": token})
				logger.Suger.Infof("user login successful,username: %s", matchUser.Username)
			}
		} else {
			logger.Suger.Warnf("not exist user login: %s", user.Username)
			ctx.JSON(401, gin.H{"msg": "user is not exist"})
			return
		}

	})
}
