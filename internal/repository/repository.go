package repository

import (
	"hitshop/internal/config"
	"hitshop/internal/core"

	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

// Authorization имплементирует логику хранения пользователей в базе
type Authorization interface {
	TestDB()
	CreateUser(user core.User) (int, error)
	GetUser(username, password string) (*core.User, error)
}

// Repository главная структура слоя репозиторий
type Repository struct {
	Authorization
}

// NewRepository конструктор
func NewRepository(db *sqlx.DB, cfg *config.Cfg) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
