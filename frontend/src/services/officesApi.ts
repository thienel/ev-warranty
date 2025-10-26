import api from './api'
import { API_ENDPOINTS } from '@constants/common-constants'
import type { ApiSuccessResponse, Office, CreateOfficeRequest, UpdateOfficeRequest } from '@/types'

export const officesApi = {
  getAll: (): Promise<ApiSuccessResponse<Office[]>> => {
    return api.get(API_ENDPOINTS.OFFICES)
  },

  getById: (id: string): Promise<ApiSuccessResponse<Office>> => {
    return api.get(`${API_ENDPOINTS.OFFICES}/${id}`)
  },

  create: (data: CreateOfficeRequest): Promise<ApiSuccessResponse<Office>> => {
    return api.post(API_ENDPOINTS.OFFICES, data)
  },

  update: (id: string, data: UpdateOfficeRequest): Promise<void> => {
    return api.put(`${API_ENDPOINTS.OFFICES}/${id}`, data)
  },

  delete: (id: string): Promise<void> => {
    return api.delete(`${API_ENDPOINTS.OFFICES}/${id}`)
  },
}

export default officesApi
