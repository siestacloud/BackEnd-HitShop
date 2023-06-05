package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/viper"
)

type Cfg struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`
	ServerPort     string `mapstructure:"SERVER_PORT"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`
	LogJsonMod   bool   `mapstructure:"LOG_JSON_MOD"` // log format in json
	LogLevel     string `mapstructure:"LOG_LEVEL"`    // info,debug
}

func LoadConfig(path string) (config Cfg, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	err = env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}
	return
}
