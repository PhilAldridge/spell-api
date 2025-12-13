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

type UserService struct {
    repository *repository.Repository
}

func NewUserService(
    repository *repository.Repository,
    ) *UserService {
    return &UserService{repository: repository}
}

func (s *UserService) Register(ctx context.Context, req dtos.RegistrationRequest) (*ent.User, *apperrors.AppError) {
    //validate req

    groupIDs:= []int{}
    schoolIDs:=[]int{}
    if req.JoinCode != nil {
        //check join code exists
    }

    passwordHash, errHash := auth.HashPassword(req.Password)
    if errHash!=nil {
        return nil, apperrors.Internal("password hash failed")
    }

    //TODO
    // tx, errTx:= s.client.Tx(ctx);
    // if errTx != nil {
    //     return nil, apperrors.Internal("unable to start transaction")
    // }

    userObject:= ent.User{
        Name: req.Name,
        PasswordHash: passwordHash,
        Email: req.Email,
    }

    if req.AccountType != nil {
        userObject.AccountType = user.AccountType(*req.AccountType)
    }
    
    newUser, err:= s.repository.UserRepository.CreateUser(ctx,&userObject, groupIDs,schoolIDs)
    if err !=nil {
        // tx.Rollback()

        return nil, err
    }

    if req.NewSchoolName != nil {
        _, err:= s.repository.SchoolRepository.Create(ctx, req.Name, newUser.ID)
        if err != nil {
            // tx.Rollback()

            return nil,err
        }
    }

    // commitErr := tx.Commit()

    // if commitErr != nil {
    //     return newUser, apperrors.Internal(commitErr.Error())
    // }

    return newUser, nil
}

func (s *UserService) Login(ctx context.Context, req dtos.LoginRequest) (*dtos.LoginResponse, *apperrors.AppError) {
    //Validate request

    user, err := s.repository.UserRepository.GetUserByEmail(ctx, req.Email)
    if err != nil {
        return nil, apperrors.Unauthorised("invalid credentials")
    }

    if err:= auth.ComparePassword(user.PasswordHash,req.Password); err !=nil {
        return nil, apperrors.Unauthorised("invalid credentials")
    }

    accessExpiryMins := 15
    accessToken, errToken := auth.GenerateAccessToken(user.ID, accessExpiryMins, string(user.AccountType))
    if errToken != nil {
        return nil, apperrors.Internal("could not generate access token")
    }

    tokenRaw:= make([]byte, 64)
    if _,err:= rand.Read(tokenRaw); err!=nil {
        return nil, apperrors.Internal("could not generate refresh token")
    }

    refeshToken:= hex.EncodeToString(tokenRaw)
    tokenHash:= auth.HashRefreshToken(refeshToken)
    expiresAt:= time.Now().Add(24*30*time.Hour)

    err = s.repository.RefreshTokenRepository.Create(ctx, tokenHash,expiresAt,user)
    if err!= nil {
        //TODO: implement proper logging
        fmt.Println(err)
    }

    // TODO Also return user information
    return &dtos.LoginResponse{
        RefreshToken: refeshToken,
        AccessToken: accessToken,
        ExpiresIn: accessExpiryMins*60,
    },nil
}

func (s *UserService) Logout(ctx context.Context, req dtos.LoginRequest) (*apperrors.AppError) {
    user, ok := auth.UserFromContext(ctx)
	if !ok {
		return apperrors.BadRequest("could not find user information")
	}
    fmt.Println(user)

    // TODO
    // Revoke token, returning error if not available


    return nil
}

func (s *UserService) RefreshAccess(ctx context.Context, refreshToken string) (*dtos.RefreshAccessResponse,*apperrors.AppError) {
    user, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, apperrors.BadRequest("could not find user information")
	}

    hash := auth.HashRefreshToken(refreshToken)

    err := s.repository.RefreshTokenRepository.IsValid(ctx, hash, user.ID)
    if err!=nil {
        return nil, err
    }

    accessExpiryMins := 15
    accessToken, errToken := auth.GenerateAccessToken(user.ID, accessExpiryMins, string(user.AccountType))
    if errToken != nil {
        return nil, apperrors.Internal("could not generate access token")
    }

    return &dtos.RefreshAccessResponse{
        AccessToken: accessToken,
        ExpiresIn: accessExpiryMins*60,
    },nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*ent.User, *apperrors.AppError) {
    user,err:= s.repository.UserRepository.GetUserByID(ctx, id)
    if err != nil {
        return user,err
    }

    return user,nil
}