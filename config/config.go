package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func InitConfig() {
	work, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	//fmt.Println(work)
	viper.AddConfigPath(work + "/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.WatchConfig()
	//viper.ReadRemoteConfig()
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

}
