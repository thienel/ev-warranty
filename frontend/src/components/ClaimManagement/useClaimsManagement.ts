import { useCallback, useEffect, useState } from 'react'
import { claimsApi } from '@/services/claimsApi'
import useDelay from '@/hooks/useDelay'
import useHandleApiError from '@/hooks/useHandleApiError'
import type { ErrorResponse } from '@/constants/error-messages'
import type { Claim, PaginationParams } from '@/types'
import { allowRoles } from '@/utils/navigationHelpers'
import { USER_ROLES } from '@/constants/common-constants'

interface UseClaimsManagementReturn {
  claims: Claim[]
  setClaims: (claims: Claim[]) => void
  loading: boolean
  setLoading: (loading: boolean) => void
  searchText: string
  setSearchText: (text: string) => void
  updateClaim: Claim | null
  setUpdateClaim: (claim: Claim | null) => void
  isUpdate: boolean
  setIsUpdate: (isUpdate: boolean) => void
  isOpenModal: boolean
  setIsOpenModal: (isOpen: boolean) => void
  pagination: PaginationParams
  setPagination: (pagination: PaginationParams) => void
  handleOpenModal: (claim?: Claim | null, isUpdate?: boolean) => void
  fetchClaims: (params?: PaginationParams) => Promise<void>
  handleReset: () => Promise<void>
  handleSubmit: (claimId: string) => Promise<void>
  handleCancel: (claimId: string) => Promise<void>
  handleComplete: (claimId: string) => Promise<void>
  handleReview: (claimId: string) => Promise<void>
  handleRequestInfo: (claimId: string) => Promise<void>
  allowCreate: boolean
}

const useClaimsManagement = (): UseClaimsManagementReturn => {
  const [claims, setClaims] = useState<Claim[]>([])
  const [loading, setLoading] = useState(false)
  const [searchText, setSearchText] = useState('')
  const [updateClaim, setUpdateClaim] = useState<Claim | null>(null)
  const [isUpdate, setIsUpdate] = useState(false)
  const [isOpenModal, setIsOpenModal] = useState(false)
  const [pagination, setPagination] = useState<PaginationParams>({
    page: 1,
    limit: 10,
  })
  const [allowCreate, setAllowCreate] = useState(false)
  const handleError = useHandleApiError()
  const delay = useDelay(300)

  const handleOpenModal = (claim: Claim | null = null, isUpdate = false) => {
    setUpdateClaim(claim)
    setIsUpdate(isUpdate)
    setIsOpenModal(true)
  }

  const fetchClaims = useCallback(
    async (params?: PaginationParams) => {
      try {
        setLoading(true)
        const response = await claimsApi.getAll(params || pagination)
        const claimsData = Array.isArray(response.data?.claims) ? response.data.claims : []
        console.log('Fetched claims:', response)
        setClaims(claimsData)
      } catch (error) {
        handleError(error as ErrorResponse)
        setClaims([])
      } finally {
        setLoading(false)
      }
    },
    [pagination, handleError],
  )

  const handleReset = async () => {
    setLoading(true)
    delay(async () => {
      setSearchText('')
      setIsOpenModal(false)
      setUpdateClaim(null)
      await fetchClaims()
      setLoading(false)
    })
  }

  // Claim action handlers
  const handleSubmit = async (claimId: string) => {
    try {
      setLoading(true)
      await claimsApi.submit(claimId)
      await fetchClaims()
    } catch (error) {
      handleError(error as ErrorResponse)
    } finally {
      setLoading(false)
    }
  }

  const handleCancel = async (claimId: string) => {
    try {
      setLoading(true)
      await claimsApi.cancel(claimId)
      await fetchClaims()
    } catch (error) {
      handleError(error as ErrorResponse)
    } finally {
      setLoading(false)
    }
  }

  const handleComplete = async (claimId: string) => {
    try {
      setLoading(true)
      await claimsApi.complete(claimId)
      await fetchClaims()
    } catch (error) {
      handleError(error as ErrorResponse)
    } finally {
      setLoading(false)
    }
  }

  const handleReview = async (claimId: string) => {
    try {
      setLoading(true)
      await claimsApi.review(claimId)
      await fetchClaims()
    } catch (error) {
      handleError(error as ErrorResponse)
    } finally {
      setLoading(false)
    }
  }

  const handleRequestInfo = async (claimId: string) => {
    try {
      setLoading(true)
      await claimsApi.requestInfo(claimId)
      await fetchClaims()
    } catch (error) {
      handleError(error as ErrorResponse)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (allowRoles(location.pathname, [USER_ROLES.SC_STAFF])) {
      setAllowCreate(true)
    }

    fetchClaims()
  }, [fetchClaims])

  return {
    claims,
    setClaims,
    loading,
    setLoading,
    searchText,
    setSearchText,
    updateClaim,
    setUpdateClaim,
    isUpdate,
    setIsUpdate,
    isOpenModal,
    setIsOpenModal,
    pagination,
    setPagination,
    handleOpenModal,
    fetchClaims,
    handleReset,
    handleSubmit,
    handleCancel,
    handleComplete,
    handleReview,
    handleRequestInfo,
    allowCreate,
  }
}

export default useClaimsManagement
