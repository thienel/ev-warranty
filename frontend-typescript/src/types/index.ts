// API Response Types
export interface ApiSuccessResponse<T = unknown> {
  data: T
}

export interface ApiErrorResponse {
  error: string
}

export interface PaginationParams {
  page?: number
  limit?: number
  status?: string
}

// Base types
export interface BaseEntity {
  id: string
  created_at?: string
  updated_at?: string
}

// User types (matching dtos.UserDTO from Swagger)
export interface User extends Record<string, unknown> {
  id: string
  email: string
  name: string
  role: string
  office_id: string
  is_active: boolean
}

export interface UserFormData {
  name: string
  email: string
  password?: string
  role: string
  office_id: string
  is_active: boolean
}

// User API DTOs
export interface CreateUserRequest {
  email: string
  name: string
  password: string
  role: string
  office_id: string
  is_active: boolean
}

export interface UpdateUserRequest {
  name?: string
  role?: string
  office_id?: string
  is_active?: boolean
}

// Auth types
export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface ValidateTokenResponse {
  valid: boolean
  user: User
}

export interface RefreshTokenResponse {
  token: string
}

// Auth state types
export interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  remember: boolean
  isLoading?: boolean
}

// Login payload type
export interface LoginPayload {
  user: User
  token: string
  remember: boolean
}

// Office types (matching entities.Office from Swagger)
export interface Office extends Record<string, unknown> {
  id: string
  office_name: string
  office_type: string
  address: string
  is_active: boolean
  created_at?: string
  updated_at?: string
}

export interface OfficeFormData {
  office_name: string
  office_type: 'evm' | 'sc'
  address: string
  is_active: boolean
}

// Office API DTOs
export interface CreateOfficeRequest {
  office_name: string
  office_type: string
  address: string
  is_active?: boolean
}

export interface UpdateOfficeRequest {
  office_name?: string
  office_type?: string
  address?: string
  is_active?: boolean
}

// Claim types (matching entities.Claim from Swagger)
export interface Claim extends Record<string, unknown> {
  id: string
  customer_id: string
  vehicle_id: string
  description: string
  status: string
  total_cost: number
  approved_by?: string
  created_at?: string
  updated_at?: string
}

export interface ClaimFormData {
  customer_id: string
  vehicle_id: string
  description: string
}

// Claim API DTOs
export interface CreateClaimRequest {
  customer_id: string
  vehicle_id: string
  description: string
}

export interface UpdateClaimRequest {
  description: string
}

// Claim Item types (matching entities.ClaimItem from Swagger)
export interface ClaimItem extends Record<string, unknown> {
  id: string
  claim_id: string
  part_category_id: number
  faulty_part_id: string
  replacement_part_id?: string
  issue_description: string
  type: string
  cost: number
  status: string
  created_at?: string
  updated_at?: string
}

// Claim Item API DTOs
export interface CreateClaimItemRequest {
  part_category_id: number
  faulty_part_id: string
  replacement_part_id?: string
  issue_description: string
  type: string
  cost: number
}

export interface ClaimItemListResponse {
  items: ClaimItem[]
  total: number
}

// Claim Attachment types (matching entities.ClaimAttachment from Swagger)
export interface ClaimAttachment extends Record<string, unknown> {
  id: string
  claimID: string
  url: string
  type: string
  created_at?: string
}

export interface ClaimAttachmentListResponse {
  attachments: ClaimAttachment[]
  total: number
}

// Claim History types (matching entities.ClaimHistory from Swagger)
export interface ClaimHistory extends Record<string, unknown> {
  id: string
  claim_id: string
  status: string
  changed_by: string
  changedAt: string
}

// Table column types
export interface SortInfo {
  columnKey?: string
  order?: 'ascend' | 'descend' | null
}

export interface FilterInfo {
  [key: string]: React.Key[] | null
}

// Additional props for table columns
export interface TableAdditionalProps {
  getOfficeName?: (officeId: string) => string
  [key: string]: unknown
}

// Modal component props
export interface BaseModalProps {
  loading: boolean
  setLoading: (loading: boolean) => void
  onClose: () => void
  opened: boolean
  isUpdate: boolean
}

export interface UserModalProps extends BaseModalProps {
  user?: User | null
  offices: Office[]
}

export interface OfficeModalProps extends BaseModalProps {
  office?: Office | null
}