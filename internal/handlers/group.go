package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/PhilAldridge/spell-api/internal/dtos"
	"github.com/PhilAldridge/spell-api/internal/service"
)

type groupHandler struct {
	service *service.Service
}

func NewGroupHandler(service *service.Service) *groupHandler {
	return &groupHandler{
		service: service,
	}
}

func (h *groupHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body dtos.GroupCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	group, err := h.service.GroupService.Create(r.Context(), body)
	if err != nil {
		http.Error(w, err.Message, err.StatusCode)

		return
	}

	json.NewEncoder(w).Encode(group.ID)
}

func (h *groupHandler) Update(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *groupHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// TODO
}
