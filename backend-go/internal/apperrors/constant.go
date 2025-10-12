package apperrors

const (
	ErrorCodeInternalServerError = "COMMON_INTERNAL_ERROR"

	ErrorCodeDBOperation  = "DB_OPERATION_ERROR"
	ErrorCodeDuplicateKey = "DB_DUPLICATE_KEY"
	ErrorCodeHashPassword = "DB_HASH_PASSWORD_ERROR"

	ErrorCodeInvalidJsonRequest = "API_INVALID_JSON_REQUEST"

	ErrorCodeInvalidAccessToken         = "AUTH_INVALID_ACCESS_TOKEN"
	ErrorCodeExpiredAccessToken         = "AUTH_EXPIRED_ACCESS_TOKEN"
	ErrorCodeFailedHashToken            = "AUTH_FAILED_HASH_TOKEN"
	ErrorCodeInvalidRefreshToken        = "AUTH_INVALID_REFRESH_TOKEN"
	ErrorCodeExpiredRefreshToken        = "AUTH_EXPIRED_REFRESH_TOKEN"
	ErrorCodeRevokedRefreshToken        = "AUTH_REVOKED_REFRESH_TOKEN"
	ErrorCodeInvalidAuthHeader          = "AUTH_INVALID_AUTH_HEADER"
	ErrorCodeInvalidCredentials         = "AUTH_INVALID_CREDENTIALS"
	ErrorCodeFailedSignAccessToken      = "AUTH_FAILED_SIGN_ACCESS_TOKEN"
	ErrorCodeFailedGenerateRefreshToken = "AUTH_FAILED_GENERATE_REFRESH_TOKEN"
	ErrorCodeUnexpectedSigningMethod    = "AUTH_UNEXPECTED_SIGNING_METHOD"

	ErrorCodeRefreshTokenNotFound = "REFRESH_TOKEN_NOT_FOUND"

	ErrorCodeInvalidOfficeType = "OFFICE_INVALID_TYPE"
	ErrorCodeOfficeNotFound    = "OFFICE_NOT_FOUND"

	ErrorCodeUserNotFound        = "USER_NOT_FOUND"
	ErrorCodeUserInactive        = "USER_INACTIVE"
	ErrorCodeUserPasswordInvalid = "USER_PASSWORD_INVALID"

	ErrorCodeClaimItemNotFound       = "CLAIM_ITEM_NOT_FOUND"
	ErrorCodeClaimHistoryNotFound    = "CLAIM_HISTORY_NOT_FOUND"
	ErrorCodeClaimAttachmentNotFound = "CLAIM_ATTACHMENT_NOT_FOUND"
	ErrorCodeClaimNotFound           = "CLAIM_NOT_FOUND"

	ErrorCodeClaimStatusNotAllowedUpdate = "CLAIM_STATUS_NOT_ALLOWED_TO_UPDATE"
	ErrorCodeClaimStatusNotAllowedDelete = "CLAIM_STATUS_NOT_ALLOWED_TO_DELETE"
	ErrorCodeInvalidClaimAction          = "CLAIM_INVALID_ACTION"
)
