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
    userRepo *repository.UserRepository
    schoolRepo *repository.SchoolRepository
    refreshTokenRepo *repository.RefreshTokenRepository
    client *ent.Client
}

func NewUserService(
    userRepo *repository.UserRepository,
    schoolRepo *repository.SchoolRepository,
    refeshTokenRepo *repository.RefreshTokenRepository,
    client *ent.Client,
    ) *UserService {
    return &UserService{
        userRepo:userRepo,
        schoolRepo: schoolRepo,
        refreshTokenRepo: refeshTokenRepo,
        client: client}
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

    tx, errTx:= s.client.Tx(ctx);
    if errTx != nil {
        return nil, apperrors.Internal("unable to start transaction")
    }

    user:= ent.User{
        Name: req.Name,
        PasswordHash: passwordHash,
        Email: req.Email,
        AccountType: user.AccountType(*req.AccountType),
    }
    
    newUser, err:= s.userRepo.CreateUser(ctx, tx.Client(),&user, groupIDs,schoolIDs)
    if err !=nil {
        tx.Rollback()

        return nil, err
    }

    if req.NewSchoolName != nil {
        _, err:= s.schoolRepo.Create(ctx, tx.Client(), req.Name, newUser.ID)
        if err != nil {
            tx.Rollback()

            return nil,err
        }
    }

    commitErr := tx.Commit()

    if commitErr != nil {
        return newUser, apperrors.Internal(commitErr.Error())
    }

    return newUser, nil
}

func (s *UserService) Login(ctx context.Context, req dtos.LoginRequest) (*dtos.LoginResponse, *apperrors.AppError) {
    //Validate request

    user, err := s.userRepo.GetUserByEmail(ctx, s.client, req.Email)
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

    err = s.refreshTokenRepo.Create(ctx, s.client, tokenHash,expiresAt,user)
    if err!= nil {
        //TODO: implement proper logging
        fmt.Println(err)
    }

    return &dtos.LoginResponse{
        RefreshToken: refeshToken,
        AccessToken: accessToken,
        ExpiresIn: accessExpiryMins*60,
    },nil
}

func (s *UserService) RefreshAccess(ctx context.Context, refreshToken string) (*dtos.RefreshAccessResponse,*apperrors.AppError) {
    user, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, apperrors.BadRequest("could not find user information")
	}

    hash := auth.HashRefreshToken(refreshToken)

    err := s.refreshTokenRepo.IsValid(ctx, s.client, hash, user.ID)
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