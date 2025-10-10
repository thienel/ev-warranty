package apperrors

const (
	ErrorCodeBadRequest          = "COMMON_BAD_REQUEST"
	ErrorCodeInternalServerError = "COMMON_INTERNAL_ERROR"
	ErrorCodeNotFound            = "COMMON_NOT_FOUND"
	ErrorCodeConflict            = "COMMON_CONFLICT"
	ErrorCodeUnauthorized        = "COMMON_UNAUTHORIZED"
	ErrorCodeForbidden           = "COMMON_FORBIDDEN"
	ErrorCodeTimeout             = "COMMON_TIMEOUT"

	ErrorCodeInvalidAccessToken         = "AUTH_INVALID_ACCESS_TOKEN"
	ErrorCodeExpiredAccessToken         = "AUTH_EXPIRED_ACCESS_TOKEN"
	ErrorCodeInvalidRefreshToken        = "AUTH_INVALID_REFRESH_TOKEN"
	ErrorCodeExpiredRefreshToken        = "AUTH_EXPIRED_REFRESH_TOKEN"
	ErrorCodeRevokedRefreshToken        = "AUTH_REVOKED_REFRESH_TOKEN"
	ErrorCodeInvalidAuthHeader          = "AUTH_INVALID_AUTH_HEADER"
	ErrorCodeInvalidCredentials         = "AUTH_INVALID_CREDENTIALS"
	ErrorCodeFailedSignAccessToken      = "AUTH_FAILED_SIGN_ACCESS_TOKEN"
	ErrorCodeFailedGenerateRefreshToken = "AUTH_FAILED_GENERATE_REFRESH_TOKEN"

	ErrorCodeClaimNotFound      = "CLAIM_NOT_FOUND"
	ErrorCodeClaimInvalidStatus = "CLAIM_INVALID_STATUS"
	ErrorCodeClaimDeleteFailed  = "CLAIM_DELETE_FAILED"
	ErrorCodeClaimCreateFailed  = "CLAIM_CREATE_FAILED"

	ErrorCodeUserNotFound        = "USER_NOT_FOUND"
	ErrorCodeUserAlreadyExists   = "USER_ALREADY_EXISTS"
	ErrorCodeUserPasswordInvalid = "USER_PASSWORD_INVALID"

	ErrorCodeDBOperation  = "DB_OPERATION_ERROR"
	ErrorCodeDuplicateKey = "DB_DUPLICATE_KEY"
	ErrorCodeHashPassword = "DB_HASH_PASSWORD_ERROR"
)
