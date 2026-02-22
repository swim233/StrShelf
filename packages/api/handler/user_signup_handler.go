package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gopkg.ilharper.com/strshelf/api/db"
	"gopkg.ilharper.com/strshelf/api/lib"
	"gopkg.ilharper.com/strshelf/api/logger"
	"gorm.io/gorm"
)

func UserSignUpHandler(r *gin.Engine, DB db.ShelfDB) {
	r.POST("/v1/user.signup", func(ctx *gin.Context) {
		logger.Suger.Debugf("current allow signup status:%v", viper.GetBool("allow_signup"))
		if !viper.GetBool("allow_signup") {
			ctx.Status(404)
			return
		}
		newUser := lib.UserInfo{}
		result := gorm.WithResult()
		err := ctx.BindJSON(&newUser)
		if err != nil {
			ctx.JSON(500, gin.H{"msg": "internal error"})
			return
		}
		if newUser.Username == "" || newUser.Password == "" {
			ctx.JSON(200, gin.H{"code": "400", "msg": "username and password not allow empty"})
			return
		}
		users, err := DB.GetMatchUser(newUser.Username)
		if err != nil {
			ctx.JSON(500, gin.H{"msg": "internal error"})
			return
		}
		if len(users) != 0 {
			ctx.JSON(200, gin.H{"code": "400", "msg": "user already exist"})
			return
		}
		HashedPassword, err := HashPassword(newUser.Password)
		if err != nil {
			ctx.JSON(500, gin.H{"msg": "internal error"})
			return
		}

		err = DB.PushNewUser(newUser.Username, string(HashedPassword))
		if err != nil {
			ctx.JSON(500, gin.H{"msg": "internal error"})
			logger.Suger.Errorf("error in insert database: %s", err.Error())
			return
		}
		logger.Suger.Debugf("dev: %v", newUser)
		ctx.JSON(200, lib.PostRequestResponse{
			Code:   200,
			Result: result,
		})

	})

}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
