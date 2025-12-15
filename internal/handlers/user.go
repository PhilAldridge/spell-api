package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/PhilAldridge/spell-api/internal/dtos"
	"github.com/PhilAldridge/spell-api/internal/service"
)

type UserHandler struct {
	service *service.Service
}

func NewUserHandler(service *service.Service) *UserHandler {
	return &UserHandler{service:service}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body dtos.RegistrationRequest

	if err:= json.NewDecoder(r.Body).Decode(&body);err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	user, err:= h.service.UserService.Register(r.Context(), body)
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

	res, err:= h.service.UserService.Login(r.Context(), body)
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

	res, err:= h.service.UserService.RefreshAccess(r.Context(),body.RefreshToken)
	if err != nil {
		http.Error(w, err.Message, err.StatusCode)

		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	err := h.service.UserService.Logout(r.Context())
	if err != nil {
		http.Error(w, err.Message, err.StatusCode)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type Test struct {
	Hello string `json:"hello"`
}

func (h *UserHandler) Test(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Test{Hello: "hello"})
}