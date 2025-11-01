import api from './api'
import { API_ENDPOINTS } from '@constants/common-constants'
import type {
  ApiSuccessResponse,
  LoginRequest,
  LoginResponse,
  ValidateTokenResponse,
  RefreshTokenResponse,
} from '@/types'

export const authApi = {
  login: (credentials: LoginRequest): Promise<ApiSuccessResponse<LoginResponse>> => {
    return api.post(API_ENDPOINTS.AUTH.LOGIN, credentials)
  },

  logout: (): Promise<void> => {
    return api.post(API_ENDPOINTS.AUTH.LOGOUT)
  },

  googleLogin: (): void => {
    window.location.href = `${api.defaults.baseURL}${API_ENDPOINTS.AUTH.GOOGLE}`
  },

  validateToken: (): Promise<ApiSuccessResponse<ValidateTokenResponse>> => {
    return api.get(API_ENDPOINTS.AUTH.TOKEN)
  },

  refreshToken: (): Promise<ApiSuccessResponse<RefreshTokenResponse>> => {
    return api.post(API_ENDPOINTS.AUTH.TOKEN)
  },
}

export default authApi
