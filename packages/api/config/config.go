package config

import (
	"github.com/spf13/viper"
	"gopkg.ilharper.com/strshelf/api/logger"
)

var ViperInstance *viper.Viper

func InitConfig() {
	viper.AddConfigPath("packages/api/config/")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("./")
	viper.AddConfigPath(".")
	viper.SetConfigName("config.json")
	viper.SetConfigType("json")
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		logger.Suger.Errorf("can not load config: %s", err.Error())
	} else {
		logger.Suger.Infof("config load successful: %s", viper.ConfigFileUsed())
		logger.Suger.Debugln(viper.AllKeys())
	}
}
