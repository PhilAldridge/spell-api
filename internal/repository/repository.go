package repository

import (
	"context"

	"github.com/PhilAldridge/spell-api/ent"
	"github.com/PhilAldridge/spell-api/internal/apperrors"
)

type Repository struct {
	RefreshTokenRepository *RefreshTokenRepository
	SchoolRepository       *SchoolRepository
	UserRepository         *UserRepository
	client *ent.Client
}

func NewRepository(client *ent.Client) *Repository {
	return &Repository{
		RefreshTokenRepository: NewRefreshTokenRepository(client),
		SchoolRepository: NewSchoolRepository(client),
		UserRepository: NewUserRepository(client),
		client: client,
	}
}

func Transaction[T any](ctx context.Context, repo *Repository, fn func(txRepo *Repository) (T, *apperrors.AppError)) (T,*apperrors.AppError) {
	var null T

	tx, err:= repo.client.Tx(ctx)
	if err != nil {
		return null, apperrors.ParseEntError(err,"unable to create transactions")
	}

	txRepo:= NewRepository(tx.Client())

	result, errFn:= fn(txRepo)
	if errFn != nil {
		if errRb:= tx.Rollback(); errRb != nil {
			return null, apperrors.ParseEntError(errRb, "unable to rollback transaction")
		}
		return null, errFn
	}

	if err:= tx.Commit(); err != nil {
		return null, apperrors.ParseEntError(err,"unable to commit transaction")
	}

	return result, nil
}
