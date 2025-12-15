package handlers

import "github.com/PhilAldridge/spell-api/internal/service"

type Handler struct {
	UserHandler  *userHandler
	GroupHandler *groupHandler
	SchoolHandler *schoolHandler
	ResultHandler *resultHandler
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		UserHandler:  NewUserHandler(service),
		GroupHandler: NewGroupHandler(service),
		SchoolHandler: NewSchoolHandler(service),
		ResultHandler: NewResultHandler(service),
	}
}
