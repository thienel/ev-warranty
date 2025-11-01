import api from './api'
import { API_ENDPOINTS } from '@constants/common-constants'
import type {
  DotNetApiResponse,
  Customer,
  CreateCustomerRequest,
  UpdateCustomerRequest,
} from '@/types'

interface CustomerQueryParams {
  name?: string
  phone?: string
  email?: string
}

export const customersApi = {
  getAll: (params?: CustomerQueryParams): Promise<DotNetApiResponse<Customer[]>> => {
    const searchParams = new URLSearchParams()
    if (params?.name) searchParams.append('name', params.name)
    if (params?.phone) searchParams.append('phone', params.phone)
    if (params?.email) searchParams.append('email', params.email)

    const query = searchParams.toString()
    return api.get(`${API_ENDPOINTS.CUSTOMERS}${query ? `?${query}` : ''}`)
  },

  getById: (id: string): Promise<DotNetApiResponse<Customer>> => {
    return api.get(`${API_ENDPOINTS.CUSTOMERS}/${id}`)
  },

  create: (data: CreateCustomerRequest): Promise<DotNetApiResponse<Customer>> => {
    return api.post(API_ENDPOINTS.CUSTOMERS, data)
  },

  update: (id: string, data: UpdateCustomerRequest): Promise<DotNetApiResponse<Customer>> => {
    return api.put(`${API_ENDPOINTS.CUSTOMERS}/${id}`, data)
  },

  delete: (id: string): Promise<DotNetApiResponse<Customer>> => {
    return api.delete(`${API_ENDPOINTS.CUSTOMERS}/${id}`)
  },

  restore: (id: string): Promise<DotNetApiResponse<Customer>> => {
    return api.post(`${API_ENDPOINTS.CUSTOMERS}/${id}/restore`)
  },
}

export default customersApi
