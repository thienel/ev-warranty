import api from './api'
import { API_ENDPOINTS } from '@constants/common-constants'
import type {
  ApiSuccessResponse,
  LoginRequest,
  LoginResponse,
  ValidateTokenResponse,
  RefreshTokenResponse,
} from '@/types'

// Auth API service matching Swagger endpoints
export const authApi = {
  // User login with email and password
  login: (credentials: LoginRequest): Promise<ApiSuccessResponse<LoginResponse>> => {
    return api.post(API_ENDPOINTS.AUTH.LOGIN, credentials)
  },

  // User logout (invalidate refresh token)
  logout: (): Promise<void> => {
    return api.post(API_ENDPOINTS.AUTH.LOGOUT)
  },

  // Initiate Google OAuth login
  googleLogin: (): void => {
    window.location.href = `${api.defaults.baseURL}${API_ENDPOINTS.AUTH.GOOGLE}`
  },

  // Validate access token
  validateToken: (): Promise<ApiSuccessResponse<ValidateTokenResponse>> => {
    return api.get(API_ENDPOINTS.AUTH.TOKEN)
  },

  // Refresh access token using refresh token from cookie
  refreshToken: (): Promise<ApiSuccessResponse<RefreshTokenResponse>> => {
    return api.post(API_ENDPOINTS.AUTH.TOKEN)
  },
}

export default authApi