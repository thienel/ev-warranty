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

func NewInternalServerError(err error) *AppError {
	return New(http.StatusInternalServerError, ErrorCodeInternalServerError, err)
}

func NewBadGateWay(err error) *AppError {
	return New(http.StatusBadGateway, ErrorCodeBadGateway, err)
}

func NewInvalidJsonRequest() *AppError {
	return New(http.StatusBadRequest, ErrorCodeInvalidJsonRequest, errors.New("invalid json request"))
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
	return New(http.StatusBadRequest, ErrorCodeInvalidCredentials, errors.New("invalid credentials"))
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

func NewUserNotFound() *AppError {
	return New(http.StatusNotFound, ErrorCodeUserNotFound, errors.New("user not found"))
}

func NewUserPasswordInvalid() *AppError {
	return New(http.StatusUnauthorized, ErrorCodeUserPasswordInvalid, errors.New("invalid user password"))
}

func NewDBOperationError(err error) *AppError {
	return New(http.StatusInternalServerError, ErrorCodeDBOperation, err)
}

func NewDBDuplicateKeyError(key string) *AppError {
	return New(http.StatusConflict, ErrorCodeDuplicateKey, errors.New(fmt.Sprintf("key %s already existed", key)))
}

func NewHashPasswordError(err error) *AppError {
	return New(http.StatusInternalServerError, ErrorCodeHashPassword, err)
}

func NewRefreshTokenNotFound() *AppError {
	return New(http.StatusNotFound, ErrorCodeRefreshTokenNotFound, errors.New("refresh token not found"))
}

func NewOfficeNotFound() *AppError {
	return New(http.StatusNotFound, ErrorCodeOfficeNotFound, errors.New("office not found"))
}

func NewClaimItemNotFound() *AppError {
	return New(http.StatusNotFound, ErrorCodeClaimItemNotFound, errors.New("claim item not found"))
}

func NewClaimHistoryNotFound() *AppError {
	return New(http.StatusNotFound, ErrorCodeClaimHistoryNotFound, errors.New("claim history not found"))
}

func NewClaimAttachmentNotFound() *AppError {
	return New(http.StatusNotFound, ErrorCodeClaimAttachmentNotFound, errors.New("claim attachment not found"))
}

func NewUserInactive() *AppError {
	return New(http.StatusForbidden, ErrorCodeUserInactive, errors.New("user inactive"))
}

func NewInvalidOfficeType() *AppError {
	return New(http.StatusBadRequest, ErrorCodeInvalidOfficeType, errors.New("invalid office type"))
}

func NewFailedHashToken() *AppError {
	return New(http.StatusInternalServerError, ErrorCodeFailedHashToken, errors.New("hash token failed"))
}

func NewUnexpectedSigningMethod(method any) *AppError {
	return New(http.StatusUnauthorized, ErrorCodeUnexpectedSigningMethod, errors.New(fmt.Sprintf("unexpected signing method: %v", method)))
}
