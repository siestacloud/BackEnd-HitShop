package repository

import (
	"fmt"
	"hitshop/internal/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB

	clientsTable                       = "clients"                                // пользователи данного сервиса
	tAccountsTable                     = "telegram_accounts"                      // аккаунты через которые ведется вся работа
	tAccountsAdditionalAttributesTable = "telegram_account_additional_attributes" // доп инфа об аккаутах
	tUntrustSessionsTable              = "telegram_untrust_sessions"              // недоверенные сессии
	tTrustSessionsTable                = "telegram_trust_sessions"                // доверенные сессии (создаются на основе недоверенных)
	tAppsTable                         = "telegram_apps"                          // данные по регистрации телеграм-клиентов
	tGroupsTable                       = "telegram_groups"                        // группы телеграм по которым работают аккаунты
	tUsersTable                        = "telegram_users"                         // целевые пользователи телеги
	tGroupsUsersTable                  = "telegram_groups_users"                  // связующая таблица (группы с пользователями, многие ко многим)
)

func ConnectDB(config *config.Cfg) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow",
		config.DBHost,
		config.DBUserName,
		config.DBUserPassword,
		config.DBName,
		config.DBPort,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}
	fmt.Println("🚀 Connected Successfully to the Database")
}
