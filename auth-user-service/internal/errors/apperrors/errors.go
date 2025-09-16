package apperrors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func NewBadRequest(message string) *AppError {
	return NewAppError(http.StatusBadRequest, message, nil)
}

func NewUnauthorized(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, message, nil)
}

func NewForbidden(message string) *AppError {
	return NewAppError(http.StatusForbidden, message, nil)
}

func NewNotFound(message string) *AppError {
	return NewAppError(http.StatusNotFound, message, nil)
}

func NewConflict(message string) *AppError {
	return NewAppError(http.StatusConflict, message, nil)
}

func NewBadGateway(message string) *AppError {
	return NewAppError(http.StatusBadGateway, message, nil)
}

func NewInternal(message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, message, err)
}

var (
	ErrScanValue = NewInternal("scan value error", nil)

	ErrUserInactive = NewForbidden("user is inactive")

	ErrInvalidAccessToken = NewUnauthorized("invalid access token")
	ErrExpiredAccessToken = NewUnauthorized("access token expired")

	ErrInvalidRefreshToken         = NewUnauthorized("invalid refresh token")
	ErrExpiredRefreshToken         = NewUnauthorized("refresh token expired")
	ErrRevokedRefreshToken         = NewUnauthorized("refresh token revoked")
	ErrInvalidAuthenticationHeader = NewUnauthorized("missing or invalid authorization header")

	ErrInvalidJSONRequest = NewBadRequest("invalid json request")
)

func ErrInvalidCredentials(message string) *AppError {
	return NewUnauthorized(message)
}

func ErrNotFound(entity string) *AppError {
	message := fmt.Sprintf("%s not found", entity)
	return NewNotFound(message)
}

func ErrDuplicateKey(dup string) *AppError {
	message := fmt.Sprintf("duplicate key value violates unique constraint %s", dup)
	return NewConflict(message)
}

func ErrDBOperation(err error) *AppError {
	return NewInternal("database operation error", err)
}

func ErrRequestTimeout(err error) *AppError {
	return NewInternal("request timeout or canceled", err)
}

func ErrHashPassword(err error) *AppError {
	return NewInternal("hashing password error", err)
}

func ErrFailedSignAccessToken(err error) *AppError {
	return NewInternal("failed to sign access token", err)
}

func ErrFailedGenerateRefreshToken(err error) *AppError {
	return NewInternal("failed to generate refresh token", err)
}
