package service

import (
	"github.com/PhilAldridge/spell-api/ent"
	"github.com/PhilAldridge/spell-api/internal/repository"
)

type Service struct {
	UserService UserService
}

func NewService(repository *repository.Repository, client *ent.Client) *Service {
	return &Service{
		UserService: *NewUserService(repository, client),
	}
}