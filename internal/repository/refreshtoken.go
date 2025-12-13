package repository

import (
	"context"
	"time"

	"github.com/PhilAldridge/spell-api/ent"
	"github.com/PhilAldridge/spell-api/ent/refreshtoken"
	"github.com/PhilAldridge/spell-api/internal/apperrors"
)

type RefreshTokenRepository struct {
	client *ent.Client
}

func NewRefreshTokenRepository (client *ent.Client) *RefreshTokenRepository {
	return &RefreshTokenRepository{client: client}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, tokenHash string, expiresAt time.Time, user *ent.User) *apperrors.AppError {
	_,err:= r.client.RefreshToken.Create().
		SetTokenHash(tokenHash).
		SetExpiresAt(expiresAt).
		SetUser(user).
		Save(ctx)
	
	if err != nil {
		return apperrors.ParseEntError(err, "could not create refresh token")
	}

	return nil
}

func (r *RefreshTokenRepository) IsValid(ctx context.Context, tokenHash string, userID int) (*apperrors.AppError) {
	refreshToken, err:= r.client.RefreshToken.Query().
		Where(
			refreshtoken.TokenHashEQ(tokenHash),
			refreshtoken.RevokedEQ(false),
			refreshtoken.ExpiresAtGT(time.Now()),
		).
		Only(ctx)
	if err!=nil {
		return apperrors.ParseEntError(err, "invalid refresh token")
	}

	if refreshTokenUser,err:= refreshToken.QueryUser().Only(ctx); err != nil || refreshTokenUser.ID != userID {
		return apperrors.ParseEntError(err, "invalid refresh token")
	}

	return nil
}

// func (r *RefreshTokenRepository) Update(ctx context.Context, )