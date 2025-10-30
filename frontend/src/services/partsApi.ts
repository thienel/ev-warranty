import api from './api'
import { API_ENDPOINTS } from '@constants/common-constants'
import type { ApiSuccessResponse, Part, CreatePartRequest, UpdatePartRequest } from '@/types'

export const partsApi = {
  getAll: (): Promise<ApiSuccessResponse<Part[]>> => {
    return api.get(API_ENDPOINTS.PARTS)
  },

  getById: (id: string): Promise<ApiSuccessResponse<Part>> => {
    return api.get(`${API_ENDPOINTS.PARTS}/${id}`)
  },

  create: (data: CreatePartRequest): Promise<ApiSuccessResponse<Part>> => {
    return api.post(API_ENDPOINTS.PARTS, data)
  },

  update: (id: string, data: UpdatePartRequest): Promise<void> => {
    return api.put(`${API_ENDPOINTS.PARTS}/${id}`, data)
  },

  delete: (id: string): Promise<void> => {
    return api.delete(`${API_ENDPOINTS.PARTS}/${id}`)
  },
}

export default partsApi
