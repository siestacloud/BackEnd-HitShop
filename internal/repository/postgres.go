package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"tservice-checker/pkg"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
	clientsTable                       = "clients"                                // пользователи данного сервиса
	tAccountsTable                     = "telegram_accounts"                      // аккаунты через которые ведется вся работа
	tAccountsAdditionalAttributesTable = "telegram_account_additional_attributes" // доп инфа об аккаутах
	tUntrustSessionsTable              = "telegram_untrust_sessions"              // недоверенные сессии
	tTrustSessionsTable                = "telegram_trust_sessions"                // доверенные сессии (создаются на основе недоверенных)
	tAppsTable                         = "telegram_apps"                          // данные по регистрации телеграм-клиентов
	tGroupsTable                       = "telegram_groups"                        // группы телеграм по которым работают аккаунты
	tUsersTable                        = "telegram_users"                         // целевые пользователи телеги
	tGroupsUsersTable                  = "telegram_groups_users"                  // связующая таблица (группы с пользователями, многие ко многим)
)

// NewPostgresDB создание всех необходимых таблиц в БД
func NewPostgresDB(urlDB string) (*sqlx.DB, error) {
	if urlDB == "" {
		return nil, errors.New("url not set")
	}
	db, err := sqlx.Open("postgres", urlDB)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	logrus.Info("Success connect to postgres.")

	// делаем запрос на создание таблицы
	if err := createTable(db, clientsTable, "CREATE TABLE clients (id serial not null unique,login varchar(255) not null unique, password_hash varchar(255) not null);"); err != nil {
		return nil, err
	}
	// делаем запрос
	if err := createTable(db, tAccountsTable, "CREATE TABLE telegram_accounts (pk_telegram_account_id serial not null unique,account_id bigint not null unique, phone varchar(30),owner varchar(15) not null,status varchar(15) not null,username varchar(15) not null,firstname varchar(15) not null,lastname varchar(15) not null, create_time timestamp, delete_time timestamp);"); err != nil {
		return nil, err
	}
	// делаем запрос на создание таблицы
	if err := createTable(db, tAccountsAdditionalAttributesTable, "CREATE TABLE telegram_account_additional_attributes (pk_telegram_account_additional_attribute_id serial not null unique,fk_telegram_account_id int unique  REFERENCES telegram_accounts(pk_telegram_account_id) on delete cascade not null,bot boolean,fake boolean, scam boolean,support boolean,premium boolean,verified boolean, restricted boolean, restriction_reason text);"); err != nil {
		return nil, err
	}
	// делаем запрос на создание таблицы
	if err := createTable(db, tUntrustSessionsTable, "CREATE TABLE telegram_untrust_sessions (pk_telegram_untrust_session_id serial not null unique,fk_telegram_account_id int  REFERENCES telegram_accounts(pk_telegram_account_id) on delete cascade not null, data text not null, create_time timestamp, delete_time timestamp);"); err != nil {
		return nil, err
	}
	// делаем запрос на создание таблицы
	if err := createTable(db, tTrustSessionsTable, "CREATE TABLE telegram_trust_sessions (pk_telegram_trust_session_id serial not null unique,fk_telegram_account_id int  REFERENCES telegram_accounts(pk_telegram_account_id)on delete cascade not null,status varchar(15) not null,data text not null, create_time timestamp, delete_time timestamp);"); err != nil {
		return nil, err
	}
	// делаем запрос на создание таблицы
	if err := createTable(db, tAppsTable, "CREATE TABLE telegram_apps (pk_telegram_app_id serial not null unique,fk_telegram_account_id int unique REFERENCES telegram_accounts(pk_telegram_account_id)on delete cascade not null,app_id int ,app_hash varchar(255), create_time timestamp, delete_time timestamp);"); err != nil {
		return nil, err
	}

	// делаем запрос на создание таблицы
	if err := createTable(db, tGroupsTable, "CREATE TABLE telegram_groups (pk_telegram_group_id serial not null unique,fk_telegram_account_id int  REFERENCES telegram_accounts(pk_telegram_account_id));"); err != nil {
		return nil, err
	}

	// делаем запрос на создание таблицы
	if err := createTable(db, tUsersTable, "CREATE TABLE telegram_users(pk_telegram_user_id serial not null unique);"); err != nil {
		return nil, err
	}

	// делаем запрос на создание таблицы
	if err := createTable(db, tGroupsUsersTable, "CREATE TABLE telegram_groups_users (id serial not null unique,fk_telegram_group_id int REFERENCES telegram_groups(pk_telegram_group_id) not null,fk_telegram_user_id int REFERENCES telegram_users(pk_telegram_user_id) not null);"); err != nil {
		return nil, err
	}
	return db, nil
}

// * "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable"

func createTable(db *sqlx.DB, nameTable, query string) error {

	var checkExist bool

	row := db.QueryRow(fmt.Sprintf("SELECT EXISTS (SELECT FROM pg_tables WHERE  tablename  = '%s');", nameTable))
	if err := row.Scan(&checkExist); err != nil {
		return err
	}

	if !checkExist {
		_, err := db.Exec(query) //QueryRowContext т.к. одна запись
		if err != nil {
			return err
		}
		pkg.InfoPrint("repository", "ok", fmt.Sprintf("Table %45s successful create", nameTable))

	} else {
		pkg.WarnPrint("repository", "ok", fmt.Sprintf("Table %45s already created", nameTable))
	}

	return nil
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// INSERT INTO telegram_accounts (account_id,phone,owner,status,username,firstname,lastname) VALUES (1222322,'qwe','sps','ok','userNAME','fN','ls');
// INSERT INTO telegram_account_additional_attributes (fk_telegram_account_id,scam) VALUES (1,true);
// INSERT INTO telegram_trust_sessions (fk_telegram_trust_session_id,status,data) VALUES (1,'ok','1232dddd');

// delete from telegram_accounts where pk_telegram_account_id = 1;
