package main

import (
	"fmt"
	"hitshop/internal/config"
	"hitshop/models"

	"hitshop/internal/repository"

	"log"
)

func init() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}
	repository.ConnectDB(&cfg)
}

func main() {
	repository.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	ar := models.AccountRoles{}
	repository.DB.AutoMigrate(&ar)
	a := models.AccountStatuses{}
	repository.DB.AutoMigrate(&a)
	a1 := models.AccountStat{}
	repository.DB.AutoMigrate(&a1)
	a2 := models.AccountFavorites{}
	repository.DB.AutoMigrate(&a2)
	a3 := models.OrdersStatuses{}
	repository.DB.AutoMigrate(&a3)
	a4 := models.Stores{}
	repository.DB.AutoMigrate(&a4)
	a5 := models.Orders{}
	repository.DB.AutoMigrate(&a5)
	a7 := models.Category{}
	repository.DB.AutoMigrate(&a7)
	a8 := models.Manufactures{}
	repository.DB.AutoMigrate(&a8)
	a9 := models.Products{}
	repository.DB.AutoMigrate(&a9)
	a10 := models.PriceChanges{}
	repository.DB.AutoMigrate(&a10)
	a11 := models.Deliveries{}
	repository.DB.AutoMigrate(&a11)
	a6 := models.Accounts{}
	repository.DB.AutoMigrate(&a6)

	// err := repository.DB.Migrator().DropColumn(&c, "password")
	// if err != nil {
	// 	// Do whatever you want to do!
	// 	log.Print("ERROR: We expect the description column to be drop-able")
	// }
	fmt.Println("👍 Migration complete")
}
