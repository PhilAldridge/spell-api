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

	handler:= handlers.NewHandler(service)

	r.Route("/users", func(r chi.Router) {
		r.Route("/join", func(r chi.Router) {
			r.Use(service.AuthMiddleware)
			r.Get("/", handler.UserHandler.Test)
			r.Post("/", handler.UserHandler.JoinGroupOrSchool)
		})
		r.Route("/logout", func(r chi.Router) {
			r.Use(service.AuthMiddleware)
			r.Get("/", handler.UserHandler.Test)
		})
		r.Post("/register", handler.UserHandler.Register)
		r.Post("/login", handler.UserHandler.Login)
		r.Post("/refresh", handler.UserHandler.RefreshAccess)
	})

	r.Route("/school/{schoolId}", func(r chi.Router) {
			//Groups, words, wordlists, users, competitions
	})

	r.Route("/teacher/school/{schoolId}", func(r chi.Router) {
		//Groups, words, wordlists, users, competitions
	})

	return r
}