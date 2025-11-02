import { useState, useEffect, useCallback } from 'react'
import { message } from 'antd'
import {
  claims as claimsApi,
  claimItems as claimItemsApi,
  claimAttachments as claimAttachmentsApi,
  customersApi,
  vehiclesApi,
  vehicleModelsApi,
  partCategoriesApi,
  partsApi,
} from '@services/index'
import useHandleApiError from '@/hooks/useHandleApiError'
import type {
  ClaimDetail,
  ClaimItem,
  ClaimAttachment,
  Customer,
  VehicleDetail,
  VehicleModel,
  PartCategory,
  Part,
} from '@/types/index'

interface UseClaimDataReturn {
  // Data
  claim: ClaimDetail | null
  customer: Customer | null
  vehicle: VehicleDetail | null
  claimItems: ClaimItem[]
  attachments: ClaimAttachment[]
  partCategories: PartCategory[]
  parts: Part[]

  // Loading states
  claimLoading: boolean
  customerLoading: boolean
  vehicleLoading: boolean
  itemsLoading: boolean
  attachmentsLoading: boolean

  // Refetch functions
  refetchClaim: () => Promise<void>
  refetchClaimItems: () => Promise<void>
  refetchAttachments: () => Promise<void>
}

export const useClaimData = (claimId?: string): UseClaimDataReturn => {
  const handleError = useHandleApiError()

  // Data state
  const [claim, setClaim] = useState<ClaimDetail | null>(null)
  const [customer, setCustomer] = useState<Customer | null>(null)
  const [vehicle, setVehicle] = useState<VehicleDetail | null>(null)
  const [claimItems, setClaimItems] = useState<ClaimItem[]>([])
  const [attachments, setAttachments] = useState<ClaimAttachment[]>([])
  const [partCategories, setPartCategories] = useState<PartCategory[]>([])
  const [parts, setParts] = useState<Part[]>([])

  // Loading states
  const [claimLoading, setClaimLoading] = useState(false)
  const [customerLoading, setCustomerLoading] = useState(false)
  const [vehicleLoading, setVehicleLoading] = useState(false)
  const [itemsLoading, setItemsLoading] = useState(false)
  const [attachmentsLoading, setAttachmentsLoading] = useState(false)

  // Fetch claim details
  const fetchClaim = useCallback(async () => {
    if (!claimId) return

    try {
      setClaimLoading(true)
      const response = await claimsApi.getById(claimId)
      let claimData = response.data

      if (claimData && typeof claimData === 'object' && 'data' in claimData) {
        claimData = (claimData as { data: unknown }).data as ClaimDetail
      }

      setClaim(claimData as ClaimDetail)
    } catch (error) {
      handleError(error as Error)
      message.error('Failed to load claim details')
    } finally {
      setClaimLoading(false)
    }
  }, [claimId, handleError])

  // Fetch customer details
  const fetchCustomer = useCallback(
    async (customerId: string) => {
      try {
        setCustomerLoading(true)
        const customerResponse = await customersApi.getById(customerId)
        let customerData = customerResponse.data
        if (customerData && typeof customerData === 'object' && 'data' in customerData) {
          customerData = (customerData as { data: unknown }).data as Customer
        }
        setCustomer(customerData as Customer)
      } catch (error) {
        handleError(error as Error)
      } finally {
        setCustomerLoading(false)
      }
    },
    [handleError],
  )

  // Fetch vehicle details with model info
  const fetchVehicle = useCallback(
    async (vehicleId: string) => {
      try {
        setVehicleLoading(true)
        const vehicleResponse = await vehiclesApi.getById(vehicleId)
        let vehicleData = vehicleResponse.data
        if (vehicleData && typeof vehicleData === 'object' && 'data' in vehicleData) {
          vehicleData = (vehicleData as { data: unknown }).data as VehicleDetail
        }

        const vehicleInfo = vehicleData as VehicleDetail

        // Fetch vehicle model if model_id exists
        if (vehicleInfo.model_id) {
          try {
            const modelResponse = await vehicleModelsApi.getById(vehicleInfo.model_id)
            let modelData: unknown = modelResponse.data
            if (modelData && typeof modelData === 'object' && 'data' in modelData) {
              modelData = (modelData as { data: VehicleModel }).data
            }
            vehicleInfo.model = modelData as VehicleModel
          } catch (error) {
            console.error('Failed to fetch vehicle model:', error)
          }
        }

        setVehicle(vehicleInfo)
      } catch (error) {
        handleError(error as Error)
      } finally {
        setVehicleLoading(false)
      }
    },
    [handleError],
  )

  // Fetch parts for claim items
  const fetchPartsForClaimItems = useCallback(
    async (claimItems: ClaimItem[]) => {
      if (claimItems.length === 0) return

      try {
        // Get unique part IDs from claim items
        const uniquePartIds = Array.from(new Set(claimItems.map((item) => item.faulty_part_id)))

        // Fetch each part individually
        const partPromises = uniquePartIds.map((partId) => partsApi.getById(partId))
        const partResponses = await Promise.allSettled(partPromises)

        const fetchedParts: Part[] = []
        let hasFailures = false

        partResponses.forEach((response, index) => {
          if (response.status === 'fulfilled') {
            let partData = response.value.data
            if (partData && typeof partData === 'object' && 'data' in partData) {
              partData = (partData as { data: unknown }).data as Part
            }
            fetchedParts.push(partData as Part)
          } else {
            console.warn(`Part ${uniquePartIds[index]} not found, it may have been deleted`)
            hasFailures = true
          }
        })

        // If we couldn't fetch some parts individually, fall back to fetching all parts
        if (hasFailures && fetchedParts.length < uniquePartIds.length) {
          console.warn('Some parts not found individually, falling back to fetch all parts')
          try {
            const response = await partsApi.getAll()
            let partsData = response.data
            if (partsData && typeof partsData === 'object' && 'data' in partsData) {
              partsData = (partsData as { data: unknown }).data as Part[]
            }
            setParts(partsData as Part[])
            return
          } catch (fallbackError) {
            console.error('Failed to fetch all parts as fallback:', fallbackError)
          }
        }

        setParts(fetchedParts)
      } catch (error) {
        handleError(error as Error)
        setParts([])
      }
    },
    [handleError],
  )

  // Fetch claim items
  const fetchClaimItems = useCallback(async () => {
    if (!claimId) return

    try {
      setItemsLoading(true)
      const response = await claimItemsApi.getByClaimId(claimId)
      const itemsData: unknown = response.data

      // Backend returns: { data: ClaimItem[] } not { data: { items: ClaimItem[] } }
      let items: ClaimItem[] = []
      if (itemsData && typeof itemsData === 'object' && 'data' in itemsData) {
        const nestedData = (itemsData as { data: unknown }).data
        if (Array.isArray(nestedData)) {
          items = nestedData
          setClaimItems(nestedData)
        } else {
          setClaimItems([])
        }
      } else if (Array.isArray(itemsData)) {
        items = itemsData
        setClaimItems(itemsData)
      } else {
        setClaimItems([])
      }

      // Fetch parts for the claim items
      if (items.length > 0) {
        await fetchPartsForClaimItems(items)
      }
    } catch (error) {
      handleError(error as Error)
      setClaimItems([])
    } finally {
      setItemsLoading(false)
    }
  }, [claimId, handleError, fetchPartsForClaimItems])

  // Fetch claim attachments
  const fetchAttachments = useCallback(async () => {
    if (!claimId) return

    try {
      setAttachmentsLoading(true)
      const response = await claimAttachmentsApi.getByClaimId(claimId)
      const attachmentsData: unknown = response.data

      // Backend returns: { data: ClaimAttachment[] } not { data: { attachments: ClaimAttachment[] } }
      if (attachmentsData && typeof attachmentsData === 'object' && 'data' in attachmentsData) {
        const nestedData = (attachmentsData as { data: unknown }).data
        if (Array.isArray(nestedData)) {
          setAttachments(nestedData)
        } else {
          setAttachments([])
        }
      } else if (Array.isArray(attachmentsData)) {
        setAttachments(attachmentsData)
      } else {
        setAttachments([])
      }
    } catch (error) {
      handleError(error as Error)
      setAttachments([])
    } finally {
      setAttachmentsLoading(false)
    }
  }, [claimId, handleError])

  // Fetch part categories
  const fetchPartCategories = useCallback(async () => {
    try {
      const response = await partCategoriesApi.getAll()
      let categoriesData = response.data
      if (categoriesData && typeof categoriesData === 'object' && 'data' in categoriesData) {
        categoriesData = (categoriesData as { data: unknown }).data as PartCategory[]
      }
      setPartCategories(categoriesData as PartCategory[])
    } catch (error) {
      handleError(error as Error)
      setPartCategories([])
    }
  }, [handleError])

  // Initial data fetch
  useEffect(() => {
    if (claimId) {
      fetchClaim()
      fetchClaimItems()
      fetchAttachments()
    }
    // Fetch part categories for displaying category names
    fetchPartCategories()
  }, [claimId, fetchClaim, fetchClaimItems, fetchAttachments, fetchPartCategories])

  // Fetch customer and vehicle when claim is loaded
  useEffect(() => {
    if (claim) {
      if (claim.customer_id) {
        fetchCustomer(claim.customer_id)
      }
      if (claim.vehicle_id) {
        fetchVehicle(claim.vehicle_id)
      }
    }
  }, [claim, fetchCustomer, fetchVehicle])

  return {
    // Data
    claim,
    customer,
    vehicle,
    claimItems,
    attachments,
    partCategories,
    parts,

    // Loading states
    claimLoading,
    customerLoading,
    vehicleLoading,
    itemsLoading,
    attachmentsLoading,

    // Refetch functions
    refetchClaim: fetchClaim,
    refetchClaimItems: fetchClaimItems,
    refetchAttachments: fetchAttachments,
  }
}

export default useClaimData
