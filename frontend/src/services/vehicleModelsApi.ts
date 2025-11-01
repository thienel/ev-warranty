import api from './api'
import { API_ENDPOINTS } from '@constants/common-constants'
import type {
  DotNetApiResponse,
  VehicleModel,
  CreateVehicleModelRequest,
  UpdateVehicleModelRequest,
} from '@/types'

interface VehicleModelQueryParams {
  brand?: string
  model?: string
  year?: number
}

export const vehicleModelsApi = {
  getAll: (params?: VehicleModelQueryParams): Promise<DotNetApiResponse<VehicleModel[]>> => {
    const searchParams = new URLSearchParams()
    if (params?.brand) searchParams.append('brand', params.brand)
    if (params?.model) searchParams.append('model', params.model)
    if (params?.year) searchParams.append('year', params.year.toString())

    const query = searchParams.toString()
    return api.get(`${API_ENDPOINTS.VEHICLE_MODELS}${query ? `?${query}` : ''}`)
  },

  getById: (id: string): Promise<DotNetApiResponse<VehicleModel>> => {
    return api.get(`${API_ENDPOINTS.VEHICLE_MODELS}/${id}`)
  },

  create: (data: CreateVehicleModelRequest): Promise<DotNetApiResponse<VehicleModel>> => {
    return api.post(API_ENDPOINTS.VEHICLE_MODELS, data)
  },

  update: (
    id: string,
    data: UpdateVehicleModelRequest,
  ): Promise<DotNetApiResponse<VehicleModel>> => {
    return api.put(`${API_ENDPOINTS.VEHICLE_MODELS}/${id}`, data)
  },

  delete: (id: string): Promise<DotNetApiResponse<VehicleModel>> => {
    return api.delete(`${API_ENDPOINTS.VEHICLE_MODELS}/${id}`)
  },
}

export default vehicleModelsApi
