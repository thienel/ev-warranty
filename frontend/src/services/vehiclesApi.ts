import api from './api'
import { API_ENDPOINTS } from '@constants/common-constants'
import type {
  DotNetApiResponse,
  Vehicle,
  CreateVehicleRequest,
  UpdateVehicleRequest,
} from '@/types'

interface VehicleQueryParams {
  vin?: string
  licensePlate?: string
  customerId?: string
  modelId?: string
}

export const vehiclesApi = {
  getAll: (params?: VehicleQueryParams): Promise<DotNetApiResponse<Vehicle[]>> => {
    const searchParams = new URLSearchParams()
    if (params?.vin) searchParams.append('vin', params.vin)
    if (params?.licensePlate) searchParams.append('licensePlate', params.licensePlate)
    if (params?.customerId) searchParams.append('customerId', params.customerId)
    if (params?.modelId) searchParams.append('modelId', params.modelId)

    const query = searchParams.toString()
    return api.get(`${API_ENDPOINTS.VEHICLES}${query ? `?${query}` : ''}`)
  },

  getById: (id: string): Promise<DotNetApiResponse<Vehicle>> => {
    return api.get(`${API_ENDPOINTS.VEHICLES}/${id}`)
  },

  create: (data: CreateVehicleRequest): Promise<DotNetApiResponse<Vehicle>> => {
    return api.post(API_ENDPOINTS.VEHICLES, data)
  },

  update: (id: string, data: UpdateVehicleRequest): Promise<DotNetApiResponse<Vehicle>> => {
    return api.put(`${API_ENDPOINTS.VEHICLES}/${id}`, data)
  },

  delete: (id: string): Promise<DotNetApiResponse<Vehicle>> => {
    return api.delete(`${API_ENDPOINTS.VEHICLES}/${id}`)
  },

  restore: (id: string): Promise<DotNetApiResponse<Vehicle>> => {
    return api.post(`${API_ENDPOINTS.VEHICLES}/${id}/restore`)
  },
}

export default vehiclesApi
