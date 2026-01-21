package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port int `json:"port"`
}

var ViperInstance *viper.Viper
var ConfigInstance Config

func InitConfig() {
	ViperInstance = viper.New()
	v := ViperInstance
	v.AddConfigPath("./")
	v.AddConfigPath(".")
	v.SetConfigName("config.json")

	err := v.ReadInConfig()
	if err != nil {
		println(err.Error())
	} else {
		ConfigInstance.Port = v.GetInt("port")
	}
}
