import { useEffect, useCallback } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { logout, setToken } from '@/redux/authSlice'
import { getTokenExpiration, isTokenExpired } from '@/utils/auth'
import { API_BASE_URL, API_ENDPOINTS } from '@constants/common-constants'
import { persistor } from '@/redux/store'
import axios from 'axios'
import type { RootState } from '@/redux/store'

const TOKEN_REFRESH_INTERVAL = 5 * 60 * 1000 // 5 minutes
const TOKEN_REFRESH_THRESHOLD = 10 * 60 * 1000 // 10 minutes before expiration

export const useTokenRefresh = () => {
  const dispatch = useDispatch()
  const { token, isAuthenticated } = useSelector((state: RootState) => state.auth)

  const refreshToken = useCallback(async (): Promise<boolean> => {
    try {
      console.log('Attempting to refresh token...')
      const response = await axios.post(
        `${API_BASE_URL}${API_ENDPOINTS.AUTH.TOKEN}`,
        {},
        { withCredentials: true },
      )

      const newToken = response.data.data?.access_token
      if (newToken) {
        dispatch(setToken(newToken))
        console.log('✓ Token refreshed successfully')
        return true
      }

      console.warn('No access token in refresh response')
      return false
    } catch (error: unknown) {
      const err = error as { response?: { status?: number; data?: unknown }; message?: string }
      console.error('✗ Token refresh failed:', err?.response?.data || err?.message)

      // Only logout if it's a real authentication error
      // Don't logout on network errors or temporary failures
      if (err?.response?.status === 401 || err?.response?.status === 404) {
        console.log('Refresh token invalid or not found, logging out')
        dispatch(logout())
        await persistor.purge()
        window.location.replace('/login')
      }

      return false
    }
  }, [dispatch])

  const checkTokenExpiration = useCallback(async () => {
    if (!token || !isAuthenticated) return

    if (isTokenExpired(token)) {
      console.log('⚠ Access token is expired, attempting refresh...')
      const refreshed = await refreshToken()
      if (!refreshed) {
        console.log('✗ Token refresh failed after expiration')
      }
      return
    }

    const expiration = getTokenExpiration(token)
    if (expiration) {
      const timeUntilExpiration = expiration.getTime() - Date.now()
      const minutesUntilExpiration = Math.floor(timeUntilExpiration / 60000)

      // If token expires within threshold, refresh it proactively
      if (timeUntilExpiration <= TOKEN_REFRESH_THRESHOLD && timeUntilExpiration > 0) {
        console.log(
          `⏰ Access token expires in ${minutesUntilExpiration} minutes, refreshing proactively...`,
        )
        await refreshToken()
      }
    }
  }, [token, isAuthenticated, refreshToken])

  // Set up periodic token check
  useEffect(() => {
    if (!isAuthenticated || !token) return

    // Check immediately
    checkTokenExpiration()

    // Set up interval for periodic checks
    const interval = setInterval(checkTokenExpiration, TOKEN_REFRESH_INTERVAL)

    return () => clearInterval(interval)
  }, [isAuthenticated, token, checkTokenExpiration])

  return {
    refreshToken,
    checkTokenExpiration,
  }
}
