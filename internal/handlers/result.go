package handlers

import (
	"net/http"

	"github.com/PhilAldridge/spell-api/internal/service"
)

type resultHandler struct {
	service *service.Service
}

func NewResultHandler(service *service.Service) *resultHandler {
	return &resultHandler{
		service: service,
	}
}

func (h *resultHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO
}