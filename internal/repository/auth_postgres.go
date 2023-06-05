package repository

import (
	"errors"
	"fmt"
	"hitshop/internal/core"
	"hitshop/pkg"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/sirupsen/logrus"
)

// AuthPostgres реализует логику авторизации и аутентификации
type AuthPostgres struct {
	db *sqlx.DB
}

// NewAuthPostgres конструктор
func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

// Тестирование доступности слоя repository
func (r *AuthPostgres) TestDB() {
	logrus.Info("Info from DB layer")
}

// CreateUser создание пользователя
func (r *AuthPostgres) CreateAccount(acc core.Account) (uuid.UUID, error) {
	if r.db == nil {
		return uuid.UUID{}, errors.New("database are not connected")
	}
	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO %s (account_email, account_password_hash,account_verify) values ($1, $2, $3) RETURNING pk_account_id", accountsTable)

	row := r.db.QueryRow(query, acc.Email, acc.Password, acc.Verify)
	if err := row.Scan(&id); err != nil {
		pkg.ErrPrintR("repository", 409, err)
		return uuid.UUID{}, errors.New("login busy")
	}

	return id, nil
}

// GetUser получить пользователя из базы
func (r *AuthPostgres) GetAccount(email, password string) (*core.Account, error) {
	if r.db == nil {
		return nil, errors.New("database are not connected")
	}
	// найденный пользователь, парсится в обьект структуры, далее он возвращается на уровень выше
	var acc core.Account
	query := fmt.Sprintf("SELECT pk_account_id FROM %s WHERE account_email=$1 AND account_password_hash=$2", accountsTable)
	if err := r.db.Get(&acc, query, email, password); err != nil {
		pkg.ErrPrintR("repository", 409, err)
		return nil, errors.New("invalid username/password pair")
	}

	return &acc, nil
}
