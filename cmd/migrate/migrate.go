package main

import (
	"fmt"
	"hitshop/internal/config"
	"hitshop/models"

	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func init() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Europe/Moscow",
		cfg.DBHost,
		cfg.DBUserName,
		cfg.DBUserPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSslMode,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}
	fmt.Println("🚀 Connected Successfully to the Database")
}

func main() {
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	ar := models.AccountRoles{}
	DB.AutoMigrate(&ar)
	a := models.AccountStatuses{}
	DB.AutoMigrate(&a)
	a1 := models.AccountStat{}
	DB.AutoMigrate(&a1)
	a2 := models.AccountFavorites{}
	DB.AutoMigrate(&a2)
	a3 := models.OrdersStatuses{}
	DB.AutoMigrate(&a3)
	a4 := models.Stores{}
	DB.AutoMigrate(&a4)
	a5 := models.Orders{}
	DB.AutoMigrate(&a5)
	a7 := models.Category{}
	DB.AutoMigrate(&a7)
	a8 := models.Manufactures{}
	DB.AutoMigrate(&a8)
	a9 := models.Products{}
	DB.AutoMigrate(&a9)
	a10 := models.PriceChanges{}
	DB.AutoMigrate(&a10)
	a11 := models.Deliveries{}
	DB.AutoMigrate(&a11)
	a6 := models.Accounts{}
	DB.AutoMigrate(&a6)

	// err := DB.Migrator().DropColumn(&c, "password")
	// if err != nil {
	// 	// Do whatever you want to do!
	// 	log.Print("ERROR: We expect the description column to be drop-able")
	// }
	fmt.Println("👍 Migration complete")
}
