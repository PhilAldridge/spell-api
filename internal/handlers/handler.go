package handlers

import "github.com/PhilAldridge/spell-api/internal/service"

type Handler struct {
	UserHandler  *userHandler
	GroupHandler *groupHandler
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		UserHandler:  NewUserHandler(service),
		GroupHandler: NewGroupHandler(service),
	}
}
