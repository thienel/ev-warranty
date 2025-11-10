import { useCallback } from 'react'
import { useSelector } from 'react-redux'
import { CLAIM_STATUSES, USER_ROLES } from '@constants/common-constants'
import type { RootState } from '@redux/store'
import type { ClaimDetail, ClaimItem } from '@/types/index'

interface UseClaimPermissionsReturn {
  canAddItems: boolean
  canEditClaim: boolean
  canDeleteClaim: boolean
  canCancelClaim: boolean
  canApproveClaim: boolean
  canRejectClaim: boolean
  canStartReview: boolean
  canApproveClaimItems: boolean
  canRejectClaimItems: boolean
  canCompleteClaim: boolean
  canViewWarrantyPolicy: boolean
  canViewPolicyCoverage: boolean
  canViewClaim: boolean
  canAddAttachments: boolean
  canSubmitClaim: boolean
}

export const useClaimPermissions = (
  claim: ClaimDetail | null,
  claimItems: ClaimItem[] = [],
): UseClaimPermissionsReturn => {
  const { user } = useSelector((state: RootState) => state.auth)

  const canAddItems = useCallback((): boolean => {
    if (!user || !claim) return false

    // Only SC_STAFF and SC_TECHNICIAN can add items
    const allowedRoles = [USER_ROLES.SC_STAFF, USER_ROLES.SC_TECHNICIAN] as const
    if (!allowedRoles.includes(user.role as (typeof allowedRoles)[number])) return false

    // Only allow adding items when claim is in DRAFT status
    const allowedStatuses = [CLAIM_STATUSES.DRAFT] as const
    return allowedStatuses.includes(claim.status as (typeof allowedStatuses)[number])
  }, [user, claim])

  const canEditClaim = useCallback((): boolean => {
    if (!user || !claim) return false

    // Only SC_STAFF and SC_TECHNICIAN can edit claims
    const allowedRoles = [USER_ROLES.SC_STAFF, USER_ROLES.SC_TECHNICIAN] as const
    if (!allowedRoles.includes(user.role as (typeof allowedRoles)[number])) return false

    // Only allow editing when claim is in DRAFT status
    const allowedStatuses = [CLAIM_STATUSES.DRAFT] as const
    return allowedStatuses.includes(claim.status as (typeof allowedStatuses)[number])
  }, [user, claim])

  const canDeleteClaim = useCallback((): boolean => {
    if (!user || !claim) return false

    // Only SC_STAFF can delete claims
    if (user.role !== USER_ROLES.SC_STAFF) return false

    // Only allow deleting when claim is in DRAFT status
    return claim.status === CLAIM_STATUSES.DRAFT
  }, [user, claim])

  const canCancelClaim = useCallback((): boolean => {
    if (!user || !claim) return false

    // Only SC_STAFF can cancel claims
    if (user.role !== USER_ROLES.SC_STAFF) return false

    // Only allow canceling when claim is in SUBMITTED status
    const allowedStatuses = [CLAIM_STATUSES.SUBMITTED] as const
    return allowedStatuses.includes(claim.status as (typeof allowedStatuses)[number])
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

  const canStartReview = useCallback((): boolean => {
    if (!user || !claim) return false

    // Only EVM_STAFF can start review
    if (user.role !== USER_ROLES.EVM_STAFF) return false

    // Only allow starting review when claim is in SUBMITTED status
    return claim.status === CLAIM_STATUSES.SUBMITTED
  }, [user, claim])

  const canApproveClaimItems = useCallback((): boolean => {
    if (!user || !claim) return false

    // Only EVM_STAFF can approve claim items
    if (user.role !== USER_ROLES.EVM_STAFF) return false

    // Only allow approving claim items when claim is in REVIEWING status
    return claim.status === CLAIM_STATUSES.REVIEWING
  }, [user, claim])

  const canRejectClaimItems = useCallback((): boolean => {
    if (!user || !claim) return false

    // Only EVM_STAFF can reject claim items
    if (user.role !== USER_ROLES.EVM_STAFF) return false

    // Only allow rejecting claim items when claim is in REVIEWING status
    return claim.status === CLAIM_STATUSES.REVIEWING
  }, [user, claim])

  const canCompleteClaim = useCallback((): boolean => {
    if (!user || !claim) return false

    // Only EVM_STAFF can complete claims
    if (user.role !== USER_ROLES.EVM_STAFF) return false

    // Only allow completing when claim is in REVIEWING status
    if (claim.status !== CLAIM_STATUSES.REVIEWING) return false

    // Check if all claim items have been processed (approved or rejected)
    if (claimItems.length === 0) return false

    const allItemsProcessed = claimItems.every(
      (item) => item.status === 'APPROVED' || item.status === 'REJECTED',
    )

    return allItemsProcessed
  }, [user, claim, claimItems])

  const canViewWarrantyPolicy = useCallback((): boolean => {
    if (!user || !claim) return false

    // Only EVM_STAFF can view warranty policy during review
    if (user.role !== USER_ROLES.EVM_STAFF) return false

    // Only show warranty policy when claim is in REVIEWING status
    return claim.status === CLAIM_STATUSES.REVIEWING
  }, [user, claim])

  const canViewPolicyCoverage = useCallback((): boolean => {
    if (!user || !claim) return false

    // Only EVM_STAFF can view policy coverage during review
    if (user.role !== USER_ROLES.EVM_STAFF) return false

    // Only show policy coverage when claim is in REVIEWING status
    return claim.status === CLAIM_STATUSES.REVIEWING
  }, [user, claim])

  const canViewClaim = useCallback((): boolean => {
    if (!user) return false

    // All authenticated users can view claims
    return true
  }, [user])

  const canAddAttachments = useCallback((): boolean => {
    if (!user || !claim) return false

    // Only SC_TECHNICIAN can add attachments
    if (user.role !== USER_ROLES.SC_TECHNICIAN) return false

    // Only allow adding attachments when claim is in DRAFT status
    const allowedStatuses = [CLAIM_STATUSES.DRAFT] as const
    return allowedStatuses.includes(claim.status as (typeof allowedStatuses)[number])
  }, [user, claim])

  const canSubmitClaim = useCallback((): boolean => {
    if (!user || !claim) return false

    // Only SC_STAFF can submit claims
    if (user.role !== USER_ROLES.SC_STAFF) return false

    // Only allow submitting when claim is in DRAFT status
    const allowedStatuses = [CLAIM_STATUSES.DRAFT] as const
    return allowedStatuses.includes(claim.status as (typeof allowedStatuses)[number])
  }, [user, claim])

  return {
    canAddItems: canAddItems(),
    canEditClaim: canEditClaim(),
    canDeleteClaim: canDeleteClaim(),
    canCancelClaim: canCancelClaim(),
    canApproveClaim: canApproveClaim(),
    canRejectClaim: canRejectClaim(),
    canStartReview: canStartReview(),
    canApproveClaimItems: canApproveClaimItems(),
    canRejectClaimItems: canRejectClaimItems(),
    canCompleteClaim: canCompleteClaim(),
    canViewWarrantyPolicy: canViewWarrantyPolicy(),
    canViewPolicyCoverage: canViewPolicyCoverage(),
    canViewClaim: canViewClaim(),
    canAddAttachments: canAddAttachments(),
    canSubmitClaim: canSubmitClaim(),
  }
}

export default useClaimPermissions
