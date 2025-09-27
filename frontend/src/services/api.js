import axios from 'axios'
import { message } from 'antd'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost/api/v1'

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
})

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('accessToken')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('accessToken')
      window.location.href = '/login'
    }
    message.error(error.response?.data?.message || 'Unexpected error occurred')
    return Promise.reject(error)
  }
)

export default api
