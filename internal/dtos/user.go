package dtos

import (
	"github.com/PhilAldridge/spell-api/ent/user"
	"github.com/PhilAldridge/spell-api/internal/apperrors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type RegistrationRequest struct {
	Name          string  `json:"name"`
	Password      string  `json:"password"`
	Email         string  `json:"email"`
	AccountType   *string `json:"account_type,omitempty"`
	NewSchoolName *string `json:"new_school_name,omitempty"`
}

func (r RegistrationRequest) Validate() *apperrors.AppError {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1,30)),
		validation.Field(&r.Password, validation.Required, validation.Length(8,255)),
		validation.Field(&r.Email, validation.Required, is.EmailFormat),
		validation.Field(&r.AccountType,validation.By(func(value interface{}) error {
			return user.AccountTypeValidator(value.(user.AccountType))
		})),
		validation.Field(&r.NewSchoolName, validation.Length(0,255)),
	)

	if err != nil {
		return apperrors.ParseValidationError(err, "request validation failed")
	}
		
	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r LoginRequest) Validate() *apperrors.AppError {
	err:= validation.ValidateStruct(&r,
		validation.Field(&r.Password, validation.Required, validation.Length(8,255)),
		validation.Field(&r.Email, validation.Required, is.EmailFormat),
	)

	if err != nil {
		return apperrors.ParseValidationError(err, "request validation failed")
	}
		
	return nil
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in_seconds"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshAccessRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshAccessResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in_seconds"`
}

type JoinRequest struct {
	JoinCode string `json:"join_code"`
}