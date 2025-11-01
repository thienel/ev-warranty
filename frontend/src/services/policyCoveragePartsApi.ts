import api from './api'
import { API_ENDPOINTS } from '@constants/common-constants'
import type {
  ApiSuccessResponse,
  PolicyCoveragePart,
  CreatePolicyCoveragePartRequest,
  UpdatePolicyCoveragePartRequest,
} from '@/types'

export const policyCoveragePartsApi = {
  getAll: (policyId?: string): Promise<ApiSuccessResponse<PolicyCoveragePart[]>> => {
    const params = policyId ? { policyId } : {}
    return api.get(API_ENDPOINTS.POLICY_COVERAGE_PARTS, { params })
  },

  getById: (id: string): Promise<ApiSuccessResponse<PolicyCoveragePart>> => {
    return api.get(`${API_ENDPOINTS.POLICY_COVERAGE_PARTS}/${id}`)
  },

  create: (
    data: CreatePolicyCoveragePartRequest,
  ): Promise<ApiSuccessResponse<PolicyCoveragePart>> => {
    return api.post(API_ENDPOINTS.POLICY_COVERAGE_PARTS, data)
  },

  update: (id: string, data: UpdatePolicyCoveragePartRequest): Promise<void> => {
    return api.put(`${API_ENDPOINTS.POLICY_COVERAGE_PARTS}/${id}`, data)
  },

  delete: (id: string): Promise<void> => {
    return api.delete(`${API_ENDPOINTS.POLICY_COVERAGE_PARTS}/${id}`)
  },
}

export default policyCoveragePartsApi
