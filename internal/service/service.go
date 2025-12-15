package service

import (
	"github.com/PhilAldridge/spell-api/internal/repository"
)

type Service struct {
	UserService *userService
	GroupService *groupService
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		UserService: NewUserService(repository),
		GroupService: NewGroupService(repository),
	}
}
