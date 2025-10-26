import api from './api'
import { API_ENDPOINTS } from '@constants/common-constants'
import type {
  ApiSuccessResponse,
  VehicleModel,
  CreateVehicleModelRequest,
  UpdateVehicleModelRequest,
} from '@/types'

interface VehicleModelQueryParams {
  brand?: string
  model?: string
  year?: number
}

// Vehicle Models API service matching Swagger endpoints
export const vehicleModelsApi = {
  // Get all vehicle models with optional filtering
  getAll: (params?: VehicleModelQueryParams): Promise<ApiSuccessResponse<VehicleModel[]>> => {
    const searchParams = new URLSearchParams()
    if (params?.brand) searchParams.append('brand', params.brand)
    if (params?.model) searchParams.append('model', params.model)
    if (params?.year) searchParams.append('year', params.year.toString())
    
    const query = searchParams.toString()
    return api.get(`${API_ENDPOINTS.VEHICLE_MODELS}${query ? `?${query}` : ''}`)
  },

  // Get vehicle model by ID
  getById: (id: string): Promise<ApiSuccessResponse<VehicleModel>> => {
    return api.get(`${API_ENDPOINTS.VEHICLE_MODELS}/${id}`)
  },

  // Create new vehicle model
  create: (data: CreateVehicleModelRequest): Promise<ApiSuccessResponse<VehicleModel>> => {
    return api.post(API_ENDPOINTS.VEHICLE_MODELS, data)
  },

  // Update vehicle model
  update: (id: string, data: UpdateVehicleModelRequest): Promise<ApiSuccessResponse<VehicleModel>> => {
    return api.put(`${API_ENDPOINTS.VEHICLE_MODELS}/${id}`, data)
  },

  // Delete vehicle model
  delete: (id: string): Promise<ApiSuccessResponse<VehicleModel>> => {
    return api.delete(`${API_ENDPOINTS.VEHICLE_MODELS}/${id}`)
  },
}

export default vehicleModelsApi