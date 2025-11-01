import api from './api'
import { API_ENDPOINTS } from '@constants/common-constants'
import type {
  ApiSuccessResponse,
  PartCategory,
  CreatePartCategoryRequest,
  UpdatePartCategoryRequest,
} from '@/types'

export const partCategoriesApi = {
  getAll: (): Promise<ApiSuccessResponse<PartCategory[]>> => {
    return api.get(API_ENDPOINTS.PART_CATEGORIES)
  },

  getById: (id: string): Promise<ApiSuccessResponse<PartCategory>> => {
    return api.get(`${API_ENDPOINTS.PART_CATEGORIES}/${id}`)
  },

  create: (data: CreatePartCategoryRequest): Promise<ApiSuccessResponse<PartCategory>> => {
    return api.post(API_ENDPOINTS.PART_CATEGORIES, data)
  },

  update: (id: string, data: UpdatePartCategoryRequest): Promise<void> => {
    return api.put(`${API_ENDPOINTS.PART_CATEGORIES}/${id}`, data)
  },

  delete: (id: string): Promise<void> => {
    return api.delete(`${API_ENDPOINTS.PART_CATEGORIES}/${id}`)
  },

  getHierarchy: (): Promise<ApiSuccessResponse<PartCategory[]>> => {
    return api.get(`${API_ENDPOINTS.PART_CATEGORIES}/hierarchy`)
  },

  getHierarchyById: (id: string): Promise<ApiSuccessResponse<PartCategory>> => {
    return api.get(`${API_ENDPOINTS.PART_CATEGORIES}/${id}/hierarchy`)
  },
}

export default partCategoriesApi
