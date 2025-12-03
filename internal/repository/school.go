package repository

import (
	"context"

	"github.com/PhilAldridge/spell-api/ent"
	"github.com/PhilAldridge/spell-api/ent/school"
	"github.com/PhilAldridge/spell-api/internal/apperrors"
)

type SchoolRepository struct {
}

func NewSchoolRepository() *SchoolRepository {
	return &SchoolRepository{}
}

func (r *SchoolRepository) Create(ctx context.Context, client *ent.Client, name string, ownerID int) (*ent.School, *apperrors.AppError) {
	school, err:= client.School.Create().
		SetName(name).SetOwnerID(ownerID).Save(ctx)

	if err !=nil {
		return nil, apperrors.ParseEntError(err,"could not create school")
	}

	return school,nil
}

func (r *SchoolRepository) GetByName(ctx context.Context, client *ent.Client, name string) (*ent.School, *apperrors.AppError) {
	school, err:= client.School.Query().
		Where(school.NameEqualFold(name)).
		Only(ctx)

	if err !=nil {
		return nil, apperrors.ParseEntError(err,"could not get school")
	}

	return school,nil
}