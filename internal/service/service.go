package service

import (
	"hitshop/internal/core"
	"hitshop/internal/repository"
	"mime/multipart"
)

// Authorization имплементорует логику авторизации
//
//go:generate mockgen -source=service.go -destination=mocks/mock.go
type Authorization interface {
	Test()
	CreateUser(user core.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

// Session имплементирует логику извлечения и проверки сессии
type TSession interface {
	ValidateSession(sess *core.Session) error
}

type TAccount interface {
	/* MultipartSave метод извлекает сессии из переданного слайса []*multipart.FileHeader. Работает только с архивами .zip
	1. разархивирует архивы;
	2. ищет в них tdata;
	3. извлекает сессии;
	4. валидирует сессии;
	5. создает обьект tAccount, добавляет валидные сессии в поле sessions
	6. сохраняет все данные обькта tAccount (включая валидные сессии) в базу;
	7. возвращает итог*/
	MultipartSave([]*multipart.FileHeader) (*core.ExtractResult, error)
}

type Service struct {
	Authorization
	TSession
	TAccount
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		TSession:      NewTSessionService(repo.TSession, repo.TClient),
		TAccount:      NewTAccountService(repo.TAccount, repo.TClient),
	}
}
