package repository

import (
	"errors"
	"fmt"
	"hitshop/internal/config"
	"hitshop/pkg"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	accountsTable         = "accounts" // пользователи данного сервиса
	accountFavoritesTable = "account_favorites"
	accountRolesTable     = "account_roles"
	accountStatsTable     = "account_stats"
	accountStatusesTable  = "account_statuses"
	categoriesTable       = "categories"
	deliveriesTable       = "deliveries"
	manufacturesTable     = "manufactures"
	ordersTable           = "orders"
	ordersStatusesTable   = "orders_statuses"
	priceChangesTable     = "price_changes"
	productsTable         = "products"
	storesTable           = "stores"
)

func ConnectDB(config *config.Cfg) (*sqlx.DB, error) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Europe/Moscow",
		config.DBHost,
		config.DBUserName,
		config.DBUserPassword,
		config.DBName,
		config.DBPort,
		config.DBSslMode,
	)
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	tables := []string{
		accountsTable,
		accountFavoritesTable,
		accountRolesTable,
		accountStatsTable,
		accountStatusesTable,
		categoriesTable,
		deliveriesTable,
		manufacturesTable,
		ordersTable,
		ordersStatusesTable,
		priceChangesTable,
		productsTable,
		storesTable,
	}

	err = checkTablesExist(db, tables)
	if err != nil {
		return nil, err
	}

	fmt.Println("🚀 Connected Successfully to the Database")
	return db, nil

}

func checkTablesExist(db *sqlx.DB, tables []string) error {

	var checkExist bool
	for _, t := range tables {

		row := db.QueryRow(fmt.Sprintf("SELECT EXISTS (SELECT FROM pg_tables WHERE  tablename  = '%s');", t))
		if err := row.Scan(&checkExist); err != nil {
			return err
		}
		if !checkExist {
			pkg.ErrPrintR("error", "🚀 Table not exists in Database: "+t)
			return errors.New("Database error")
		}
	}

	return nil
}
