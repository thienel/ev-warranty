import api from './api'
import { API_ENDPOINTS } from '@constants/common-constants'
import type {
  ApiSuccessResponse,
  Customer,
  CreateCustomerRequest,
  UpdateCustomerRequest,
} from '@/types'

interface CustomerQueryParams {
  name?: string
  phone?: string
  email?: string
}

// Customers API service matching Swagger endpoints
export const customersApi = {
  // Get all customers with optional filtering
  getAll: (params?: CustomerQueryParams): Promise<ApiSuccessResponse<Customer[]>> => {
    const searchParams = new URLSearchParams()
    if (params?.name) searchParams.append('name', params.name)
    if (params?.phone) searchParams.append('phone', params.phone)
    if (params?.email) searchParams.append('email', params.email)
    
    const query = searchParams.toString()
    return api.get(`${API_ENDPOINTS.CUSTOMERS}${query ? `?${query}` : ''}`)
  },

  // Get customer by ID
  getById: (id: string): Promise<ApiSuccessResponse<Customer>> => {
    return api.get(`${API_ENDPOINTS.CUSTOMERS}/${id}`)
  },

  // Create new customer
  create: (data: CreateCustomerRequest): Promise<ApiSuccessResponse<Customer>> => {
    return api.post(API_ENDPOINTS.CUSTOMERS, data)
  },

  // Update customer
  update: (id: string, data: UpdateCustomerRequest): Promise<ApiSuccessResponse<Customer>> => {
    return api.put(`${API_ENDPOINTS.CUSTOMERS}/${id}`, data)
  },

  // Delete customer (soft delete)
  delete: (id: string): Promise<ApiSuccessResponse<Customer>> => {
    return api.delete(`${API_ENDPOINTS.CUSTOMERS}/${id}`)
  },

  // Restore deleted customer
  restore: (id: string): Promise<ApiSuccessResponse<Customer>> => {
    return api.post(`${API_ENDPOINTS.CUSTOMERS}/${id}/restore`)
  },
}

export default customersApi