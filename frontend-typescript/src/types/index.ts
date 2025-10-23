// Base types
export interface BaseEntity {
  id: string
  created_at?: string
  updated_at?: string
}

// User types
export interface User {
  id: string
  email: string
  name: string
  role: string
  office_id?: string
  is_active: boolean
  created_at?: string
  updated_at?: string
}

export interface UserFormData {
  name: string
  email: string
  password?: string
  role: string
  office_id: string
  is_active: boolean
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

// Office types
export interface Office {
  id: string
  office_name: string
  office_type: 'evm' | 'sc'
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

// Claim types
export interface Claim {
  id: string
  status: 'SUBMITTED' | 'APPROVED' | 'REJECTED' | 'PROCESSING' | 'COMPLETED'
  customer_id: string
  customer_name: string
  vehicle_id: string
  vehicle_info: string
  description: string
  total_cost: number
  created_at?: string
  updated_at?: string
}

export interface ClaimFormData {
  status: string
  customer_id: string
  vehicle_id: string
  description: string
  total_cost: number
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