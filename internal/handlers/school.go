package handlers

import (
	"net/http"

	"github.com/PhilAldridge/spell-api/internal/service"
)

type schoolHandler struct {
	service *service.Service
}

func NewSchoolHandler(service *service.Service) *schoolHandler {
	return &schoolHandler{
		service: service,
	}
}

func (s *schoolHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	//TODO	
}

func (s *schoolHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (s *schoolHandler) Update(w http.ResponseWriter, r *http.Request) {
	// TODO
}