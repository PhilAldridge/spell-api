package handlers

import "github.com/PhilAldridge/spell-api/internal/service"

type groupHandler struct {
	service *service.Service
}

func NewGroupHandler(service *service.Service) *groupHandler {
	return &groupHandler{
		service: service,
	}
}
