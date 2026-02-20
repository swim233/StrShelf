package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.ilharper.com/strshelf/api/logger"
)

var ViperInstance *viper.Viper

func InitConfig() {
	viper.AddConfigPath("packages/api/config/")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("./")
	viper.AddConfigPath(".")
	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Suger.Errorf("can not load config: %s", err.Error())
	} else {
		logger.Suger.Infof("config load successful: %s", viper.ConfigFileUsed())
		logger.Suger.Debugln(viper.AllKeys())
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		logger.Suger.Infof("detaching config changed: %s", in.Name)
		logger.Suger.Debugf("config file event operation: %s", in.Op)

	})
}
