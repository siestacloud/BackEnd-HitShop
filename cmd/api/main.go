package main

import (
	"log"

	"tservice-checker/internal/config"
	"tservice-checker/internal/repository"
	"tservice-checker/internal/service"
	"tservice-checker/internal/transport/rest"
	"tservice-checker/internal/transport/rest/handler"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	cfg config.Cfg
)

// @title Telegram Checker Service
// @version 1.0
// @description REST API Server for parsing tdata folders, validate and save sessions

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	// парсинг flags, env
	err := config.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	// подключение базы
	db, err := repository.NewPostgresDB(cfg.URLPostgres)
	if err != nil {
		logrus.Warnf("failed to initialize postrges: %s", err.Error())
	}
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv, err := rest.NewServer(handlers)
	if err != nil {
		log.Fatal()
	}

	if err := srv.Run(viper.GetString("port")); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}

	// if err := db.Close(); err != nil {
	// 	logrus.Errorf("error occured on db connection close: %s", err.Error())
	// }
}
