package repository

import (
	"hitshop/internal/config"
	"hitshop/internal/core"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

// Authorization имплементирует логику хранения пользователей в базе
type Authorization interface {
	TestDB()
	CreateAccount(acc *core.Account) (uuid.UUID, error)
	UpdateAccount(acc *core.Account) (uuid.UUID, error)
	GetAccountByEmail(email, password string) (*core.Account, error)
	GetAccountByCode(verification_code string) (*core.Account, error)
	GetAccountByUUID(UUID uuid.UUID) (*core.Account, error)
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
