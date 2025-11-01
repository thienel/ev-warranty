import { useState } from 'react'
import { message } from 'antd'
import { useLocation, useNavigate } from 'react-router-dom'
import { claimsApi } from '@services/claimsApi'
import type { Customer, Vehicle, CreateClaimRequest } from '@/types'
import useHandleApiError from '@/hooks/useHandleApiError'
import { getClaimsBasePath } from '@/utils/navigationHelpers'

interface ClaimFormValues {
  description: string
}

export const useClaimSubmission = () => {
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()
  const handleError = useHandleApiError()
  const location = useLocation()

  const submitClaim = async (
    values: ClaimFormValues,
    selectedCustomer: Customer | null,
    selectedVehicle: Vehicle | null,
  ) => {
    if (!selectedCustomer) {
      message.error('Please select a customer')
      return { success: false, shouldNavigateToStep: 0 }
    }

    if (!selectedVehicle) {
      message.error('Please select a vehicle')
      return { success: false, shouldNavigateToStep: 1 }
    }

    if (!values.description?.trim()) {
      message.error('Please provide a claim description')
      return { success: false, shouldNavigateToStep: 2 }
    }

    try {
      setLoading(true)

      const claimData: CreateClaimRequest = {
        customer_id: selectedCustomer.id,
        vehicle_id: selectedVehicle.id,
        description: values.description.trim(),
      }

      await claimsApi.create(claimData)

      message.success('Claim created successfully!')
      navigate(getClaimsBasePath(location.pathname))
      return { success: true }
    } catch (error) {
      console.error('Failed to create claim:', error)
      handleError(error as Error)
      return { success: false }
    } finally {
      setLoading(false)
    }
  }

  return {
    loading,
    submitClaim,
  }
}
