import { useCallback, useEffect, useState } from 'react'
import { claimsApi } from '@/services/claimsApi'
import { customersApi } from '@/services/customersApi'
import useDelay from '@/hooks/useDelay'
import useHandleApiError from '@/hooks/useHandleApiError'
import type { ErrorResponse } from '@/constants/error-messages'
import type { Claim, PaginationParams, Customer } from '@/types'
import { allowRoles } from '@/utils/navigationHelpers'
import { USER_ROLES } from '@/constants/common-constants'

interface EnrichedClaim extends Claim {
  customer_name?: string
}

interface UseClaimsManagementReturn {
  claims: EnrichedClaim[]
  setClaims: (claims: EnrichedClaim[]) => void
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
  const [claims, setClaims] = useState<EnrichedClaim[]>([])
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
        const claimsData = Array.isArray(response.data) ? response.data : []

        const customersResponse = await customersApi.getAll()

        let customers: Customer[] = []
        if (customersResponse.data) {
          if (Array.isArray(customersResponse.data)) {
            customers = customersResponse.data
          } else if (
            typeof customersResponse.data === 'object' &&
            'data' in customersResponse.data
          ) {
            const nestedData = (customersResponse.data as Record<string, unknown>).data
            customers = Array.isArray(nestedData) ? nestedData : []
          }
        }

        const customerMap = new Map<string, string>()
        if (customers.length > 0) {
          customers.forEach((customer: Customer) => {
            const fullName = `${customer.first_name} ${customer.last_name}`.trim()
            customerMap.set(customer.id, fullName)
          })
        }

        const enrichedClaims: EnrichedClaim[] = claimsData.map((claim) => ({
          ...claim,
          customer_name: customerMap.get(claim.customer_id) || claim.customer_id,
        }))

        setClaims(enrichedClaims)
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
