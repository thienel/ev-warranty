import api from './api'
import { API_ENDPOINTS } from '@constants/common-constants'
import type {
  ApiSuccessResponse,
  CreateClaimRequest,
  UpdateClaimRequest,
  Claim,
  ClaimItem,
  ClaimAttachment,
  ClaimHistory,
  CreateClaimItemRequest,
  ClaimItemListResponse,
  ClaimAttachmentListResponse,
  PaginationParams,
} from '@/types'

// Claims API
export const claimsApi = {
  // Get all claims with optional filtering and pagination
  getAll: (params?: PaginationParams): Promise<ApiSuccessResponse<Claim[]>> => {
    const searchParams = new URLSearchParams()
    if (params?.page) searchParams.append('page', params.page.toString())
    if (params?.limit) searchParams.append('limit', params.limit.toString())
    if (params?.status) searchParams.append('status', params.status)
    
    const query = searchParams.toString()
    return api.get(`${API_ENDPOINTS.CLAIMS}${query ? `?${query}` : ''}`)
  },

  // Get claim by ID
  getById: (id: string): Promise<ApiSuccessResponse<Claim>> => {
    return api.get(`${API_ENDPOINTS.CLAIMS}/${id}`)
  },

  // Create new claim (SC Technician/Staff only)
  create: (data: CreateClaimRequest): Promise<ApiSuccessResponse<Claim>> => {
    return api.post(API_ENDPOINTS.CLAIMS, data)
  },

  // Update claim (SC Staff only)
  update: (id: string, data: UpdateClaimRequest): Promise<void> => {
    return api.put(`${API_ENDPOINTS.CLAIMS}/${id}`, data)
  },

  // Delete claim (SC Staff: hard delete, EVM Staff: soft delete)
  delete: (id: string): Promise<void> => {
    return api.delete(`${API_ENDPOINTS.CLAIMS}/${id}`)
  },

  // Claim Actions
  submit: (id: string): Promise<void> => {
    return api.post(API_ENDPOINTS.CLAIM_ACTIONS.SUBMIT(id))
  },

  cancel: (id: string): Promise<void> => {
    return api.post(API_ENDPOINTS.CLAIM_ACTIONS.CANCEL(id))
  },

  complete: (id: string): Promise<void> => {
    return api.post(API_ENDPOINTS.CLAIM_ACTIONS.COMPLETE(id))
  },

  review: (id: string): Promise<void> => {
    return api.post(API_ENDPOINTS.CLAIM_ACTIONS.REVIEW(id))
  },

  requestInfo: (id: string): Promise<void> => {
    return api.post(API_ENDPOINTS.CLAIM_ACTIONS.REQUEST_INFO(id))
  },

  // Get claim history
  getHistory: (id: string): Promise<ApiSuccessResponse<ClaimHistory[]>> => {
    return api.get(API_ENDPOINTS.CLAIM_ACTIONS.HISTORY(id))
  },
}

// Claim Items API
export const claimItemsApi = {
  // Get all items for a claim
  getByClaimId: (claimId: string): Promise<ApiSuccessResponse<ClaimItemListResponse>> => {
    return api.get(API_ENDPOINTS.CLAIM_ACTIONS.ITEMS(claimId))
  },

  // Get specific claim item
  getById: (claimId: string, itemId: string): Promise<ApiSuccessResponse<ClaimItem>> => {
    return api.get(`${API_ENDPOINTS.CLAIM_ACTIONS.ITEMS(claimId)}/${itemId}`)
  },

  // Create new claim item (SC Staff only)
  create: (claimId: string, data: CreateClaimItemRequest): Promise<ApiSuccessResponse<ClaimItem>> => {
    return api.post(API_ENDPOINTS.CLAIM_ACTIONS.ITEMS(claimId), data)
  },

  // Delete claim item (SC Staff only)
  delete: (claimId: string, itemId: string): Promise<void> => {
    return api.delete(`${API_ENDPOINTS.CLAIM_ACTIONS.ITEMS(claimId)}/${itemId}`)
  },

  // Approve claim item (EVM Staff only)
  approve: (claimId: string, itemId: string): Promise<void> => {
    return api.post(API_ENDPOINTS.CLAIM_ITEM_ACTIONS.APPROVE(claimId, itemId))
  },

  // Reject claim item (EVM Staff only)
  reject: (claimId: string, itemId: string): Promise<void> => {
    return api.post(API_ENDPOINTS.CLAIM_ITEM_ACTIONS.REJECT(claimId, itemId))
  },
}

// Claim Attachments API
export const claimAttachmentsApi = {
  // Get all attachments for a claim
  getByClaimId: (claimId: string): Promise<ApiSuccessResponse<ClaimAttachmentListResponse>> => {
    return api.get(API_ENDPOINTS.CLAIM_ACTIONS.ATTACHMENTS(claimId))
  },

  // Get specific attachment
  getById: (claimId: string, attachmentId: string): Promise<ApiSuccessResponse<ClaimAttachment>> => {
    return api.get(`${API_ENDPOINTS.CLAIM_ACTIONS.ATTACHMENTS(claimId)}/${attachmentId}`)
  },

  // Upload attachments (SC Technician only)
  upload: (claimId: string, files: FileList): Promise<ApiSuccessResponse<ClaimAttachment[]>> => {
    const formData = new FormData()
    Array.from(files).forEach((file) => {
      formData.append('files', file)
    })

    return api.post(API_ENDPOINTS.CLAIM_ACTIONS.ATTACHMENTS(claimId), formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
  },

  // Delete attachment (SC Technician only)
  delete: (claimId: string, attachmentId: string): Promise<void> => {
    return api.delete(`${API_ENDPOINTS.CLAIM_ACTIONS.ATTACHMENTS(claimId)}/${attachmentId}`)
  },
}

export default {
  claims: claimsApi,
  claimItems: claimItemsApi,
  claimAttachments: claimAttachmentsApi,
}