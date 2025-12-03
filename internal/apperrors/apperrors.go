package apperrors

import (
	"net/http"

	"github.com/PhilAldridge/spell-api/ent"
)

type AppError struct {
	StatusCode int
	Message    string
}

func BadRequest(msg string) *AppError {
	return &AppError{StatusCode: http.StatusBadRequest, Message: msg}
}

func NotFound(msg string) *AppError {
	return &AppError{StatusCode: http.StatusNotFound, Message: msg}
}

func Conflict(msg string) *AppError {
	return &AppError{StatusCode: http.StatusConflict, Message: msg}
}

func Internal(msg string) *AppError {
	return &AppError{StatusCode: http.StatusInternalServerError, Message: msg}
}

func Unauthorised(msg string) *AppError {
	return &AppError{StatusCode: http.StatusUnauthorized, Message: msg}
}

func ParseEntError(err error, msg string) *AppError {
	errMsg := msg+": "+err.Error()

	if ent.IsConstraintError(err) {
		return Conflict(errMsg)
	}

	if ent.IsNotFound(err) {
		return NotFound(errMsg)
	}

	return Internal(errMsg)
}
