package repository

import (
	"context"
	"time"

	"github.com/PhilAldridge/spell-api/ent"
	"github.com/PhilAldridge/spell-api/ent/group"
	"github.com/PhilAldridge/spell-api/internal/apperrors"
	"github.com/PhilAldridge/spell-api/internal/utils"
)

type groupRepository struct {
	client *ent.Client
}

func NewGroupRepository(client *ent.Client) *groupRepository {
	return &groupRepository{client: client}
}

func (r *groupRepository) Create(ctx context.Context, name string, schoolID int) (*ent.Group, *apperrors.AppError) {
	group, err := r.client.Group.Create().
		SetName(name).SetSchoolID(schoolID).Save(ctx)

	if err != nil {
		return nil, apperrors.ParseEntError(err, "could not create group")
	}

	return group, nil
}

func (r *groupRepository) GetByName(ctx context.Context, name string) (*ent.Group, *apperrors.AppError) {
	group, err := r.client.Group.Query().
		Where(group.NameEqualFold(name)).
		WithUsers().
		WithCompetitions().
		WithWordLists(func(wlq *ent.WordListQuery) {
			wlq.WithWords()
		}).
		Only(ctx)

	if err != nil {
		return nil, apperrors.ParseEntError(err, "could not get group")
	}

	return group, nil
}

func (r *groupRepository) GetByID(ctx context.Context, id int) (*ent.Group, *apperrors.AppError) {
	school, err := r.client.Group.Query().
		Where(group.IDEQ(id)).
		WithUsers().
		WithCompetitions().
		WithWordLists(func(wlq *ent.WordListQuery) {
			wlq.WithWords()
		}).
		Only(ctx)

	if err != nil {
		return nil, apperrors.ParseEntError(err, "could not get school")
	}

	return school, nil
}

func (r *groupRepository) GetByJoinCode(ctx context.Context, joinCode string) (*ent.Group, *apperrors.AppError) {
	group, err := r.client.Group.Query().
		Where(
			group.JoinCodeEqualFold(joinCode),
			group.JoinCodeValidUntilTimestampGTE(time.Now().UTC()),
		).
		Only(ctx)

	if err != nil {
		return nil, apperrors.ParseEntError(err, "could not get group")
	}

	return group, nil
}

func (r *groupRepository) RefreshJoinCode(ctx context.Context, id int) (*ent.Group, *apperrors.AppError) {
	str, err := utils.RandomString()
	if err != nil {
		return nil, apperrors.Internal("could not create join code string")
	}

	err = r.client.Group.UpdateOneID(id).
		SetJoinCode(str).
		SetJoinCodeValidUntilTimestamp(time.Now().UTC().AddDate(0, 0, 14)).
		Exec(ctx)

	if err != nil {
		return nil, apperrors.ParseEntError(err, "could not update join code")
	}

	return r.GetByID(ctx, id)
}
