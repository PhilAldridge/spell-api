package repository

type Repository struct {
	RefreshTokenRepository *RefreshTokenRepository
	SchoolRepository       *SchoolRepository
	UserRepository         *UserRepository
}

func NewRepository() *Repository {
	return &Repository{
		RefreshTokenRepository: NewRefreshTokenRepository(),
		SchoolRepository: NewSchoolRepository(),
		UserRepository: NewUserRepository(),
	}
}
