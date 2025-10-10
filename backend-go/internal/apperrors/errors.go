package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	HttpCode  int    `json:"-"`
	ErrorCode string `json:"error_code"`
	Err       error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.ErrorCode, e.Err)
	}
	return e.ErrorCode
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(httpCode int, errorCode string, err error) *AppError {
	return &AppError{
		HttpCode:  httpCode,
		ErrorCode: errorCode,
		Err:       err,
	}
}

func NewBadRequest(err error) *AppError {
	return New(http.StatusBadRequest, ErrorCodeBadRequest, err)
}

func NewInternalServerError(err error) *AppError {
	return New(http.StatusInternalServerError, ErrorCodeInternalServerError, err)
}

func NewNotFound(err error) *AppError {
	return New(http.StatusNotFound, ErrorCodeNotFound, err)
}

func NewConflict(err error) *AppError {
	return New(http.StatusConflict, ErrorCodeConflict, err)
}

func NewUnauthorized(err error) *AppError {
	return New(http.StatusUnauthorized, ErrorCodeUnauthorized, err)
}

func NewForbidden(err error) *AppError {
	return New(http.StatusForbidden, ErrorCodeForbidden, err)
}

func NewTimeout(err error) *AppError {
	return New(http.StatusRequestTimeout, ErrorCodeTimeout, err)
}

func NewInvalidAccessToken() *AppError {
	return New(http.StatusUnauthorized, ErrorCodeInvalidAccessToken, errors.New("invalid access token"))
}

func NewExpiredAccessToken() *AppError {
	return New(http.StatusUnauthorized, ErrorCodeExpiredAccessToken, errors.New("expired access token"))
}

func NewInvalidRefreshToken() *AppError {
	return New(http.StatusUnauthorized, ErrorCodeInvalidRefreshToken, errors.New("invalid refresh token"))
}

func NewExpiredRefreshToken() *AppError {
	return New(http.StatusUnauthorized, ErrorCodeExpiredRefreshToken, errors.New("expired refresh token"))
}

func NewRevokedRefreshToken() *AppError {
	return New(http.StatusUnauthorized, ErrorCodeRevokedRefreshToken, errors.New("revoked refresh token"))
}

func NewInvalidAuthHeader() *AppError {
	return New(http.StatusUnauthorized, ErrorCodeInvalidAuthHeader, errors.New("invalid authorization header"))
}

func NewInvalidCredentials() *AppError {
	return New(http.StatusUnauthorized, ErrorCodeInvalidCredentials, errors.New("invalid credentials"))
}

func NewFailedSignAccessToken(err error) *AppError {
	return New(http.StatusInternalServerError, ErrorCodeFailedSignAccessToken, err)
}

func NewFailedGenerateRefreshToken(err error) *AppError {
	return New(http.StatusInternalServerError, ErrorCodeFailedGenerateRefreshToken, err)
}

func NewClaimNotFound() *AppError {
	return New(http.StatusNotFound, ErrorCodeClaimNotFound, errors.New("claim not found"))
}

func NewClaimInvalidStatus() *AppError {
	return New(http.StatusBadRequest, ErrorCodeClaimInvalidStatus, errors.New("invalid claim status"))
}

func NewClaimDeleteFailed(err error) *AppError {
	return New(http.StatusInternalServerError, ErrorCodeClaimDeleteFailed, err)
}

func NewClaimCreateFailed(err error) *AppError {
	return New(http.StatusInternalServerError, ErrorCodeClaimCreateFailed, err)
}

func NewUserNotFound() *AppError {
	return New(http.StatusNotFound, ErrorCodeUserNotFound, errors.New("user not found"))
}

func NewUserAlreadyExists() *AppError {
	return New(http.StatusConflict, ErrorCodeUserAlreadyExists, errors.New("user already exists"))
}

func NewUserPasswordInvalid() *AppError {
	return New(http.StatusUnauthorized, ErrorCodeUserPasswordInvalid, errors.New("invalid user password"))
}

func NewDBOperationError(err error) *AppError {
	return New(http.StatusInternalServerError, ErrorCodeDBOperation, err)
}

func NewDBDuplicateKeyError(err error) *AppError {
	return New(http.StatusConflict, ErrorCodeDuplicateKey, err)
}

func NewHashPasswordError(err error) *AppError {
	return New(http.StatusInternalServerError, ErrorCodeHashPassword, err)
}
