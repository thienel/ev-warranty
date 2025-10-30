import api from './api'
import { API_ENDPOINTS } from '@constants/common-constants'
import type {
  ApiSuccessResponse,
  WarrantyPolicy,
  CreateWarrantyPolicyRequest,
  UpdateWarrantyPolicyRequest,
} from '@/types'

export const warrantyPoliciesApi = {
  getAll: (): Promise<ApiSuccessResponse<WarrantyPolicy[]>> => {
    return api.get(API_ENDPOINTS.WARRANTY_POLICIES)
  },

  getById: (id: string): Promise<ApiSuccessResponse<WarrantyPolicy>> => {
    return api.get(`${API_ENDPOINTS.WARRANTY_POLICIES}/${id}`)
  },

  create: (data: CreateWarrantyPolicyRequest): Promise<ApiSuccessResponse<WarrantyPolicy>> => {
    return api.post(API_ENDPOINTS.WARRANTY_POLICIES, data)
  },

  update: (id: string, data: UpdateWarrantyPolicyRequest): Promise<void> => {
    return api.put(`${API_ENDPOINTS.WARRANTY_POLICIES}/${id}`, data)
  },

  delete: (id: string): Promise<void> => {
    return api.delete(`${API_ENDPOINTS.WARRANTY_POLICIES}/${id}`)
  },

  getDetails: (id: string): Promise<ApiSuccessResponse<WarrantyPolicy>> => {
    return api.get(`${API_ENDPOINTS.WARRANTY_POLICIES}/${id}/details`)
  },
}

export default warrantyPoliciesApi
