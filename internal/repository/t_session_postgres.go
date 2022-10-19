package repository

import (
	"errors"
	"tservice-checker/internal/core"

	"github.com/jmoiron/sqlx"
)

//AuthPostgres реализует логику авторизации и аутентификации
type SessionPostgres struct {
	db *sqlx.DB
}

//NewAuthPostgres конструктор
func NewSessionPostgres(db *sqlx.DB) *SessionPostgres {
	return &SessionPostgres{db: db}
}

//CreateUser создание пользователя
func (r *SessionPostgres) SaveSession(session string) (int, error) {
	if r.db == nil {
		return 0, errors.New("database are not connected")
	}
	var id int
	// todo реализовать логику сохранения сессии в базе
	return id, nil
}

//GetUser получить пользователя из базы
func (r *SessionPostgres) GetSession(id int) (*core.Session, error) {
	if r.db == nil {
		return nil, errors.New("database are not connected")
	}
	var session core.Session
	// todo реализовать логику извлечения сессии из базы
	return &session, nil
}
