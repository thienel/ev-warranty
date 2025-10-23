import axios, { type AxiosRequestConfig, type AxiosResponse, AxiosError } from 'axios'
import store, { persistor } from '@/redux/store'
import { setToken, logout } from '@/redux/authSlice'
import { API_BASE_URL, API_ENDPOINTS } from '@constants/common-constants'

interface ExtendedAxiosRequestConfig extends AxiosRequestConfig {
  _retry?: boolean
}

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  withCredentials: true,
})

api.interceptors.request.use(
  (config) => {
    const { token } = store.getState().auth
    if (token) {
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
      if (originalRequest) {
        originalRequest._retry = true
      }
      try {
        const res = await axios.post(
          `${API_BASE_URL}${API_ENDPOINTS.AUTH.TOKEN}`,
          {},
          { withCredentials: true }
        )
        const newToken = res.data.data.access_token

        store.dispatch(setToken(newToken))

        if (originalRequest?.headers) {
          originalRequest.headers.Authorization = `Bearer ${newToken}`
        }
        return api(originalRequest as AxiosRequestConfig)
      } catch (err) {
        store.dispatch(logout())
        await persistor.purge()
        return Promise.reject(err)
      }
    }

    console.log(error.response?.data || 'Unexpected error occurred')
    return Promise.reject(error)
  }
)

export default api