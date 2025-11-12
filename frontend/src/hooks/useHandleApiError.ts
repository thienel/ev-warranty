import { useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import { useDispatch } from 'react-redux'
import { message } from 'antd'
import { logout } from '@/redux/authSlice'
import { persistor } from '@/redux/store'
import {
  getErrorMessage,
  isAuthError,
  shouldRedirectToLogin,
  type ErrorResponse,
} from '@/utils/errorHandler'

interface HandleApiErrorOptions {
  showNotification?: boolean
  duration?: number
  onAuthError?: (error: ErrorResponse) => void
  onError?: (error: ErrorResponse) => void
}

interface ErrorResult {
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
    async (error: ErrorResponse, errorMessage: string, options: HandleApiErrorOptions) => {
      const { showNotification, duration, onAuthError } = options

      if (showNotification) {
        message.error(errorMessage, duration)
      }

      if (shouldRedirectToLogin(error)) {
        dispatch(logout())
        await persistor.purge()

        setTimeout(() => {
          navigate('/login', {
            replace: true,
            state: { from: window.location.pathname },
          })
        }, 1000)
      }

      onAuthError?.(error)
    },
    [navigate, dispatch],
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

      const errorMessage = getErrorMessage(error)
      const errorStatus = error?.response?.status

      console.error('API Error:', {
        message: errorMessage,
        status: errorStatus,
        data: error?.response?.data,
      })

      const isAuth = isAuthError(error)
      const shouldRedirect = shouldRedirectToLogin(error)

      if (isAuth) {
        await handleAuthError(error, errorMessage, mergedOptions)
      } else {
        handleGeneralError(errorMessage, mergedOptions)
      }

      onError?.(error)

      return {
        message: errorMessage,
        status: errorStatus,
        isAuthError: isAuth,
        shouldRedirect,
      }
    },
    [handleAuthError, handleGeneralError],
  )
}

export default useHandleApiError
