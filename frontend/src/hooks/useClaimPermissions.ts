import { useCallback } from 'react'
import { useSelector } from 'react-redux'
import { CLAIM_STATUSES, USER_ROLES } from '@constants/common-constants'
import type { RootState } from '@redux/store'
import type { ClaimDetail } from '@/types/index'

interface UseClaimPermissionsReturn {
  canAddItems: boolean
  canEditClaim: boolean
  canDeleteClaim: boolean
  canApproveClaim: boolean
  canRejectClaim: boolean
  canViewClaim: boolean
}

export const useClaimPermissions = (claim: ClaimDetail | null): UseClaimPermissionsReturn => {
  const { user } = useSelector((state: RootState) => state.auth)

  const canAddItems = useCallback((): boolean => {
    if (!user || !claim) return false

    // Only SC_STAFF and SC_TECHNICIAN can add items
    const allowedRoles = [USER_ROLES.SC_STAFF, USER_ROLES.SC_TECHNICIAN] as const
    if (!allowedRoles.includes(user.role as (typeof allowedRoles)[number])) return false

    // Only allow adding items when claim is in DRAFT or REQUEST_INFO status
    const allowedStatuses = [CLAIM_STATUSES.DRAFT, CLAIM_STATUSES.REQUEST_INFO] as const
    return allowedStatuses.includes(claim.status as (typeof allowedStatuses)[number])
  }, [user, claim])

  const canEditClaim = useCallback((): boolean => {
    if (!user || !claim) return false

    // Only SC_STAFF and SC_TECHNICIAN can edit claims
    const allowedRoles = [USER_ROLES.SC_STAFF, USER_ROLES.SC_TECHNICIAN] as const
    if (!allowedRoles.includes(user.role as (typeof allowedRoles)[number])) return false

    // Only allow editing when claim is in DRAFT or REQUEST_INFO status
    const allowedStatuses = [CLAIM_STATUSES.DRAFT, CLAIM_STATUSES.REQUEST_INFO] as const
    return allowedStatuses.includes(claim.status as (typeof allowedStatuses)[number])
  }, [user, claim])

  const canDeleteClaim = useCallback((): boolean => {
    if (!user || !claim) return false

    // Only SC_STAFF can delete claims
    if (user.role !== USER_ROLES.SC_STAFF) return false

    // Only allow deleting when claim is in DRAFT status
    return claim.status === CLAIM_STATUSES.DRAFT
  }, [user, claim])

  const canApproveClaim = useCallback((): boolean => {
    if (!user || !claim) return false

    // Only EVM_STAFF can approve claims
    if (user.role !== USER_ROLES.EVM_STAFF) return false

    // Only allow approving when claim is in REVIEWING status
    return claim.status === CLAIM_STATUSES.REVIEWING
  }, [user, claim])

  const canRejectClaim = useCallback((): boolean => {
    if (!user || !claim) return false

    // Only EVM_STAFF can reject claims
    if (user.role !== USER_ROLES.EVM_STAFF) return false

    // Only allow rejecting when claim is in REVIEWING status
    return claim.status === CLAIM_STATUSES.REVIEWING
  }, [user, claim])

  const canViewClaim = useCallback((): boolean => {
    if (!user) return false

    // All authenticated users can view claims
    return true
  }, [user])

  return {
    canAddItems: canAddItems(),
    canEditClaim: canEditClaim(),
    canDeleteClaim: canDeleteClaim(),
    canApproveClaim: canApproveClaim(),
    canRejectClaim: canRejectClaim(),
    canViewClaim: canViewClaim(),
  }
}

export default useClaimPermissions
