package config

import (
	"github.com/spf13/viper"
	"gopkg.ilharper.com/strshelf/api/logger"
)

type Config struct {
	Port int `json:"port"`
}

var ViperInstance *viper.Viper
var ConfigInstance Config

func InitConfig() {
	ViperInstance = viper.New()
	v := ViperInstance
	v.AddConfigPath("packages/api/config/")
	v.AddConfigPath("./config/")
	v.AddConfigPath("./")
	v.AddConfigPath(".")
	v.SetConfigName("config.json")
	v.SetConfigType("json")
	err := v.ReadInConfig()
	if err != nil {
		logger.Suger.Errorf("can not load config: %s", err.Error())
	} else {
		logger.Suger.Infof("config load successful")
		ConfigInstance.Port = v.GetInt("port")
	}
}
