import api from './api'
import { API_ENDPOINTS } from '@constants/common-constants'
import type { ApiSuccessResponse, Technician } from '@/types'

export const techniciansApi = {
  getAvailable: (): Promise<ApiSuccessResponse<Technician[]>> => {
    return api.get(API_ENDPOINTS.AVAILABLE_TECHNICIANS)
  },
}

export default techniciansApi
