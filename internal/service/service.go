package service

import (
	"hitshop/internal/core"
	"hitshop/internal/repository"

	"github.com/google/uuid"
)

// Authorization имплементорует логику авторизации
//
//go:generate mockgen -source=service.go -destination=mocks/mock.go
type Authorization interface {
	Test()
	CreateAccount(acc *core.Account) (uuid.UUID, error)
	UpdateAccount(acc *core.Account) (uuid.UUID, error)
	GetAccountByEmail(email, password string) (*core.Account, error)
	GetAccountByCode(verification_code string) (*core.Account, error)
	ChangePassword(UUID uuid.UUID, password, passwordNew string) (uuid.UUID, error)

	GenerateToken(email, password string) (string, error)
	ParseToken(token string) (uuid.UUID, error)
}

type Service struct {
	Authorization
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
	}
}
