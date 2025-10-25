import api from './api'
import { API_ENDPOINTS } from '@constants/common-constants'
import type {
  ApiSuccessResponse,
  Office,
  CreateOfficeRequest,
  UpdateOfficeRequest,
} from '@/types'

// Offices API service matching Swagger endpoints
export const officesApi = {
  // Get all offices
  getAll: (): Promise<ApiSuccessResponse<Office[]>> => {
    return api.get(API_ENDPOINTS.OFFICES)
  },

  // Get office by ID
  getById: (id: string): Promise<ApiSuccessResponse<Office>> => {
    return api.get(`${API_ENDPOINTS.OFFICES}/${id}`)
  },

  // Create new office (Admin only)
  create: (data: CreateOfficeRequest): Promise<ApiSuccessResponse<Office>> => {
    return api.post(API_ENDPOINTS.OFFICES, data)
  },

  // Update office (Admin only)
  update: (id: string, data: UpdateOfficeRequest): Promise<void> => {
    return api.put(`${API_ENDPOINTS.OFFICES}/${id}`, data)
  },

  // Delete office (Admin only)
  delete: (id: string): Promise<void> => {
    return api.delete(`${API_ENDPOINTS.OFFICES}/${id}`)
  },
}

export default officesApi