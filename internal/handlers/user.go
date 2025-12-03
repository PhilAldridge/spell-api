package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/PhilAldridge/spell-api/internal/dtos"
	"github.com/PhilAldridge/spell-api/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service:service}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body dtos.RegistrationRequest

	if err:= json.NewDecoder(r.Body).Decode(&body);err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	user, err:= h.service.Register(r.Context(), body)
	if err!=nil {
		http.Error(w, err.Message, err.StatusCode)

		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)	
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body dtos.LoginRequest

	if err:= json.NewDecoder(r.Body).Decode(&body);err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	res, err:= h.service.Login(r.Context(), body)
	if err!=nil {
		http.Error(w, err.Message, err.StatusCode)

		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) RefreshAccess(w http.ResponseWriter, r *http.Request) {
	var body dtos.RefreshAccessRequest

	if err:= json.NewDecoder(r.Body).Decode(&body);err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	res, err:= h.service.RefreshAccess(r.Context(),body.RefreshToken)
	if err != nil {
		http.Error(w, err.Message, err.StatusCode)

		return
	}

	json.NewEncoder(w).Encode(res)
}