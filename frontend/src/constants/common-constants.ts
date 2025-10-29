import type { Rule } from 'antd/es/form'

export const PASSWORD_RULES: Rule[] = [
  {
    required: true,
    message: 'Please input your password!',
  },
  {
    min: 8,
    message: 'Password must be at least 8 characters long!',
  },
  {
    pattern: /[a-z]/,
    message: 'Password must contain at least one lowercase letter!',
  },
  {
    pattern: /[A-Z]/,
    message: 'Password must contain at least one uppercase letter!',
  },
  {
    pattern: /\d/,
    message: 'Password must contain at least one digit!',
  },
  {
    pattern: /[^A-Za-z0-9]/,
    message: 'Password must contain at least one special character!',
  },
]

export const EMAIL_RULES: Rule[] = [
  {
    required: true,
    message: 'Please input your email!',
  },
  {
    type: 'email',
    message: 'Please enter a valid email address!',
  },
]

export const USER_ROLES = {
  ADMIN: 'ADMIN',
  SC_STAFF: 'SC_STAFF',
  SC_TECHNICIAN: 'SC_TECHNICIAN',
  EVM_STAFF: 'EVM_STAFF',
} as const

export type UserRole = (typeof USER_ROLES)[keyof typeof USER_ROLES]

export const ROLE_LABELS: Record<UserRole, string> = {
  [USER_ROLES.ADMIN]: 'Admin',
  [USER_ROLES.SC_STAFF]: 'SC Staff',
  [USER_ROLES.SC_TECHNICIAN]: 'SC Technician',
  [USER_ROLES.EVM_STAFF]: 'EVM Staff',
}

export const ERROR_MESSAGES: Record<number, string> = {
  403: 'Sorry, you are not authorized to access this page.',
  500: 'Sorry, something went wrong on the server.',
  404: 'Sorry, the page you visited does not exist.',
}

export const API_BASE_URL: string = import.meta.env.VITE_API_URL || 'https://localhost'

export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: '/api/v1/auth/login',
    LOGOUT: '/api/v1/auth/logout',
    GOOGLE: '/api/v1/auth/google',
    GOOGLE_CALLBACK: '/api/v1/auth/google/callback',
    TOKEN: '/api/v1/auth/token',
  },
  USERS: '/api/v1/users',
  OFFICES: '/api/v1/offices',
  CLAIMS: '/api/v1/claims',
  CUSTOMERS: '/api/v1/customers',
  VEHICLES: '/api/v1/vehicles',
  VEHICLE_MODELS: '/api/v1/vehicle-models',

  CLAIM_ACTIONS: {
    SUBMIT: (id: string) => `/api/v1/claims/${id}/submit`,
    CANCEL: (id: string) => `/api/v1/claims/${id}/cancel`,
    COMPLETE: (id: string) => `/api/v1/claims/${id}/complete`,
    REVIEW: (id: string) => `/api/v1/claims/${id}/review`,
    REQUEST_INFO: (id: string) => `/api/v1/claims/${id}/request-information`,
    HISTORY: (id: string) => `/api/v1/claims/${id}/history`,
    ITEMS: (id: string) => `/api/v1/claims/${id}/items`,
    ATTACHMENTS: (id: string) => `/api/v1/claims/${id}/attachments`,
  },

  CLAIM_ITEM_ACTIONS: {
    APPROVE: (claimId: string, itemId: string) =>
      `/api/v1/claims/${claimId}/items/${itemId}/approve`,
    REJECT: (claimId: string, itemId: string) => `/api/v1/claims/${claimId}/items/${itemId}/reject`,
  },
} as const

export const CLAIM_STATUSES = {
  DRAFT: 'DRAFT',
  SUBMITTED: 'SUBMITTED',
  REVIEWING: 'REVIEWING',
  REQUEST_INFO: 'REQUEST_INFO',
  APPROVED: 'APPROVED',
  PARTIALLY_APPROVED: 'PARTIALLY_APPROVED',
  REJECTED: 'REJECTED',
  CANCELLED: 'CANCELLED',
  COMPLETED: 'COMPLETED',
} as const

export type ClaimStatus = (typeof CLAIM_STATUSES)[keyof typeof CLAIM_STATUSES]

export const CLAIM_STATUS_LABELS: Record<ClaimStatus, string> = {
  [CLAIM_STATUSES.DRAFT]: 'Draft',
  [CLAIM_STATUSES.SUBMITTED]: 'Submitted',
  [CLAIM_STATUSES.REVIEWING]: 'Reviewing',
  [CLAIM_STATUSES.REQUEST_INFO]: 'Request Info',
  [CLAIM_STATUSES.APPROVED]: 'Approved',
  [CLAIM_STATUSES.PARTIALLY_APPROVED]: 'Partially Approved',
  [CLAIM_STATUSES.REJECTED]: 'Rejected',
  [CLAIM_STATUSES.CANCELLED]: 'Cancelled',
  [CLAIM_STATUSES.COMPLETED]: 'Completed',
}

export const CLAIM_ITEM_TYPES = {
  REPAIR: 'REPAIR',
  REPLACEMENT: 'REPLACEMENT',
  INSPECTION: 'INSPECTION',
} as const

export type ClaimItemType = (typeof CLAIM_ITEM_TYPES)[keyof typeof CLAIM_ITEM_TYPES]

export const CLAIM_ITEM_STATUSES = {
  PENDING: 'PENDING',
  APPROVED: 'APPROVED',
  REJECTED: 'REJECTED',
  COMPLETED: 'COMPLETED',
} as const

export type ClaimItemStatus = (typeof CLAIM_ITEM_STATUSES)[keyof typeof CLAIM_ITEM_STATUSES]

export const CLAIM_ITEM_STATUS_LABELS: Record<ClaimItemStatus, string> = {
  [CLAIM_ITEM_STATUSES.PENDING]: 'Pending',
  [CLAIM_ITEM_STATUSES.APPROVED]: 'Approved',
  [CLAIM_ITEM_STATUSES.REJECTED]: 'Rejected',
  [CLAIM_ITEM_STATUSES.COMPLETED]: 'Completed',
}

export const OFFICE_TYPES = {
  EVM: 'EVM',
  SC: 'SC',
} as const

export type OfficeType = (typeof OFFICE_TYPES)[keyof typeof OFFICE_TYPES]

export const OFFICE_TYPE_LABELS: Record<OfficeType, string> = {
  [OFFICE_TYPES.EVM]: 'EVM Office',
  [OFFICE_TYPES.SC]: 'Service Center',
}

export const FILE_UPLOAD_CONFIG = {
  MAX_FILE_SIZE: 10 * 1024 * 1024, // 10MB
  ALLOWED_TYPES: [
    'image/jpeg',
    'image/png',
    'image/gif',
    'application/pdf',
    'application/msword',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
  ],
  MAX_FILES: 5,
} as const

export const PAGINATION_CONFIG = {
  DEFAULT_PAGE_SIZE: 10,
  PAGE_SIZE_OPTIONS: ['10', '20', '50', '100'],
  SHOW_SIZE_CHANGER: true,
  SHOW_QUICK_JUMPER: true,
} as const

export const ATTACHMENTS_TYPES = {
  VIDEO: 'video',
  IMAGE: 'image',
} as const

export type AttachmentType = (typeof ATTACHMENTS_TYPES)[keyof typeof ATTACHMENTS_TYPES]

export const ATTACHMENT_TYPE_LABELS: Record<AttachmentType, string> = {
  [ATTACHMENTS_TYPES.VIDEO]: 'Video',
  [ATTACHMENTS_TYPES.IMAGE]: 'Image',
}
