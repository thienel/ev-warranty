package apperror

import (
	"fmt"
	"net/http"
)

var (
	ErrInternalServerError = New(http.StatusInternalServerError, "CMN_INTERNAL_ERROR",
		"Internal server error")
	ErrNotFoundError = New(http.StatusNotFound, "CMN_NOT_FOUND_ERROR", "Record not found")

	ErrExternalServiceMarshalRequest = New(http.StatusInternalServerError,
		"XTN_SERVICE_MARSHAL_REQUEST", "Failed to marshal request payload")
	ErrExternalServiceRequestFailed = New(http.StatusBadGateway, "XTN_SERVICE_REQUEST_FAILED",
		"External service request failed")
	ErrExternalServiceUnavailable = New(http.StatusServiceUnavailable,
		"XTN_SERVICE_UNAVAILABLE", "External service is unavailable")
	ErrExternalServiceReadResponse = New(http.StatusBadGateway, "XTN_SERVICE_READ_RESPONSE",
		"Failed to read external service response")
	ErrExternalServiceUnmarshalResponse = New(http.StatusBadGateway,
		"XTN_SERVICE_UNMARSHAL_RESPONSE",
		"Failed to unmarshal external service response")

	ErrInvalidParams        = New(http.StatusBadRequest, "API_INVALID_PARAMS", "Invalid param url")
	ErrInvalidInput         = New(http.StatusBadRequest, "API_INVALID_INPUT", "Invalid input")
	ErrInvalidMultipartForm = New(http.StatusBadRequest, "API_INVALID_MULTIPART_FORM_REQUEST",
		"Invalid multipart form request")
	ErrInvalidFile = New(http.StatusBadRequest, "API_INVALID_FILE",
		"Invalid or unreadable file")
	ErrInvalidJsonRequest = New(http.StatusBadRequest, "API_INVALID_JSON_REQUEST",
		"Invalid JSON request")

	ErrDBOperation = New(http.StatusInternalServerError, "DB_OPERATION_ERROR",
		"Database operation failed")
	ErrDuplicateKey = New(http.StatusConflict, "DB_DUPLICATE_KEY", "Key already exists")

	ErrInvalidAccessToken = New(http.StatusUnauthorized, "AUTH_INVALID_ACCESS_TOKEN",
		"Invalid access token")
	ErrExpiredAccessToken = New(http.StatusUnauthorized, "AUTH_EXPIRED_ACCESS_TOKEN",
		"Expired access token")
	ErrInvalidRefreshToken = New(http.StatusUnauthorized, "AUTH_INVALID_REFRESH_TOKEN",
		"Invalid refresh token")
	ErrExpiredRefreshToken = New(http.StatusUnauthorized, "AUTH_EXPIRED_REFRESH_TOKEN",
		"Expired refresh token")
	ErrRevokedRefreshToken = New(http.StatusUnauthorized, "AUTH_REVOKED_REFRESH_TOKEN",
		"Revoked refresh token")
	ErrInvalidAuthHeader = New(http.StatusUnauthorized, "AUTH_INVALID_AUTH_HEADER",
		"Invalid authorization header")
	ErrInvalidCredentials = New(http.StatusUnauthorized, "AUTH_INVALID_CREDENTIALS",
		"Email or password is not correct")
	ErrFailedSignAccessToken = New(http.StatusInternalServerError,
		"AUTH_FAILED_SIGN_ACCESS_TOKEN", "Failed to sign access token")
	ErrFailedGenerateRefreshToken = New(http.StatusInternalServerError,
		"AUTH_FAILED_GENERATE_REFRESH_TOKEN", "Failed to generate refresh token")
	ErrFailedHashToken = New(http.StatusInternalServerError, "AUTH_FAILED_HASH_TOKEN",
		"Failed to hash token")
	ErrUnexpectedSigningMethod = New(http.StatusUnauthorized, "AUTH_UNEXPECTED_SIGNING_METHOD",
		"Unexpected signing method")
	ErrInvalidUserID = New(http.StatusUnauthorized, "AUTH_INVALID_USER_ID",
		"Invalid user id in claim")

	ErrUserInactive  = New(http.StatusForbidden, "USER_INACTIVE", "User is inactive")
	ErrMissingUserID = New(http.StatusBadRequest, "AUTH_MISSING_USER_ID",
		"Missing X-User-ID header")
	ErrMissingUserRole = New(http.StatusBadRequest, "AUTH_MISSING_USER_ROLE",
		"Missing X-User-Role header")
	ErrUnauthorizedRole = New(http.StatusForbidden, "AUTH_UNAUTHORIZED_ROLE",
		"User does not have permission to perform this action")

	ErrHashPassword = New(http.StatusInternalServerError, "DB_HASH_PASSWORD_ERROR",
		"Failed to hash password")

	ErrInvalidClaimAction = New(http.StatusConflict, "CLAIM_INVALID_ACTION",
		"Invalid claim action")
	ErrMissingInformationClaim = New(http.StatusBadRequest, "CLAIM_MISSING_INFORMATION",
		"Claim does not have enough information to submit")

	ErrFailedInitializeCloudinary = New(http.StatusInternalServerError,
		"CLOUDINARY_FAILED_INITIALIZE", "Failed to initialize Cloudinary")
	ErrInvalidCloudinaryURL = New(http.StatusBadRequest, "CLOUDINARY_INVALID_URL",
		"Invalid Cloudinary URL")
	ErrFailedUploadCloudinary = New(http.StatusServiceUnavailable, "CLOUDINARY_FAILED_UPLOAD",
		"Failed to upload to Cloudinary")
	ErrFailedDeleteCloudinary = New(http.StatusServiceUnavailable, "CLOUDINARY_FAILED_DELETE",
		"Failed to delete from Cloudinary")
	ErrEmptyCloudinaryParameter = New(http.StatusBadRequest, "CLOUDINARY_EMPTY_PARAMETER",
		"Empty Cloudinary parameter")
)

func NewExternalServiceRequestFailed(statusCode int, err error) *AppError {
	return ErrExternalServiceRequestFailed.
		WithMessage(fmt.Sprintf("External service returned status %d", statusCode)).
		WithError(err)
}
