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
	CreateUser(user core.Account) (uuid.UUID, error)
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
