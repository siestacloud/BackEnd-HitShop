package service

import (
	"gitlab.com/siteasservice/project-architecture/templates/template-svc-golang/internal/core"
	"gitlab.com/siteasservice/project-architecture/templates/template-svc-golang/internal/repository"
)

type Authorization interface {
	Test()
	CreateUser(user core.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
