package repository

import (
	"context"
	"time"

	"github.com/PhilAldridge/spell-api/ent"
	"github.com/PhilAldridge/spell-api/ent/result"
	"github.com/PhilAldridge/spell-api/ent/user"
	"github.com/PhilAldridge/spell-api/internal/apperrors"
)

type ListParams struct {
	SchoolID      *int
	GroupID       *int
	CompetitionID *int
}

type userRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *userRepository {
	return &userRepository{client: client}
}

// Create a new user
func (r *userRepository) CreateUser(ctx context.Context, u *ent.User) (*ent.User, *apperrors.AppError) {
	query := r.client.User.Create().
		SetName(u.Name).
		SetEmail(u.Email).
		SetPasswordHash(u.PasswordHash)

	if u.AccountType != "" {
		query = query.SetAccountType(u.AccountType)
	}

	user, err := query.Save(ctx)
	if err != nil {
		return nil, apperrors.ParseEntError(err, "unable to create user")
	}

	return user, nil
}

func (r *userRepository) GetStudentByID(ctx context.Context, id int) (*ent.User, *apperrors.AppError) {
	user, err := r.client.User.Query().
		Where(user.IDEQ(id), user.AccountTypeEQ(user.AccountTypeStudent)).
		WithGroups(func(gq *ent.GroupQuery) {
			gq.WithSchool()
		}).
		WithResults(func(rq *ent.ResultQuery) {
			rq.GroupBy(
				result.WordColumn,
			).Aggregate(
				ent.Count(),
				ent.Sum(result.FieldCorrect),
				ent.Mean(result.FieldTimeTakenInSeconds),
			)
			result.TestedAtTimestampGT(time.Now().UTC().AddDate(0, -2, 0))
		}).
		First(ctx)
	if err != nil {
		return nil, apperrors.ParseEntError(err, "unable to get user")
	}

	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*ent.User, *apperrors.AppError) {
	user, err := r.client.User.Query().
		Where(user.EmailEQ(email)).
		Only(ctx)

	if err != nil {
		return nil, apperrors.ParseEntError(err, "unable to get user: ")
	}

	return user, nil
}

func (r *userRepository) JoinGroup(ctx context.Context, userID int, groupID int) *apperrors.AppError {
	err := r.client.User.UpdateOneID(userID).AddGroupIDs(groupID).Exec(ctx)
	if err != nil {
		return apperrors.ParseEntError(err, "could not join group")
	}

	return nil
}

func (r *userRepository) JoinSchool(ctx context.Context, userID int, schoolID int) *apperrors.AppError {
	err := r.client.User.UpdateOneID(userID).AddSchoolIDs(schoolID).Exec(ctx)
	if err != nil {
		return apperrors.ParseEntError(err, "could not join school")
	}

	return nil
}

func (r *userRepository) LeaveGroup(ctx context.Context, userID int, groupID int) *apperrors.AppError {
	err := r.client.User.UpdateOneID(userID).RemoveGroupIDs(groupID).Exec(ctx)
	if err != nil {
		return apperrors.ParseEntError(err, "could not Leave group")
	}

	return nil
}

func (r *userRepository) LeaveSchool(ctx context.Context, userID int, schoolID int) *apperrors.AppError {
	err := r.client.User.UpdateOneID(userID).RemoveSchoolIDs(schoolID).Exec(ctx)
	if err != nil {
		return apperrors.ParseEntError(err, "could not Leave school")
	}

	return nil
}
