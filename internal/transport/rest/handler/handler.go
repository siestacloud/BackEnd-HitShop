package handler

import (
	"hitshop/internal/config"
	"hitshop/internal/service"
)

type Handler struct {
	services *service.Service
	cfg      *config.Cfg
}

func NewHandler(services *service.Service, cfg *config.Cfg) *Handler {
	return &Handler{services: services, cfg: cfg}
}
