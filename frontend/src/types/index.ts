import type {
  UserRole,
  OfficeType,
  ClaimStatus,
  ClaimItemStatus,
  ClaimItemType,
  AttachmentType,
} from '../constants/common-constants'

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

export interface BaseEntity {
  id: string
  created_at?: string
  updated_at?: string
}

export interface User {
  id: string
  email: string
  name: string
  role: UserRole
  office_id: string
  is_active?: boolean
}

export interface UserFormData {
  name: string
  email: string
  password?: string
  role: UserRole
  office_id: string
  is_active: boolean
}

export interface CreateUserRequest {
  email: string
  name: string
  password: string
  role: UserRole
  office_id: string
  is_active?: boolean
}

export interface UpdateUserRequest {
  name?: string
  role?: UserRole
  office_id?: string
  is_active?: boolean
}

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

export interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  remember: boolean
  isLoading?: boolean
}

export interface LoginPayload {
  user: User
  token: string
  remember: boolean
}

export interface Office {
  id: string
  office_name: string
  office_type: OfficeType
  address: string
  is_active?: boolean
  created_at?: string
  updated_at?: string
}

export interface OfficeFormData {
  office_name: string
  office_type: OfficeType
  address: string
  is_active: boolean
}

export interface CreateOfficeRequest {
  office_name: string
  office_type: OfficeType
  address: string
  is_active?: boolean
}

export interface UpdateOfficeRequest {
  office_name?: string
  office_type?: OfficeType
  address?: string
  is_active?: boolean
}

export interface Claim {
  id: string
  customer_id: string
  vehicle_id: string
  description: string
  status: ClaimStatus
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

export interface CreateClaimRequest {
  customer_id: string
  vehicle_id: string
  description: string
}

export interface UpdateClaimRequest {
  description: string
}

export interface ClaimListResponse {
  claims: Claim[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

export interface ClaimItem {
  id: string
  claim_id: string
  part_category_id: number
  faulty_part_id: string
  replacement_part_id?: string
  issue_description: string
  type: ClaimItemType
  cost: number
  status: ClaimItemStatus
  created_at?: string
  updated_at?: string
}

export interface CreateClaimItemRequest {
  part_category_id: number
  faulty_part_id: string
  replacement_part_id?: string
  issue_description: string
  type: ClaimItemType
  cost: number
}

export interface ClaimItemListResponse {
  items: ClaimItem[]
  total: number
}

export interface ClaimAttachment {
  id: string
  claimID: string
  url: string
  type: AttachmentType
  created_at?: string
}

export interface ClaimAttachmentListResponse {
  attachments: ClaimAttachment[]
  total: number
}

export interface ClaimHistory {
  id: string
  claim_id: string
  status: ClaimStatus
  changed_by: string
  changed_at: string
}

export interface SortInfo {
  columnKey?: string
  order?: 'ascend' | 'descend' | null
}

export interface FilterInfo {
  [key: string]: React.Key[] | null
}

export interface TableAdditionalProps {
  getOfficeName?: (officeId: string) => string
  [key: string]: unknown
}

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
  officesLoading?: boolean
}

export interface OfficeModalProps extends BaseModalProps {
  office?: Office | null
}

export interface Customer {
  id: string
  first_name: string
  last_name: string
  phone_number?: string
  email?: string
  address?: string
  created_at: string
  updated_at?: string
  deleted_at?: string
  is_deleted: boolean
  full_name?: string
}

export interface CustomerFormData {
  first_name: string
  last_name: string
  phone_number?: string
  email?: string
  address?: string
}

export interface CreateCustomerRequest {
  first_name: string
  last_name: string
  email?: string
  phone_number?: string
  address?: string
}

export interface UpdateCustomerRequest {
  first_name: string
  last_name: string
  email?: string
  phone_number?: string
  address?: string
}

export interface VehicleModel {
  id: string
  brand: string
  model_name: string
  year: number
  created_at: string
  updated_at?: string
}

export interface VehicleModelFormData {
  brand: string
  model_name: string
  year: number
}

export interface CreateVehicleModelRequest {
  brand: string
  model_name: string
  year: number
}

export interface UpdateVehicleModelRequest {
  brand: string
  model_name: string
  year: number
}

export interface Vehicle {
  id: string
  vin: string
  license_plate?: string
  customer_id: string
  model_id: string
  purchase_date?: string
  created_at: string
  updated_at?: string
}

export interface VehicleFormData {
  vin: string
  license_plate?: string
  customer_id: string
  model_id: string
  purchase_date?: unknown
}

export interface CreateVehicleRequest {
  vin: string
  license_plate?: string
  customer_id: string
  model_id: string
  purchase_date?: string
}

export interface UpdateVehicleRequest {
  vin: string
  license_plate?: string
  customer_id: string
  model_id: string
  purchase_date?: string
}

export interface DotNetApiResponse<T = unknown> {
  is_success: boolean
  message?: string
  error?: string
  data?: T
}

export interface CustomerModalProps extends BaseModalProps {
  customer?: Customer | null
}

export interface VehicleModelModalProps extends BaseModalProps {
  vehicleModel?: VehicleModel | null
}

export interface VehicleModalProps extends BaseModalProps {
  vehicle?: Vehicle | null
  customers: Customer[]
  vehicleModels: VehicleModel[]
  customersLoading?: boolean
  vehicleModelsLoading?: boolean
}
