package repository

import (
	"context"
	"time"

	"github.com/PhilAldridge/spell-api/ent"
	"github.com/PhilAldridge/spell-api/internal/apperrors"
)

type resultRepository struct {
	client *ent.Client
}

func NewResultRepository(client *ent.Client) *resultRepository {
	return &resultRepository{
		client: client,
	}
}

func (r *resultRepository) Create(ctx context.Context, result *ent.Result, userID int, wordID int) *apperrors.AppError {
	err := r.client.Result.Create().
		SetCorrect(result.Correct).
		SetTestedAtTimestamp(time.Now().UTC()).
		SetUserID(userID).
		SetType(result.Type).
		SetWordID(wordID).
		Exec(ctx)

	if err != nil {
		return apperrors.ParseEntError(err, "error saving result")
	}

	return nil
}
