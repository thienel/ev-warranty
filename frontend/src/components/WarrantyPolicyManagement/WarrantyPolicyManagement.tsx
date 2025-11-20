import React, { useCallback, useEffect, useState } from 'react'
import { API_ENDPOINTS } from '@constants/common-constants'
import { type WarrantyPolicy, type VehicleModel } from '@/types/index'
import { warrantyPoliciesApi, vehicleModelsApi } from '@services/index'
import WarrantyPolicyModal from './WarrantyPolicyModal/WarrantyPolicyModal'
import useHandleApiError from '@/hooks/useHandleApiError'
import useManagement from '@/hooks/useManagement'
import GenericActionBar from '@components/common/GenericActionBar/GenericActionBar'
import GenericTable from '@components/common/GenericTable/GenericTable'
import GenerateColumns from './warrantyPolicyTableColumns'

const WarrantyPolicyManagement: React.FC = () => {
  const {
    items: policies,
    setItems: setPolicies,
    loading,
    setLoading,
    searchText,
    setSearchText,
    updateItem: updatePolicy,
    isUpdate,
    isOpenModal,
    setIsOpenModal,
    setUpdateItem,
    handleOpenModal,
  } = useManagement<WarrantyPolicy>(API_ENDPOINTS.WARRANTY_POLICIES)

  const [vehicleModels, setVehicleModels] = useState<VehicleModel[]>([])
  const handleError = useHandleApiError()

  // Fetch vehicle models and enrich policies with vehicle model data
  const enrichPoliciesWithVehicleModels = useCallback(async () => {
    try {
      // Fetch both policies and vehicle models
      const [policiesResponse, vehicleModelsResponse] = await Promise.all([
        warrantyPoliciesApi.getAll(),
        vehicleModelsApi.getAll(),
      ])

      let policiesData = policiesResponse.data
      if (policiesData && typeof policiesData === 'object' && 'data' in policiesData) {
        policiesData = (policiesData as { data: unknown }).data as WarrantyPolicy[]
      }

      let vehicleModelsData = vehicleModelsResponse.data
      if (
        vehicleModelsData &&
        typeof vehicleModelsData === 'object' &&
        'data' in vehicleModelsData
      ) {
        vehicleModelsData = (vehicleModelsData as { data: unknown }).data as VehicleModel[]
      }

      if (Array.isArray(policiesData) && Array.isArray(vehicleModelsData)) {
        // Store vehicle models for the modal
        setVehicleModels(vehicleModelsData)

        // Map vehicle models to their policies
        const policiesWithModels = policiesData.map((policy) => {
          const relatedModels = vehicleModelsData
            .filter((model: VehicleModel & { policy_id?: string }) => model.policy_id === policy.id)
            .map((model: VehicleModel) => ({
              id: model.id,
              brand: model.brand,
              model_name: model.model_name,
              year: model.year,
            }))

          return {
            ...policy,
            vehicle_models: relatedModels,
          }
        })

        setPolicies(policiesWithModels)
      } else {
        setPolicies([])
        setVehicleModels([])
      }
    } catch (error) {
      handleError(error as Error)
      setPolicies([])
      setVehicleModels([])
    }
  }, [handleError, setPolicies])

  const handleReset = useCallback(async () => {
    setLoading(true)
    setSearchText('')
    setIsOpenModal(false)
    setUpdateItem(null)
    await enrichPoliciesWithVehicleModels()
    setLoading(false)
  }, [setLoading, setSearchText, setIsOpenModal, setUpdateItem, enrichPoliciesWithVehicleModels])

  useEffect(() => {
    enrichPoliciesWithVehicleModels()
  }, [enrichPoliciesWithVehicleModels])

  const searchFields = [
    'policy_name',
    (policy: Record<string, unknown> & { id: string | number }) => {
      const policyRecord = policy as unknown as WarrantyPolicy
      if (policyRecord.vehicle_models && policyRecord.vehicle_models.length > 0) {
        return policyRecord.vehicle_models
          .map((model) => `${model.brand} ${model.model_name}`)
          .join(' ')
      }
      return ''
    },
  ]

  return (
    <>
      <GenericActionBar
        searchText={searchText}
        setSearchText={setSearchText}
        onReset={handleReset}
        onOpenModal={handleOpenModal}
        loading={loading}
        searchPlaceholder="Search by policy name or vehicle model..."
        addButtonText="Add Policy"
        allowCreate={true}
      />

      <GenericTable
        loading={loading}
        setLoading={setLoading}
        searchText={searchText}
        data={policies as unknown as (Record<string, unknown> & { id: string | number })[]}
        onOpenModal={
          handleOpenModal as (
            item?: (Record<string, unknown> & { id: string | number }) | null,
            isUpdate?: boolean,
          ) => void
        }
        onRefresh={enrichPoliciesWithVehicleModels}
        generateColumns={GenerateColumns}
        searchFields={searchFields}
        deleteEndpoint={API_ENDPOINTS.WARRANTY_POLICIES}
        deleteSuccessMessage="Warranty policy deleted successfully"
      />

      <WarrantyPolicyModal
        loading={loading}
        setLoading={setLoading}
        onClose={handleReset}
        policy={updatePolicy ? (updatePolicy as unknown as WarrantyPolicy) : null}
        opened={isOpenModal}
        isUpdate={isUpdate}
        vehicleModels={vehicleModels}
      />
    </>
  )
}

export default WarrantyPolicyManagement
