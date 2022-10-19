package repository

import (
	"tservice-checker/internal/config"
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
type TSession interface {
	SaveSession(session string) (int, error)
	GetSession(id int) (*core.Session, error)
}

type TAccount interface {
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
