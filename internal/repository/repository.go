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

// Session имплементирует логику хранения телеграм сессий в базе
type TSession interface {
	SaveSession(session string) (int, error)
	GetSession(id int) (*core.Session, error)
}

type TAccount interface {
	/* Save метод сохраняет телеграмм аккаунт в базу
	1. основная инфа об аккаунте сохр в telegram_accounts;
	2. дополнительная инфа об аккаунте сохр в telegram_account_additional_attributes;
	3. валидная недоверенная сессия этого аккаунта сохр в telegram_untrust_sessions;
	4. метод записывает данные в режиме транзакции (все или ничего)*/
	Save(tAccount *core.TelegramAccount) error
}
type TClient interface {
	ValidateTSession(tSession *core.Session) error
	GetAccountInfo(tSession *core.Session) (*core.TelegramAccount, error)
}

// Repository главная структура слоя репозиторий
type Repository struct {
	Authorization
	TSession
	TAccount
	TClient
}

// NewRepository конструктор
func NewRepository(db *sqlx.DB, cfg *config.Cfg) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TSession:      NewSessionPostgres(db),
		TAccount:      NewTAccountPostgres(db),
		TClient:       NewTClientAPI(cfg),
	}
}
