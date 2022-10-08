package main

import (
	"log"
	"os"

	"gitlab.com/siteasservice/project-architecture/templates/template-svc-golang/internal/repository"
	"gitlab.com/siteasservice/project-architecture/templates/template-svc-golang/internal/service"
	"gitlab.com/siteasservice/project-architecture/templates/template-svc-golang/internal/transport/rest"
	"gitlab.com/siteasservice/project-architecture/templates/template-svc-golang/internal/transport/rest/handler"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title Template App API
// @version 1.0
// @description API Server for Template Application

// @host localhost:9999
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

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

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
