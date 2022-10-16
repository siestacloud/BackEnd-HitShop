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
	SaveZip(file *multipart.FileHeader) (string, error) //* сохраняю полученный архив
	Unzip(src string) (string, error)                   //* разархивирую из архива директорию tdata
	ExtractSession(src string) ([]core.Session, error)  //* вытаскувую из директории tdata сессию
	ValidateSession(session *core.Session) error        //* проверяю жива ли сессия, сохраняю ее в базе
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
