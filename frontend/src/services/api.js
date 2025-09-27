import axios from 'axios'
import { message } from 'antd'
import store from '@/redux/store'
import { setToken, logout } from '@/redux/authSlice'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost/api/v1'

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  withCredentials: true,
})

api.interceptors.request.use(
  (config) => {
    const { token } = store.getState().auth
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config
    const { status } = error.response || {}

    if (status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      try {
        const res = await axios.post(`${API_BASE_URL}/auth/refresh`, {}, { withCredentials: true })
        const newToken = res.data.accessToken

        store.dispatch(setToken(newToken))

        originalRequest.headers.Authorization = `Bearer ${newToken}`
        return api(originalRequest)
      } catch (err) {
        store.dispatch(logout())
        window.location.href = '/login'
        return Promise.reject(err)
      }
    }

    message.error(error.response?.data?.message || 'Unexpected error occurred')
    return Promise.reject(error)
  }
)

export default api
