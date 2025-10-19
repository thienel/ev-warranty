const ErrorMessages = {
  // Common Errors
  COMMON_INTERNAL_ERROR: "Something went wrong on our end. Please try again in a moment.",

  // Database Errors
  DB_OPERATION_ERROR: "We're having trouble saving your changes. Please try again.",
  DB_DUPLICATE_KEY: "This information already exists in our system.",
  DB_HASH_PASSWORD_ERROR: "Unable to process your password. Please try again.",

  // API Request Errors
  API_INVALID_JSON_REQUEST: "We couldn't process your request. Please check your information and try again.",
  API_INVALID_QUERY_PARAMETER: "Some information provided is invalid. Please check and try again.",
  API_INVALID_MULTIPART_FORM_REQUEST: "There's an issue with your file upload. Please check your files and try again.",
  API_INVALID_UUID: "The provided ID is invalid. Please try again.",

  // Authentication Errors - Session Related
  AUTH_INVALID_ACCESS_TOKEN: "Your session is invalid. Please sign in again.",
  AUTH_EXPIRED_ACCESS_TOKEN: "Your session has expired. Please sign in again.",
  AUTH_FAILED_HASH_TOKEN: "Unable to authenticate. Please try signing in again.",
  AUTH_INVALID_REFRESH_TOKEN: "Your session is no longer valid. Please sign in again.",
  AUTH_EXPIRED_REFRESH_TOKEN: "Your session has expired. Please sign in again.",
  AUTH_REVOKED_REFRESH_TOKEN: "Your session has been ended. Please sign in again.",
  AUTH_INVALID_AUTH_HEADER: "Authentication failed. Please sign in again.",

  // Authentication Errors - Credentials
  AUTH_INVALID_CREDENTIALS: "The email or password you entered is incorrect. Please try again.",
  AUTH_FAILED_SIGN_ACCESS_TOKEN: "We're having trouble signing you in. Please try again.",
  AUTH_FAILED_GENERATE_REFRESH_TOKEN: "Unable to maintain your session. Please sign in again.",
  AUTH_UNEXPECTED_SIGNING_METHOD: "An authentication error occurred. Please contact support.",
  AUTH_MISSING_USER_ID: "Your account information is missing. Please sign in again.",
  AUTH_INVALID_USER_ID: "Your account ID is invalid. Please sign in again.",

  // Refresh Token Errors
  REFRESH_TOKEN_NOT_FOUND: "Your session could not be found. Please sign in again.",

  // Office Errors
  OFFICE_INVALID_TYPE: "The office type you selected is invalid.",
  OFFICE_NOT_FOUND: "We couldn't find that office.",

  // User Errors
  USER_NOT_FOUND: "We couldn't find your account. Please check your information.",
  USER_INACTIVE: "Your account is currently inactive. Please contact support for help.",
  USER_PASSWORD_INVALID: "The password you entered is incorrect. Please try again.",
  USER_INVALID_INPUT: "Some of your information is invalid. Please check and try again.",

  // Claim Errors - Not Found
  CLAIM_ITEM_NOT_FOUND: "We couldn't find that claim item.",
  CLAIM_HISTORY_NOT_FOUND: "No history found for this claim.",
  CLAIM_ATTACHMENT_NOT_FOUND: "The attachment you're looking for couldn't be found.",
  CLAIM_NOT_FOUND: "We couldn't find that claim.",

  // Claim Status Errors
  CLAIM_STATUS_NOT_ALLOWED_TO_UPDATE: "This claim can't be updated at this stage.",
  CLAIM_STATUS_NOT_ALLOWED_TO_DELETE: "This claim can't be deleted at this stage.",
  CLAIM_INVALID_ACTION: "That action isn't available for this claim right now.",
  CLAIM_MISSING_INFORMATION: "Please add all required items and attachments before submitting your claim.",
  CLAIM_INVALID_STATUS: "The status you selected is not valid for this claim.",
  CLAIM_ITEM_INVALID_STATUS: "That status isn't valid for this claim item.",
  CLAIM_ITEM_INVALID_TYPE: "The item type you selected is invalid.",
  CLAIM_ATTACHMENT_INVALID_TYPE: "That file type isn't supported. Please upload a valid file.",

  // Cloudinary Errors
  CLOUDINARY_FAILED_INITIALIZE: "We're having trouble with our file storage. Please try again later.",
  CLOUDINARY_FAILED_UPLOAD: "Your file couldn't be uploaded. Please try again.",
  CLOUDINARY_FAILED_DELETE: "We couldn't delete that file. Please try again.",
  CLOUDINARY_INVALID_URL: "The file URL is invalid.",
  CLOUDINARY_EMPTY_PARAMETER: "Required file information is missing.",
};

export function getErrorMessage(errorCode, defaultMessage = "Something went wrong. Please try again.") {
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