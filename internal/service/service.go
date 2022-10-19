package service

import (
	"mime/multipart"
	"tservice-checker/internal/core"
	"tservice-checker/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
// Authorization имплементорует логику авторизации
type Authorization interface {
	Test()
	CreateUser(user core.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

// Session имплементирует логику извлечения и проверки сессии
type Session interface {
	/* Extract метод извлекает сессии из переданного слайса []*multipart.FileHeader. Работает только с архивами .zip
	1. разархивирует архивы;
	2. ищет в них tdata;
	3. извлекает сессии;
	4. валидирует сессии;
	5. сохраняет сессии в базе;
	6. возвращает итог*/
	Extract([]*multipart.FileHeader) (*core.ExtractResult, error)
}

type Service struct {
	Authorization
	Session
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Session:       NewSessionService(repos.Session),
	}
}
