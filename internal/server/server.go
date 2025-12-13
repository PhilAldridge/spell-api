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

	repository := repository.NewRepository(client)

	service:= service.NewService(repository)

	userHandler:= handlers.NewUserHandler(service)

	r.Route("/users", func(r chi.Router) {
		r.Route("/test", func(r chi.Router) {
			r.Use(service.AuthMiddleware)
			r.Get("/", userHandler.Test)
		})
		r.Route("/logout", func(r chi.Router) {
			r.Use(service.AuthMiddleware)
			r.Get("/", userHandler.Test)
		})
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
		r.Post("/refresh", userHandler.RefreshAccess)
	})

	return r
}