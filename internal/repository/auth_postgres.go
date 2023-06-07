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
func (r *AuthPostgres) CreateAccount(acc *core.Account) (uuid.UUID, error) {
	if r.db == nil {
		return uuid.UUID{}, errors.New("database are not connected")
	}
	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO %s (account_email, account_password_hash,account_verify,account_verify_code,account_create_at,account_update_at) values ($1, $2, $3,$4,$5,$6) RETURNING pk_account_id", accountsTable)

	row := r.db.QueryRow(query, acc.Email, acc.Password, acc.Verified, acc.VerificationCode, acc.CreateAt, acc.UpdateAt)
	if err := row.Scan(&id); err != nil {
		pkg.ErrPrintR("unhealfy", err)
		return uuid.UUID{}, errors.New("login busy")
	}

	return id, nil
}

// GetUser получить пользователя из базы
func (r *AuthPostgres) GetAccountByEmail(email, password string) (*core.Account, error) {
	if r.db == nil {
		return nil, errors.New("database are not connected")
	}
	var acc core.Account
	query := fmt.Sprintf("SELECT pk_account_id, account_email, account_verify, account_password_hash,  account_create_at, account_update_at   FROM %s WHERE account_email=$1 AND account_password_hash=$2", accountsTable)
	if err := r.db.Get(&acc, query, email, password); err != nil {
		pkg.ErrPrintR("unhealfy", err)

		return nil, errors.New("invalid username/password pair")
	}

	return &acc, nil
}

// GetUser получить пользователя из базы
func (r *AuthPostgres) GetAccountByCode(verification_code string) (*core.Account, error) {
	if r.db == nil {
		return nil, errors.New("database are not connected")
	}

	var acc core.Account
	query := fmt.Sprintf("SELECT  pk_account_id, account_email, account_verify, account_password_hash,  account_create_at, account_update_at  FROM %s WHERE account_verify_code=$1", accountsTable)
	if err := r.db.Get(&acc, query, verification_code); err != nil {
		pkg.ErrPrintR("unhealfy", err)
		return nil, errors.New("Invalid verification code or account doesn't exists")
	}

	return &acc, nil
}

// GetUser получить пользователя из базы
func (r *AuthPostgres) GetAccountByUUID(UUID uuid.UUID) (*core.Account, error) {
	if r.db == nil {
		return nil, errors.New("database are not connected")
	}
	var acc core.Account
	query := fmt.Sprintf("SELECT pk_account_id, account_email, account_verify, account_password_hash,  account_create_at, account_update_at   FROM %s WHERE pk_account_id=$1 ", accountsTable)
	if err := r.db.Get(&acc, query, UUID); err != nil {
		pkg.ErrPrintR("unhealfy", err)

		return nil, errors.New("invalid username/password pair")
	}

	return &acc, nil
}

// GetUser получить пользователя из базы
func (r *AuthPostgres) UpdateAccount(acc *core.Account) (uuid.UUID, error) {
	if r.db == nil {
		return uuid.UUID{}, errors.New("database are not connected")
	}
	// fk_role_id_referer  fk_account_stat_id_referer account_email account_verify account_password_hash account_phone_number account_create_at account_update_at account_delete_at account_verify_code
	query := fmt.Sprintf("UPDATE %s SET account_verify=$1, account_password_hash=$2,account_update_at=$3 WHERE account_email = $4  ", accountsTable)
	_, err := r.db.Exec(query, acc.Verified, acc.Password, acc.UpdateAt, acc.Email)
	if err != nil {
		pkg.ErrPrintR("unhealfy", err)
		return uuid.UUID{}, errors.New("login busy")
	}
	return uuid.UUID{}, nil
}
