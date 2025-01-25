package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	if !valid() {
		panic("config error")
	}
}

func valid() bool {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return viper.GetString("mongodb.host") != "" && viper.GetString("mongodb.port") != ""
}
