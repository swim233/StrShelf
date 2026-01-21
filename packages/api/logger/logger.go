package logger

import (
	"log"

	"go.uber.org/zap"
)

var Logger *zap.Logger
var Suger *zap.SugaredLogger

func InitLogger() {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("Fail to init logger")
	}
	Suger = Logger.Sugar()
}
