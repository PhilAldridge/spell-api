package repository

import (
	"context"

	"github.com/PhilAldridge/spell-api/ent"
	"github.com/PhilAldridge/spell-api/ent/user"
	"github.com/PhilAldridge/spell-api/internal/apperrors"
)

type ListParams struct {
	SchoolID *int
	GroupID *int
	CompetitionID *int
}

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Create a new user 
func (r *UserRepository) CreateUser(ctx context.Context, client *ent.Client, u *ent.User, groupIDs []int, schoolIDs []int) (*ent.User, *apperrors.AppError) {
    user,err:= client.User.Create().
        SetName(u.Name).
        SetEmail(u.Email).
        SetPasswordHash(u.PasswordHash).
		SetAccountType(u.AccountType).
		AddGroupIDs(groupIDs...).
		AddSchoolIDs(schoolIDs...).
        Save(ctx)
	if err !=nil {
		return nil, apperrors.ParseEntError(err, "unable to create user")
	}

    return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, client *ent.Client, id int) (*ent.User, *apperrors.AppError) {
	user,err:= client.User.Get(ctx,id)
	if err !=nil {
		return nil, apperrors.ParseEntError(err,"unable to get user")
	}

	return user,nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, client *ent.Client, email string) (*ent.User, *apperrors.AppError) {
	user, err:= client.User.Query().
		Where(user.EmailEQ(email)).
		Only(ctx)
	
	if err!=nil {
		return nil, apperrors.ParseEntError(err, "unable to get user: ")
	}

	return user, nil
}
