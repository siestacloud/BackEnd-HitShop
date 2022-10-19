package repository

import (
	"errors"
	"fmt"
	"tservice-checker/internal/core"
	"tservice-checker/pkg"

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

	var pk_tAccountID int
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
	if err := row.Scan(&pk_tAccountID); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	var pk_tAccountAdditionalAttributeID int
	tAccountsAdditionalAttributesQuery := fmt.Sprintf("INSERT INTO %s (fk_telegram_account_id,bot,fake,scam,support,premium,verified) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING pk_telegram_account_additional_attribute_id", tAccountsAdditionalAttributesTable)
	row = tx.QueryRow(
		tAccountsAdditionalAttributesQuery,
		pk_tAccountID,
		tAccount.Bot,
		tAccount.Fake,
		tAccount.Scam,
		tAccount.Support,
		tAccount.Premium,
		tAccount.Verified,
	)
	if err := row.Scan(&pk_tAccountAdditionalAttributeID); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	var pk_tUntrustSessionID int
	for _, tSession := range tAccount.Sessions {
		tUntrustSessionQuery := fmt.Sprintf("INSERT INTO %s (fk_telegram_account_id,data,create_time) VALUES ($1,$2,$3) RETURNING pk_telegram_untrust_session_id", tUntrustSessionsTable)
		row = tx.QueryRow(
			tUntrustSessionQuery,
			pk_tAccountID,
			pkg.Base64Encode(tSession.Data),
			NewNullString(tAccount.CreateTime),
		)
		if err := row.Scan(&pk_tUntrustSessionID); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
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
