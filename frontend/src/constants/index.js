export const USER_ROLES = {
  ADMIN: 'admin',
  SC_STAFF: 'sc staff',
  SC_TECHNICIAN: 'sc technician',
  EVM_STAFF: 'evm staff',
}

export const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost/api/v1'

export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: '/auth/login',
    LOGOUT: '/auth/logout',
    GOOGLE: '/auth/google',
    TOKEN: '/auth/token',
  },
}
