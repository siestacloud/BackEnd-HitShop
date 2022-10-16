package repository

import (
	"tservice-checker/internal/core"

	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

// Authorization имплементирует логику хранения пользователей в базе
type Authorization interface {
	TestDB()
	CreateUser(user core.User) (int, error)
	GetUser(username, password string) (*core.User, error)
}

// Session имплементирует логику хранения телеграм сессий в базе
type Session interface {
	SaveSession(session string) (int, error)
	GetSession(id int) (*core.Session, error)
}

// Repository главная структура слоя репозиторий
type Repository struct {
	Authorization
	Session
}

// NewRepository конструктор
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Session:       NewSessionPostgres(db),
	}
}
