package repository

import (
	"context"
	"time"

	"github.com/PhilAldridge/spell-api/ent"
	"github.com/PhilAldridge/spell-api/ent/school"
	"github.com/PhilAldridge/spell-api/internal/apperrors"
	"github.com/PhilAldridge/spell-api/internal/utils"
)

type SchoolRepository struct {
	client *ent.Client
}

func NewSchoolRepository(client *ent.Client) *SchoolRepository {
	return &SchoolRepository{client: client}
}

func (r *SchoolRepository) Create(ctx context.Context, name string, ownerID int) (*ent.School, *apperrors.AppError) {
	school, err:= r.client.School.Create().
		SetName(name).SetOwnerID(ownerID).Save(ctx)

	if err !=nil {
		return nil, apperrors.ParseEntError(err,"could not create school")
	}

	return school,nil
}

func (r *SchoolRepository) GetByName(ctx context.Context, name string) (*ent.School, *apperrors.AppError) {
	school, err:= r.client.School.Query().
		Where(school.NameEqualFold(name)).
		WithGroups(func(gq *ent.GroupQuery) {
			gq.WithUsers()
		}).
		Only(ctx)

	if err !=nil {
		return nil, apperrors.ParseEntError(err,"could not get school")
	}

	return school,nil
}

func (r *SchoolRepository) GetByID(ctx context.Context, id int) (*ent.School, *apperrors.AppError) {
	school, err:= r.client.School.Query().
		Where(school.IDEQ(id)).
		WithGroups(func(gq *ent.GroupQuery) {
			gq.WithUsers()
		}).
		Only(ctx)

	if err !=nil {
		return nil, apperrors.ParseEntError(err,"could not get school")
	}

	return school,nil
}

func (r *SchoolRepository) GetByJoinCode(ctx context.Context, joinCode string) (*ent.School, *apperrors.AppError) {
	school, err:= r.client.School.Query().
		Where(
			school.JoinCodeEqualFold(joinCode),
			school.JoinCodeValidUntilTimestampGTE(time.Now().UTC()),
		).
		Only(ctx)

	if err !=nil {
		return nil, apperrors.ParseEntError(err,"could not get school")
	}

	return school,nil
}

func (r *SchoolRepository) RefreshJoinCode(ctx context.Context, id int) (*ent.School, *apperrors.AppError) {
	str, err:= utils.RandomString()
	if err != nil {
		return nil, apperrors.Internal("could not create join code string")
	}

	err = r.client.School.UpdateOneID(id).
		SetJoinCode(str).
		SetJoinCodeValidUntilTimestamp(time.Now().UTC().AddDate(0,0,14)).
		Exec(ctx)

		if err !=nil {
		return nil, apperrors.ParseEntError(err,"could not update join code")
	}

	return r.GetByID(ctx, id)
}

