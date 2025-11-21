export interface ErrorResponse {
  response?: {
    status?: number
    data?: {
      message?: string
      error?: string
    }
  }
  message?: string
}

/**
 * Extract error message from API response
 * Priority: response.data.message > response.data.error > error.message > default message
 */
export function getErrorMessage(
  error: ErrorResponse,
  defaultMessage = 'An unexpected error occurred. Please try again.',
): string {
  // // First priority: message field from API response
  // if (error?.response?.data?.message) {
  //   return error.response.data.message
  // }

  // // Second priority: error field from API response (fallback)
  // if (error?.response?.data?.error) {
  //   return error.response.data.error
  // }

  // // Third priority: error message from the error object
  // if (error?.message) {
  //   return error.message
  // }

  // Default fallback message
  return error.response?.data?.message || defaultMessage
}

/**
 * Check if error is related to authentication
 */
export function isAuthError(error: ErrorResponse): boolean {
  const status = error?.response?.status
  const errorCode = error?.response?.data?.error

  // Check for auth-related status codes
  if (status === 401 || status === 403) {
    return true
  }

  // Check for auth-related error codes
  if (typeof errorCode === 'string' && errorCode.includes('AUTH_')) {
    return true
  }

  return false
}

/**
 * Check if error requires redirect to login page
 */
export function shouldRedirectToLogin(error: ErrorResponse): boolean {
  const status = error?.response?.status
  const errorCode = error?.response?.data?.error

  // Redirect on 401 Unauthorized
  if (status === 401) {
    return true
  }

  // Check for specific auth error codes that require re-login
  const reloginErrorCodes = [
    'AUTH_INVALID_ACCESS_TOKEN',
    'AUTH_EXPIRED_ACCESS_TOKEN',
    'AUTH_INVALID_REFRESH_TOKEN',
    'AUTH_EXPIRED_REFRESH_TOKEN',
    'AUTH_REVOKED_REFRESH_TOKEN',
    'AUTH_MISSING_USER_ID',
    'REFRESH_TOKEN_NOT_FOUND',
  ]

  return typeof errorCode === 'string' && reloginErrorCodes.includes(errorCode)
}
