package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.ilharper.com/strshelf/api/logger"
	"gopkg.ilharper.com/strshelf/api/token"
)

func JWTAuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(401, gin.H{"msg": "need login"})
			ctx.Abort()
			return
		}
		args := strings.Split(authHeader, " ")
		if len(args) != 2 {
			ctx.JSON(401, gin.H{"msg": "invalid number of arguments"})
			ctx.Abort()
			return
		}
		arg := args[1]
		result, err := token.VerifyJWT(arg)
		if err != nil {
			logger.Suger.Errorf("error in verifying token: %s", err.Error())
			ctx.JSON(401, gin.H{"msg": "internal error"})
			ctx.Abort()
			return
		}
		if !result {
			ctx.JSON(401, gin.H{"msg": "invalid token"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
