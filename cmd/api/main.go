package main

import (
	"log"

	"hitshop/internal/config"
	"hitshop/internal/repository"
	"hitshop/internal/service"
	"hitshop/internal/transport/rest"
	"hitshop/internal/transport/rest/handler"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
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

	// * парсинг env
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	// *подключение базы
	db, err := repository.ConnectDB(&cfg)
	if err != nil {
		return
	}

	repos := repository.NewRepository(db, &cfg)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services, &cfg)

	srv, err := rest.NewServer(&cfg, handlers)
	if err != nil {
		log.Fatal()
	}
	if err := srv.Run(); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
	// if err := db.Close(); err != nil {
	// 	logrus.Errorf("error occured on db connection close: %s", err.Error())
	// }
}
