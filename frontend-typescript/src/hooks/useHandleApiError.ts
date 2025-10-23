import { useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import { useDispatch } from 'react-redux'
import { message } from 'antd'
import { logout } from '@/redux/authSlice'
import { persistor } from '@/redux/store'
import {
  getErrorMessageFromResponse,
  isAuthError,
  shouldRedirectToLogin,
  type ErrorResponse,
  type ErrorCode,
} from '@/constants/error-messages'

interface HandleApiErrorOptions {
  showNotification?: boolean
  duration?: number
  onAuthError?: (error: ErrorResponse, errorCode: ErrorCode) => void
  onError?: (error: ErrorResponse, errorCode?: ErrorCode) => void
}

interface ErrorResult {
  errorCode?: ErrorCode
  message: string
  status?: number
  isAuthError: boolean
  shouldRedirect: boolean
}

const DEFAULT_OPTIONS: HandleApiErrorOptions = {
  showNotification: true,
  duration: 3,
}

message.config({
  top: 10,
  maxCount: 3,
})

const useHandleApiError = () => {
  const navigate = useNavigate()
  const dispatch = useDispatch()

  const handleAuthError = useCallback(
    async (error: ErrorResponse, errorCode: ErrorCode, errorMessage: string, options: HandleApiErrorOptions) => {
      const { showNotification, duration, onAuthError } = options

      if (showNotification) {
        message.error(errorMessage, duration)
      }

      if (shouldRedirectToLogin(errorCode)) {
        dispatch(logout())
        await persistor.purge()

        setTimeout(() => {
          navigate('/login', {
            replace: true,
            state: { from: window.location.pathname },
          })
        }, 1000)
      }

      onAuthError?.(error, errorCode)
    },
    [navigate, dispatch]
  )

  const handleGeneralError = useCallback((errorMessage: string, options: HandleApiErrorOptions) => {
    const { showNotification, duration } = options

    if (showNotification) {
      message.error(errorMessage, duration)
    }
  }, [])

  return useCallback(
    async (error: ErrorResponse, options: HandleApiErrorOptions = {}): Promise<ErrorResult> => {
      const mergedOptions = { ...DEFAULT_OPTIONS, ...options }
      const { onError } = mergedOptions

      const errorCode = error?.response?.data?.error
      const errorMessage = getErrorMessageFromResponse(error)
      const errorStatus = error?.response?.status

      console.error('API Error:', {
        errorCode,
        message: errorMessage,
        status: errorStatus,
        data: error?.response?.data,
      })

      if (errorCode && isAuthError(errorCode)) {
        await handleAuthError(error, errorCode, errorMessage, mergedOptions)
      } else {
        handleGeneralError(errorMessage, mergedOptions)
      }

      onError?.(error, errorCode)

      return {
        errorCode,
        message: errorMessage,
        status: errorStatus,
        isAuthError: errorCode ? isAuthError(errorCode) : false,
        shouldRedirect: errorCode ? shouldRedirectToLogin(errorCode) : false,
      }
    },
    [handleAuthError, handleGeneralError]
  )
}

export default useHandleApiError