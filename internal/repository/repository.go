package repository

import "github.com/PhilAldridge/spell-api/ent"

type Repository struct {
	RefreshTokenRepository *RefreshTokenRepository
	SchoolRepository       *SchoolRepository
	UserRepository         *UserRepository
}

func NewRepository(client *ent.Client) *Repository {
	return &Repository{
		RefreshTokenRepository: NewRefreshTokenRepository(client),
		SchoolRepository: NewSchoolRepository(client),
		UserRepository: NewUserRepository(client),
	}
}
