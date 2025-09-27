export const USER_ROLES = {
  ADMIN: 'admin',
  SC_STAFF: 'sc staff',
  SC_TECHNICIAN: 'sc technician',
  EVM_STAFF: 'evm staff',
}

export const ROLE_ROUTES = {
  [USER_ROLES.ADMIN]: '/admin',
  [USER_ROLES.SC_STAFF]: '/sc-staff',
  [USER_ROLES.SC_TECHNICIAN]: '/sc-technician',
  [USER_ROLES.EVM_STAFF]: '/evm-staff',
}

export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: '/auth/login',
    LOGOUT: '/auth/logout',
    REFRESH: '/auth/refresh',
    GOOGLE: '/auth/google',
  },
}
