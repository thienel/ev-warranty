interface JwtPayload {
  exp: number
  iat: number
  sub: string
  [key: string]: unknown
}

/**
 * Decode JWT token payload (basic implementation)
 */
function decodeJwtPayload(token: string): JwtPayload | null {
  try {
    const parts = token.split('.')
    if (parts.length !== 3) {
      return null
    }
    
    const payload = parts[1]
    const decoded = atob(payload.replace(/-/g, '+').replace(/_/g, '/'))
    return JSON.parse(decoded) as JwtPayload
  } catch (error) {
    console.error('Error decoding JWT payload:', error)
    return null
  }
}

/**
 * Check if a JWT token is expired
 */
export function isTokenExpired(token: string | null): boolean {
  if (!token) return true
  
  const decoded = decodeJwtPayload(token)
  if (!decoded || !decoded.exp) return true
  
  const currentTime = Date.now() / 1000
  
  // Add 30 seconds buffer before expiration
  return decoded.exp < (currentTime + 30)
}

/**
 * Check if a JWT token is valid (not expired and properly formatted)
 */
export function isTokenValid(token: string | null): boolean {
  if (!token) return false
  
  const decoded = decodeJwtPayload(token)
  if (!decoded) return false
  
  const currentTime = Date.now() / 1000
  
  // Check if token has required fields and is not expired
  return !!(decoded.exp && decoded.sub && decoded.exp > currentTime)
}

/**
 * Get token expiration time
 */
export function getTokenExpiration(token: string | null): Date | null {
  if (!token) return null
  
  const decoded = decodeJwtPayload(token)
  if (!decoded || !decoded.exp) return null
  
  return new Date(decoded.exp * 1000)
}

/**
 * Get user ID from token
 */
export function getUserIdFromToken(token: string | null): string | null {
  if (!token) return null
  
  const decoded = decodeJwtPayload(token)
  return decoded?.sub || null
}