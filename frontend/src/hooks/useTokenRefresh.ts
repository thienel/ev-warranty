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
        { withCredentials: true },
      )

      const newToken = response.data.data?.token
      if (newToken) {
        dispatch(setToken(newToken))
        return true
      }

      return false
    } catch (error: unknown) {
      const err = error as { response?: { status?: number; data?: unknown }; message?: string }

      if (err?.response?.status === 401 || err?.response?.status === 404) {
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
      await refreshToken()
      return
    }

    const expiration = getTokenExpiration(token)
    if (expiration) {
      const timeUntilExpiration = expiration.getTime() - Date.now()

      if (timeUntilExpiration <= TOKEN_REFRESH_THRESHOLD && timeUntilExpiration > 0) {
        await refreshToken()
      }
    }
  }, [token, isAuthenticated, refreshToken])

  useEffect(() => {
    if (!isAuthenticated || !token) return

    checkTokenExpiration()

    const interval = setInterval(checkTokenExpiration, TOKEN_REFRESH_INTERVAL)

    return () => clearInterval(interval)
  }, [isAuthenticated, token, checkTokenExpiration])

  return {
    refreshToken,
    checkTokenExpiration,
  }
}
