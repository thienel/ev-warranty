import api from './api'
import { API_ENDPOINTS } from '@constants/common-constants'
import type { ApiSuccessResponse, User, CreateUserRequest, UpdateUserRequest } from '@/types'

export const usersApi = {
  getAll: (): Promise<ApiSuccessResponse<User[]>> => {
    return api.get(API_ENDPOINTS.USERS)
  },

  getById: (id: string): Promise<ApiSuccessResponse<User>> => {
    return api.get(`${API_ENDPOINTS.USERS}/${id}`)
  },

  create: (data: CreateUserRequest): Promise<ApiSuccessResponse<User>> => {
    return api.post(API_ENDPOINTS.USERS, data)
  },

  update: (id: string, data: UpdateUserRequest): Promise<void> => {
    return api.put(`${API_ENDPOINTS.USERS}/${id}`, data)
  },

  delete: (id: string): Promise<void> => {
    return api.delete(`${API_ENDPOINTS.USERS}/${id}`)
  },
}

export default usersApi
