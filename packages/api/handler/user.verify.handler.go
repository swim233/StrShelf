package handler

import (
	"github.com/gin-gonic/gin"
	"gopkg.ilharper.com/strshelf/api/logger"
	"gopkg.ilharper.com/strshelf/api/token"
)

func UserVerifyHandler(r *gin.Engine) {
	r.POST("/v1/user.verify", func(ctx *gin.Context) {
		var tokenReq token.TokenRequest
		err := ctx.Bind(&tokenReq)
		if err != nil {
			logger.Suger.Errorf("error in verifying user token: %s", err.Error())
			ctx.JSON(500, gin.H{"msg": "internal error"})
			return
		}
		if result, err := token.VerifyJWT(string(tokenReq.Token)); err != nil {
			logger.Suger.Warnf("error in verifying user token: %s", err.Error())
			logger.Suger.Warnf("received token: %s", tokenReq.Token)
			ctx.JSON(401, gin.H{
				"msg": "token is invalid"})
			return
		} else {

			switch result {
			case true:
				ctx.JSON(200, gin.H{"msg": "successful"})
			case false:
				logger.Suger.Warnf("user post a error token: %s", tokenReq.Token)
				ctx.JSON(401, gin.H{"msg": "token is invalid"})
			}
			return
		}
	})
}
