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
	DBSslMode      string `mapstructure:"POSTGRES_SSL_MODE"`

	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	ServerPort    string `mapstructure:"SERVER_PORT"`

	LogJsonMod bool   `mapstructure:"LOG_JSON_MOD"`
	LogLevel   string `mapstructure:"LOG_LEVEL"`

	EmailFrom    string `mapstructure:"EMAIL_FROM"`
	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPUser     string `mapstructure:"SMTP_USER"`
	SMTPPass     string `mapstructure:"SMTP_PASS"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`
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
