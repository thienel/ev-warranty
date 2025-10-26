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
  ClaimListResponse,
} from '@/types'

export const claimsApi = {
  getAll: (params?: PaginationParams): Promise<ApiSuccessResponse<ClaimListResponse>> => {
    const searchParams = new URLSearchParams()
    if (params?.page) searchParams.append('page', params.page.toString())
    if (params?.limit) searchParams.append('limit', params.limit.toString())
    if (params?.status) searchParams.append('status', params.status)

    const query = searchParams.toString()
    return api
      .get(`${API_ENDPOINTS.CLAIMS}${query ? `?${query}` : ''}`)
      .then((response) => response.data)
  },

  getById: (id: string): Promise<ApiSuccessResponse<Claim>> => {
    return api.get(`${API_ENDPOINTS.CLAIMS}/${id}`)
  },

  create: (data: CreateClaimRequest): Promise<ApiSuccessResponse<Claim>> => {
    return api.post(API_ENDPOINTS.CLAIMS, data)
  },

  update: (id: string, data: UpdateClaimRequest): Promise<void> => {
    return api.put(`${API_ENDPOINTS.CLAIMS}/${id}`, data)
  },

  delete: (id: string): Promise<void> => {
    return api.delete(`${API_ENDPOINTS.CLAIMS}/${id}`)
  },

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

  getHistory: (id: string): Promise<ApiSuccessResponse<ClaimHistory[]>> => {
    return api.get(API_ENDPOINTS.CLAIM_ACTIONS.HISTORY(id))
  },
}

export const claimItemsApi = {
  getByClaimId: (claimId: string): Promise<ApiSuccessResponse<ClaimItemListResponse>> => {
    return api.get(API_ENDPOINTS.CLAIM_ACTIONS.ITEMS(claimId))
  },

  getById: (claimId: string, itemId: string): Promise<ApiSuccessResponse<ClaimItem>> => {
    return api.get(`${API_ENDPOINTS.CLAIM_ACTIONS.ITEMS(claimId)}/${itemId}`)
  },

  create: (
    claimId: string,
    data: CreateClaimItemRequest,
  ): Promise<ApiSuccessResponse<ClaimItem>> => {
    return api.post(API_ENDPOINTS.CLAIM_ACTIONS.ITEMS(claimId), data)
  },

  delete: (claimId: string, itemId: string): Promise<void> => {
    return api.delete(`${API_ENDPOINTS.CLAIM_ACTIONS.ITEMS(claimId)}/${itemId}`)
  },

  approve: (claimId: string, itemId: string): Promise<void> => {
    return api.post(API_ENDPOINTS.CLAIM_ITEM_ACTIONS.APPROVE(claimId, itemId))
  },

  reject: (claimId: string, itemId: string): Promise<void> => {
    return api.post(API_ENDPOINTS.CLAIM_ITEM_ACTIONS.REJECT(claimId, itemId))
  },
}

export const claimAttachmentsApi = {
  getByClaimId: (claimId: string): Promise<ApiSuccessResponse<ClaimAttachmentListResponse>> => {
    return api.get(API_ENDPOINTS.CLAIM_ACTIONS.ATTACHMENTS(claimId))
  },

  getById: (
    claimId: string,
    attachmentId: string,
  ): Promise<ApiSuccessResponse<ClaimAttachment>> => {
    return api.get(`${API_ENDPOINTS.CLAIM_ACTIONS.ATTACHMENTS(claimId)}/${attachmentId}`)
  },

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

  delete: (claimId: string, attachmentId: string): Promise<void> => {
    return api.delete(`${API_ENDPOINTS.CLAIM_ACTIONS.ATTACHMENTS(claimId)}/${attachmentId}`)
  },
}

export default {
  claims: claimsApi,
  claimItems: claimItemsApi,
  claimAttachments: claimAttachmentsApi,
}
