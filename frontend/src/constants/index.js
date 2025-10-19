export const PASSWORD_RULES = [
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

export const EMAIL_RULES = [
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
}

export const ROLE_LABELS = {
  [USER_ROLES.ADMIN]: 'Admin',
  [USER_ROLES.SC_STAFF]: 'SC Staff',
  [USER_ROLES.SC_TECHNICIAN]: 'SC Technician',
  [USER_ROLES.EVM_STAFF]: 'EVM Staff',
}

export const ERROR_MESSAGES = {
  [403]: 'Sorry, you are not authorized to access this page.',
  [500]: 'Sorry, something went wrong on the server.',
  [404]: 'Sorry, the page you visited does not exist.',
}

export const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost/api/v1'

export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: '/auth/login',
    LOGOUT: '/auth/logout',
    GOOGLE: '/auth/google',
    TOKEN: '/auth/token',
  },
  USER: '/users',
  OFFICE: '/offices',
}
