import api from './api'
import { API_ENDPOINTS } from '@constants/common-constants'
import type {
  ApiSuccessResponse,
  User,
  CreateUserRequest,
  UpdateUserRequest,
} from '@/types'

// Users API service matching Swagger endpoints
export const usersApi = {
  // Get all users
  getAll: (): Promise<ApiSuccessResponse<User[]>> => {
    return api.get(API_ENDPOINTS.USERS)
  },

  // Get user by ID
  getById: (id: string): Promise<ApiSuccessResponse<User>> => {
    return api.get(`${API_ENDPOINTS.USERS}/${id}`)
  },

  // Create new user (Admin only)
  create: (data: CreateUserRequest): Promise<ApiSuccessResponse<User>> => {
    return api.post(API_ENDPOINTS.USERS, data)
  },

  // Update user (Admin only)
  update: (id: string, data: UpdateUserRequest): Promise<void> => {
    return api.put(`${API_ENDPOINTS.USERS}/${id}`, data)
  },

  // Delete user (Admin only)
  delete: (id: string): Promise<void> => {
    return api.delete(`${API_ENDPOINTS.USERS}/${id}`)
  },
}

export default usersApi