import api from './api'
import { API_ENDPOINTS } from '@constants/common-constants'
import type {
  ApiSuccessResponse,
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

// Vehicles API service matching Swagger endpoints
export const vehiclesApi = {
  // Get all vehicles with optional filtering
  getAll: (params?: VehicleQueryParams): Promise<ApiSuccessResponse<Vehicle[]>> => {
    const searchParams = new URLSearchParams()
    if (params?.vin) searchParams.append('vin', params.vin)
    if (params?.licensePlate) searchParams.append('licensePlate', params.licensePlate)
    if (params?.customerId) searchParams.append('customerId', params.customerId)
    if (params?.modelId) searchParams.append('modelId', params.modelId)
    
    const query = searchParams.toString()
    return api.get(`${API_ENDPOINTS.VEHICLES}${query ? `?${query}` : ''}`)
  },

  // Get vehicle by ID
  getById: (id: string): Promise<ApiSuccessResponse<Vehicle>> => {
    return api.get(`${API_ENDPOINTS.VEHICLES}/${id}`)
  },

  // Create new vehicle
  create: (data: CreateVehicleRequest): Promise<ApiSuccessResponse<Vehicle>> => {
    return api.post(API_ENDPOINTS.VEHICLES, data)
  },

  // Update vehicle
  update: (id: string, data: UpdateVehicleRequest): Promise<ApiSuccessResponse<Vehicle>> => {
    return api.put(`${API_ENDPOINTS.VEHICLES}/${id}`, data)
  },

  // Delete vehicle (soft delete)
  delete: (id: string): Promise<ApiSuccessResponse<Vehicle>> => {
    return api.delete(`${API_ENDPOINTS.VEHICLES}/${id}`)
  },

  // Restore deleted vehicle
  restore: (id: string): Promise<ApiSuccessResponse<Vehicle>> => {
    return api.post(`${API_ENDPOINTS.VEHICLES}/${id}/restore`)
  },
}

export default vehiclesApi