package repository

import (
	"errors"
	"fmt"
	"tservice-checker/internal/core"

	"github.com/jmoiron/sqlx"
)

//AuthPostgres реализует логику авторизации и аутентификации
type TAccountPostgres struct {
	db *sqlx.DB
}

//NewAuthPostgres конструктор
func NewTAccountPostgres(db *sqlx.DB) *SessionPostgres {
	return &SessionPostgres{db: db}
}

func (s *SessionPostgres) Save(tAccount *core.TelegramAccount) error {
	if s.db == nil {
		return errors.New("database are not connected")
	}
	// todo реализовать логику сохранения аккаунта в базе
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	// INSERT INTO telegram_accounts (account_id,phone,owner,status,username,firstname,lastname) VALUES (1222322,'qwe','sps','ok','userNAME','fN','ls');

	var id int
	tAccountQuery := fmt.Sprintf("INSERT INTO %s (account_id,phone,owner,status,username,firstname,lastname,create_time) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING pk_telegram_account_id", tAccountsTable)
	row := tx.QueryRow(
		tAccountQuery,
		tAccount.AccountID,
		tAccount.Phone,
		tAccount.Owner,
		tAccount.Status,
		tAccount.UserName,
		tAccount.FirstName,
		tAccount.LastName,
		NewNullString(tAccount.CreateTime),
	)
	if err := row.Scan(&id); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()

}

//GetUser получить пользователя из базы
func (r *SessionPostgres) Get(id int) (*core.TelegramAccount, error) {
	if r.db == nil {
		return nil, errors.New("database are not connected")
	}
	var tAccount core.TelegramAccount
	// todo реализовать логику извлечения сессии из базы
	return &tAccount, nil
}
