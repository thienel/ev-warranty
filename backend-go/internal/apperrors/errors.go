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
	return New(http.StatusConflict, ErrorCodeDuplicateKey, fmt.Errorf("key %s already existed", key))
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
	return New(http.StatusUnauthorized, ErrorCodeUnexpectedSigningMethod,
		fmt.Errorf("unexpected signing method: %v", method))
}

func NewInvalidClaimAction() *AppError {
	return New(http.StatusConflict, ErrorCodeInvalidClaimAction, errors.New("invalid claim action"))
}

func NewNotAllowUpdateClaim() *AppError {
	return New(http.StatusConflict, ErrorCodeClaimStatusNotAllowedUpdate,
		errors.New("current claim's status does not allow to update"))
}

func NewNotAllowDeleteClaim() *AppError {
	return New(http.StatusConflict, ErrorCodeClaimStatusNotAllowedDelete,
		errors.New("current claim's status does not allow to delete"))
}

func NewMissingInformationClaim() *AppError {
	return New(http.StatusConflict, ErrorCodeClaimMissingInformation,
		errors.New("claim does not have enough information to submit"))
}

func NewInvalidQueryParameter(param string) *AppError {
	return New(http.StatusBadRequest, ErrorCodeInvalidQueryParameter,
		fmt.Errorf("invalid query parameter: %s", param))
}

func NewMissingUserID() *AppError {
	return New(http.StatusUnauthorized, ErrorCodeMissingUserID, errors.New("missing X-User-ID header"))
}

func NewInvalidUserID() *AppError {
	return New(http.StatusBadRequest, ErrorCodeInvalidUserID, errors.New("invalid user ID format"))
}

func NewFailedInitializeCloudinary() *AppError {
	return New(http.StatusInternalServerError, ErrorCodeFailedInitializeCloudinary, errors.New("failed to initialize Cloudinary"))
}

func NewFailedUploadCloudinary() *AppError {
	return New(http.StatusServiceUnavailable, ErrorCodeFailedUploadCloudinary, errors.New("failed to upload Cloudinary"))
}

func NewFailedDeleteCloudinary() *AppError {
	return New(http.StatusServiceUnavailable, ErrorCodeFailedDeleteCloudinary, errors.New("failed to delete Cloudinary"))
}

func NewInvalidMultipartFormRequest() *AppError {
	return New(http.StatusBadRequest, ErrorCodeInvalidMultipartFormRequest, errors.New("invalid multipart form request"))
}
