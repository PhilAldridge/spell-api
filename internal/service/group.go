package service

import (
	"context"
	"slices"

	"github.com/PhilAldridge/spell-api/ent"
	"github.com/PhilAldridge/spell-api/ent/user"
	"github.com/PhilAldridge/spell-api/internal/apperrors"
	"github.com/PhilAldridge/spell-api/internal/auth"
	"github.com/PhilAldridge/spell-api/internal/dtos"
	"github.com/PhilAldridge/spell-api/internal/repository"
	"github.com/PhilAldridge/spell-api/internal/utils"
)

type GroupService struct {
	repository *repository.Repository
}

func NewGroupService(repository *repository.Repository) *GroupService {
	return &GroupService{
		repository: repository,
	}
}

func (s *GroupService) Create(ctx context.Context, req dtos.GroupCreateRequest) (*ent.Group, *apperrors.AppError) {
	userObject, ok:= auth.UserFromContext(ctx)
	if !ok {
		return nil, apperrors.BadRequest("could not find user information")
	}

	if userObject.AccountType != user.AccountTypeAdmin {
		return nil, apperrors.Unauthorised("only admin accounts may create groups")
	}

	school, err := s.repository.SchoolRepository.GetByID(ctx, req.SchoolID)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(utils.ExtractIDsReflect(school.Edges.Admins), userObject.ID) {
		return nil, apperrors.NotFound("could not find school attached to this account")
	}

	group, err:= s.repository.GroupRepository.Create(ctx, req.Name, req.SchoolID)
	if err != nil {
		return nil, err
	}

	return group, err
}