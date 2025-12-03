package server

import (
	"net/http"

	"github.com/PhilAldridge/spell-api/ent"
	"github.com/PhilAldridge/spell-api/internal/handlers"
	"github.com/PhilAldridge/spell-api/internal/repository"
	"github.com/PhilAldridge/spell-api/internal/service"
	"github.com/go-chi/chi/v5"
)

func New(client *ent.Client) http.Handler {
	r:= chi.NewRouter()

	userRepo:= repository.NewUserRepository()
	schoolRepo:= repository.NewSchoolRepository()
	refreshTokenRepo:= repository.NewRefreshTokenRepository()

	userService:= service.NewUserService(userRepo,schoolRepo,refreshTokenRepo,client)

	userHandler:= handlers.NewUserHandler(userService)

	r.Route("/users", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
		r.Post("/refresh", userHandler.RefreshAccess)
	})

	return r
}