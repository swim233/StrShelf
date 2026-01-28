package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"gopkg.ilharper.com/strshelf/api/logger"
)

func CreateJWT(username string) string {
	claims := &jwt.RegisteredClaims{
		Issuer: username,
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(time.Hour * 720),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(viper.GetString("secret_key")))
	if err != nil {
		logger.Suger.Errorf("error when creating jwt token: %s", err.Error())
		return ""
	}
	return ss
}

func VerifyJWT(tokenString string) (success bool, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("excepted HS256: %v", t.Header["alg"])
		}
		return []byte(viper.GetString("secret_key")), nil

	})
	if err != nil {
		return false, err
	}
	return token.Valid, nil
}
