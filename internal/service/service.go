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

type TodoList interface {
	Create(userId int, list core.TodoList) (int, error)
	GetAll(userId int) ([]core.TodoList, error)
	GetById(userId, listId int) (core.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input core.UpdateListInput) error
}

type TodoItem interface {
	Create(userId, listId int, item core.TodoItem) (int, error)
	GetAll(userId, listId int) ([]core.TodoItem, error)
	GetById(userId, itemId int) (core.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input core.UpdateItemInput) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
