package service

import (
	"hitshop/internal/core"
	"hitshop/internal/repository"
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

type Service struct {
	Authorization
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
	}
}
