package repository

import (
	"context"

	"github.com/PhilAldridge/spell-api/ent"
	"github.com/PhilAldridge/spell-api/ent/school"
	"github.com/PhilAldridge/spell-api/internal/apperrors"
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
		Only(ctx)

	if err !=nil {
		return nil, apperrors.ParseEntError(err,"could not get school")
	}

	return school,nil
}