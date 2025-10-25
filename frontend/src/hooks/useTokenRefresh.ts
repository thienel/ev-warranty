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
      const response = await axios.post(
        `${API_BASE_URL}${API_ENDPOINTS.AUTH.TOKEN}`,
        {},
        { withCredentials: true }
      )
      
      const newToken = response.data.data?.access_token
      if (newToken) {
        dispatch(setToken(newToken))
        console.log('Token refreshed successfully')
        return true
      }
      
      return false
    } catch (error) {
      console.error('Token refresh failed:', error)
      dispatch(logout())
      await persistor.purge()
      return false
    }
  }, [dispatch])

  const checkTokenExpiration = useCallback(async () => {
    if (!token || !isAuthenticated) return

    if (isTokenExpired(token)) {
      console.log('Token is expired, attempting refresh...')
      const refreshed = await refreshToken()
      if (!refreshed) {
        console.log('Token refresh failed, logging out')
        window.location.replace('/login')
      }
      return
    }

    const expiration = getTokenExpiration(token)
    if (expiration) {
      const timeUntilExpiration = expiration.getTime() - Date.now()
      
      // If token expires within threshold, refresh it
      if (timeUntilExpiration <= TOKEN_REFRESH_THRESHOLD) {
        console.log('Token expires soon, refreshing...')
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