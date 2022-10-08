package handler

import (
	"gitlab.com/siteasservice/project-architecture/templates/template-svc-golang/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
