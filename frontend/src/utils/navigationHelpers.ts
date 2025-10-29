import { USER_ROLES, type UserRole } from '@/constants/common-constants'

export const getRoleBasedPath = (currentPath: string, defaultBasePath: string): string => {
  let basePath = defaultBasePath

  if (currentPath.includes('/admin/')) {
    basePath = `/admin${defaultBasePath}`
  } else if (currentPath.includes('/evm-staff/')) {
    basePath = `/evm-staff${defaultBasePath}`
  } else if (currentPath.includes('/sc-staff/')) {
    basePath = `/sc-staff${defaultBasePath}`
  } else if (currentPath.includes('/sc-technician/')) {
    basePath = `/sc-technician${defaultBasePath}`
  }

  return basePath
}

export const getClaimsBasePath = (currentPath: string): string => {
  return getRoleBasedPath(currentPath, '/claims')
}

export const isRolePath = (currentPath: string, role: UserRole): boolean => {
  return currentPath.includes(`/${role}/`)
}

export const getUserRoleFromPath = (currentPath: string): UserRole | null => {
  if (currentPath.includes('/admin/')) {
    return USER_ROLES.ADMIN
  } else if (currentPath.includes('/evm-staff/')) {
    return USER_ROLES.EVM_STAFF
  } else if (currentPath.includes('/sc-staff/')) {
    return USER_ROLES.SC_STAFF
  } else if (currentPath.includes('/sc-technician/')) {
    return USER_ROLES.SC_TECHNICIAN
  }
  return null
}

export const allowRoles = (currentPath: string, roles: UserRole[]): boolean => {
  const userRole = getUserRoleFromPath(currentPath)
  return !!userRole && roles.includes(userRole)
}
