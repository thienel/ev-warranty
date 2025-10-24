import axios, { type AxiosRequestConfig, type AxiosResponse, AxiosError } from 'axios'
import store, { persistor } from '@/redux/store'
import { setToken, logout } from '@/redux/authSlice'
import { API_BASE_URL, API_ENDPOINTS } from '@constants/common-constants'
import { isTokenExpired, isTokenValid } from '@/utils/auth'

interface ExtendedAxiosRequestConfig extends AxiosRequestConfig {
  _retry?: boolean
}

// Keep track of token refresh to prevent multiple concurrent refresh attempts
let isRefreshing = false
let failedQueue: Array<{
  resolve: (value: unknown) => void
  reject: (reason?: unknown) => void
}> = []

const processQueue = (error: unknown, token: string | null = null) => {
  failedQueue.forEach(({ resolve, reject }) => {
    if (error) {
      reject(error)
    } else {
      resolve(token)
    }
  })
  
  failedQueue = []
}

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use(
  (config) => {
    const { token } = store.getState().auth
    if (token) {
      // Check if token is valid and not expired before making request
      if (!isTokenValid(token) || isTokenExpired(token)) {
        console.warn('Invalid or expired token detected, clearing auth state')
        store.dispatch(logout())
        return Promise.reject(new Error('Invalid or expired token'))
      }
      
      config.headers = config.headers || {}
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error: AxiosError) => Promise.reject(error)
)

api.interceptors.response.use(
  (response: AxiosResponse) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as ExtendedAxiosRequestConfig
    const { status } = error.response || {}

    if (status === 401 && !originalRequest?._retry) {
      if (isRefreshing) {
        // If token is being refreshed, queue this request
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        }).then(token => {
          if (originalRequest?.headers) {
            originalRequest.headers.Authorization = `Bearer ${token}`
          }
          return api(originalRequest as AxiosRequestConfig)
        }).catch(err => {
          return Promise.reject(err)
        })
      }

      if (originalRequest) {
        originalRequest._retry = true
      }

      isRefreshing = true

      try {
        const res = await axios.post(
          `${API_BASE_URL}${API_ENDPOINTS.AUTH.TOKEN}`,
          {},
          { withCredentials: true }
        )
        const newToken = res.data.data.access_token

        store.dispatch(setToken(newToken))
        processQueue(null, newToken)

        if (originalRequest?.headers) {
          originalRequest.headers.Authorization = `Bearer ${newToken}`
        }
        
        return api(originalRequest as AxiosRequestConfig)
      } catch (refreshError) {
        processQueue(refreshError, null)
        store.dispatch(logout())
        await persistor.purge()
        
        // Redirect to login page
        if (window.location.pathname !== '/login') {
          window.location.replace('/login')
        }
        
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }

    return Promise.reject(error)
  }
)

export default api