const ErrorMessages = {
  // Common Errors
  COMMON_INTERNAL_ERROR: "An internal server error occurred. Please try again later.",

  // Database Errors
  DB_OPERATION_ERROR: "A database operation failed. Please try again.",
  DB_DUPLICATE_KEY: "This record already exists in the system.",
  DB_HASH_PASSWORD_ERROR: "Failed to process password. Please try again.",

  // API Request Errors
  API_INVALID_JSON_REQUEST: "Invalid request format. Please check your input.",
  API_INVALID_QUERY_PARAMETER: "Invalid query parameter provided.",
  API_INVALID_MULTIPART_FORM_REQUEST: "Invalid file upload request. Please check your files.",
  API_INVALID_UUID: "Invalid ID format provided.",

  // Authentication Errors
  AUTH_INVALID_ACCESS_TOKEN: "Invalid access token. Please log in again.",
  AUTH_EXPIRED_ACCESS_TOKEN: "Your session has expired. Please log in again.",
  AUTH_FAILED_HASH_TOKEN: "Failed to process authentication token.",
  AUTH_INVALID_REFRESH_TOKEN: "Invalid refresh token. Please log in again.",
  AUTH_EXPIRED_REFRESH_TOKEN: "Your session has expired. Please log in again.",
  AUTH_REVOKED_REFRESH_TOKEN: "Your session has been revoked. Please log in again.",
  AUTH_INVALID_AUTH_HEADER: "Invalid authorization header. Please log in again.",
  AUTH_INVALID_CREDENTIALS: "Invalid email or password. Please try again.",
  AUTH_FAILED_SIGN_ACCESS_TOKEN: "Failed to generate access token. Please try again.",
  AUTH_FAILED_GENERATE_REFRESH_TOKEN: "Failed to generate refresh token. Please try again.",
  AUTH_UNEXPECTED_SIGNING_METHOD: "Unexpected authentication method encountered.",
  AUTH_MISSING_USER_ID: "User ID is missing. Please log in again.",
  AUTH_INVALID_USER_ID: "Invalid user ID format.",

  // Refresh Token Errors
  REFRESH_TOKEN_NOT_FOUND: "Session not found. Please log in again.",

  // Office Errors
  OFFICE_INVALID_TYPE: "Invalid office type provided.",
  OFFICE_NOT_FOUND: "Office not found.",

  // User Errors
  USER_NOT_FOUND: "User not found.",
  USER_INACTIVE: "This user account is inactive. Please contact an administrator.",
  USER_PASSWORD_INVALID: "Invalid password. Please try again.",
  USER_INVALID_INPUT: "Invalid user information provided. Please check your input.",

  // Claim Errors
  CLAIM_ITEM_NOT_FOUND: "Claim item not found.",
  CLAIM_HISTORY_NOT_FOUND: "Claim history not found.",
  CLAIM_ATTACHMENT_NOT_FOUND: "Claim attachment not found.",
  CLAIM_NOT_FOUND: "Claim not found.",

  // Claim Status Errors
  CLAIM_STATUS_NOT_ALLOWED_TO_UPDATE: "This claim cannot be updated in its current status.",
  CLAIM_STATUS_NOT_ALLOWED_TO_DELETE: "This claim cannot be deleted in its current status.",
  CLAIM_INVALID_ACTION: "Invalid action for this claim's current status.",
  CLAIM_MISSING_INFORMATION: "This claim is missing required information. Please add all required items and attachments before submitting.",
  CLAIM_INVALID_STATUS: "Invalid claim status provided.",
  CLAIM_ITEM_INVALID_STATUS: "Invalid claim item status provided.",
  CLAIM_ITEM_INVALID_TYPE: "Invalid claim item type provided.",
  CLAIM_ATTACHMENT_INVALID_TYPE: "Invalid attachment type. Please upload a valid file.",

  // Cloudinary Errors
  CLOUDINARY_FAILED_INITIALIZE: "Failed to initialize cloud storage service.",
  CLOUDINARY_FAILED_UPLOAD: "Failed to upload file. Please try again.",
  CLOUDINARY_FAILED_DELETE: "Failed to delete file. Please try again.",
  CLOUDINARY_INVALID_URL: "Invalid cloud storage URL.",
  CLOUDINARY_EMPTY_PARAMETER: "Required cloud storage parameter is missing.",
};

export function getErrorMessage(errorCode, defaultMessage = "An error occurred. Please try again.") {
  return ErrorMessages[errorCode] || defaultMessage;
}

export function getErrorMessageFromResponse(error) {
  if (error?.response?.data?.error) {
    return getErrorMessage(error.response.data.error);
  }

  if (error?.response?.data?.message) {
    return error.response.data.message;
  }

  if (error?.message) {
    return error.message;
  }

  return "An unexpected error occurred. Please try again.";
}

export function isAuthError(errorCode) {
  return errorCode?.startsWith('AUTH_') ||
         errorCode === 'REFRESH_TOKEN_NOT_FOUND' ||
         errorCode === 'USER_INACTIVE';
}

export function shouldRedirectToLogin(errorCode) {
  const reloginCodes = [
    'AUTH_INVALID_ACCESS_TOKEN',
    'AUTH_EXPIRED_ACCESS_TOKEN',
    'AUTH_INVALID_REFRESH_TOKEN',
    'AUTH_EXPIRED_REFRESH_TOKEN',
    'AUTH_REVOKED_REFRESH_TOKEN',
    'AUTH_MISSING_USER_ID',
    'REFRESH_TOKEN_NOT_FOUND'
  ];

  return reloginCodes.includes(errorCode);
}

export default ErrorMessages;

