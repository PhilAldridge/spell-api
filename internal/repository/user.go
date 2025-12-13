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
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client: client}
}

// Create a new user 
func (r *UserRepository) CreateUser(ctx context.Context, u *ent.User, groupIDs []int, schoolIDs []int) (*ent.User, *apperrors.AppError) {
    query:= r.client.User.Create().
        SetName(u.Name).
        SetEmail(u.Email).
        SetPasswordHash(u.PasswordHash).
		AddGroupIDs(groupIDs...).
		AddSchoolIDs(schoolIDs...)

	if u.AccountType != "" {
		query = query.SetAccountType(u.AccountType)
	}

	user, err := query.Save(ctx)
	if err !=nil {
		return nil, apperrors.ParseEntError(err, "unable to create user")
	}

    return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*ent.User, *apperrors.AppError) {
	user,err:= r.client.User.Get(ctx,id)
	if err !=nil {
		return nil, apperrors.ParseEntError(err,"unable to get user")
	}

	return user,nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*ent.User, *apperrors.AppError) {
	user, err:= r.client.User.Query().
		Where(user.EmailEQ(email)).
		Only(ctx)
	
	if err!=nil {
		return nil, apperrors.ParseEntError(err, "unable to get user: ")
	}

	return user, nil
}
