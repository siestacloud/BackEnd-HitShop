package repository

import (
	"tservice-checker/internal/core"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	TestDB()
	CreateUser(user core.User) (int, error)
	GetUser(username, password string) (*core.User, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
