package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/PhilAldridge/spell-api/ent"
	"github.com/PhilAldridge/spell-api/ent/user"
	"github.com/PhilAldridge/spell-api/internal/apperrors"
	"github.com/PhilAldridge/spell-api/internal/auth"
	"github.com/PhilAldridge/spell-api/internal/dtos"
	"github.com/PhilAldridge/spell-api/internal/repository"
)

type userService struct {
	repository *repository.Repository
}

func NewUserService(
	repository *repository.Repository,
) *userService {
	return &userService{repository: repository}
}

func (s *userService) Register(ctx context.Context, req dtos.RegistrationRequest) (*ent.User, *apperrors.AppError) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	passwordHash, errHash := auth.HashPassword(req.Password)
	if errHash != nil {
		return nil, apperrors.Internal("password hash failed")
	}

	newUser, err := repository.Transaction(ctx, s.repository, func(txRepo *repository.Repository) (*ent.User, *apperrors.AppError) {
		userObject := ent.User{
			Name:         req.Name,
			PasswordHash: passwordHash,
			Email:        req.Email,
		}

		if req.AccountType != nil {
			userObject.AccountType = user.AccountType(*req.AccountType)
		}

		newUser, err := txRepo.UserRepository.CreateUser(ctx, &userObject)
		if err != nil {
			return nil, err
		}

		if req.NewSchoolName != nil {
			_, err := txRepo.SchoolRepository.Create(ctx, req.Name, newUser.ID)
			if err != nil {
				return nil, err
			}
		}

		return newUser, nil
	})

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *userService) Login(ctx context.Context, req dtos.LoginRequest) (*dtos.LoginResponse, *apperrors.AppError) {
	//TODO Validate request

	user, err := s.repository.UserRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.Unauthorised("invalid credentials")
	}

	if err := auth.ComparePassword(user.PasswordHash, req.Password); err != nil {
		return nil, apperrors.Unauthorised("invalid credentials")
	}

	accessExpiryMins := 15
	accessToken, errToken := auth.GenerateAccessToken(user.ID, accessExpiryMins, string(user.AccountType))
	if errToken != nil {
		return nil, apperrors.Internal("could not generate access token")
	}

	tokenRaw := make([]byte, 64)
	if _, err := rand.Read(tokenRaw); err != nil {
		return nil, apperrors.Internal("could not generate refresh token")
	}

	refeshToken := hex.EncodeToString(tokenRaw)
	tokenHash := auth.HashRefreshToken(refeshToken)
	expiresAt := time.Now().Add(24 * 30 * time.Hour)

	err = s.repository.RefreshTokenRepository.Create(ctx, tokenHash, expiresAt, user)
	if err != nil {
		//TODO: implement proper logging
		fmt.Println(err)
	}

	// TODO Also return user information
	return &dtos.LoginResponse{
		RefreshToken: refeshToken,
		AccessToken:  accessToken,
		ExpiresIn:    accessExpiryMins * 60,
	}, nil
}

func (s *userService) Logout(ctx context.Context) *apperrors.AppError {
	user, ok := auth.UserFromContext(ctx)
	if !ok {
		return apperrors.BadRequest("could not find user information")
	}

	return s.repository.RefreshTokenRepository.Revoke(ctx, user.ID)
}

func (s *userService) RefreshAccess(ctx context.Context, refreshToken string) (*dtos.RefreshAccessResponse, *apperrors.AppError) {
	user, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, apperrors.BadRequest("could not find user information")
	}

	hash := auth.HashRefreshToken(refreshToken)

	err := s.repository.RefreshTokenRepository.IsValid(ctx, hash, user.ID)
	if err != nil {
		return nil, err
	}

	accessExpiryMins := 15
	accessToken, errToken := auth.GenerateAccessToken(user.ID, accessExpiryMins, string(user.AccountType))
	if errToken != nil {
		return nil, apperrors.Internal("could not generate access token")
	}

	return &dtos.RefreshAccessResponse{
		AccessToken: accessToken,
		ExpiresIn:   accessExpiryMins * 60,
	}, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int) (*ent.User, *apperrors.AppError) {
	user, err := s.repository.UserRepository.GetStudentByID(ctx, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *userService) JoinGroupOrSchool(ctx context.Context, joinCode string) *apperrors.AppError {
	userObject, ok := auth.UserFromContext(ctx)
	if !ok {
		return apperrors.BadRequest("could not find user information")
	}

	if userObject.AccountType == user.AccountTypeAdmin {
		return s.joinSchool(ctx, userObject.ID, joinCode)
	}

	return s.joinGroup(ctx, userObject.ID, joinCode)
}

func (s *userService) joinGroup(ctx context.Context, userID int, joinCode string) *apperrors.AppError {
	group, err := s.repository.GroupRepository.GetByJoinCode(ctx, joinCode)
	if err != nil {
		return err
	}

	err = s.repository.UserRepository.JoinGroup(ctx, userID, group.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) joinSchool(ctx context.Context, userID int, joinCode string) *apperrors.AppError {
	school, err := s.repository.SchoolRepository.GetByJoinCode(ctx, joinCode)
	if err != nil {
		return err
	}

	err = s.repository.UserRepository.JoinSchool(ctx, userID, school.ID)
	if err != nil {
		return err
	}

	return nil
}
