import React, { useCallback, useEffect, useState } from 'react'
import { API_ENDPOINTS } from '@constants/common-constants'
import { type Vehicle, type Customer, type VehicleModel } from '@/types/index'
import { customersApi, vehicleModelsApi } from '@services/index'
import VehicleModal from '@components/VehicleManagement/VehicleModal/VehicleModal'
import useManagement from '@/hooks/useManagement'
import GenericActionBar from '@components/common/GenericActionBar/GenericActionBar'
import GenericTable from '@components/common/GenericTable/GenericTable'
import GenerateColumns from './vehicleTableColumns'
import useHandleApiError from '@/hooks/useHandleApiError'

const VehicleManagement: React.FC = () => {
  const {
    items: vehicles,
    loading,
    setLoading,
    searchText,
    setSearchText,
    updateItem: updateVehicle,
    isUpdate,
    isOpenModal,
    handleOpenModal,
    handleReset,
  } = useManagement(API_ENDPOINTS.VEHICLES)

  const [customers, setCustomers] = useState<Customer[]>([])
  const [vehicleModels, setVehicleModels] = useState<VehicleModel[]>([])
  const [customersLoading, setCustomersLoading] = useState(false)
  const [vehicleModelsLoading, setVehicleModelsLoading] = useState(false)
  const handleError = useHandleApiError()

  const fetchCustomers = useCallback(async (): Promise<void> => {
    try {
      setCustomersLoading(true)
      const response = await customersApi.getAll()
      // Handle different response structures
      let customersData = response.data
      if (customersData && typeof customersData === 'object' && 'data' in customersData) {
        customersData = (customersData as { data: unknown }).data as Customer[]
      }
      if (Array.isArray(customersData)) {
        setCustomers(customersData)
      } else {
        setCustomers([])
      }
    } catch (error) {
      handleError(error as Error)
      setCustomers([])
    } finally {
      setCustomersLoading(false)
    }
  }, [handleError])

  const fetchVehicleModels = useCallback(async (): Promise<void> => {
    try {
      setVehicleModelsLoading(true)
      const response = await vehicleModelsApi.getAll()
      // Handle different response structures
      let vehicleModelsData = response.data
      if (
        vehicleModelsData &&
        typeof vehicleModelsData === 'object' &&
        'data' in vehicleModelsData
      ) {
        vehicleModelsData = (vehicleModelsData as { data: unknown }).data as VehicleModel[]
      }
      if (Array.isArray(vehicleModelsData)) {
        setVehicleModels(vehicleModelsData)
      } else {
        setVehicleModels([])
      }
    } catch (error) {
      handleError(error as Error)
      setVehicleModels([])
    } finally {
      setVehicleModelsLoading(false)
    }
  }, [handleError])

  useEffect(() => {
    fetchCustomers()
    fetchVehicleModels()
  }, [fetchCustomers, fetchVehicleModels])

  // Refetch data when modal opens if lists are empty
  useEffect(() => {
    if (isOpenModal) {
      if (customers.length === 0 && !customersLoading) {
        fetchCustomers()
      }
      if (vehicleModels.length === 0 && !vehicleModelsLoading) {
        fetchVehicleModels()
      }
    }
  }, [
    isOpenModal,
    customers.length,
    vehicleModels.length,
    customersLoading,
    vehicleModelsLoading,
    fetchCustomers,
    fetchVehicleModels,
  ])

  const getCustomerName = (customerId: string): string => {
    if (!Array.isArray(customers)) {
      return 'N/A'
    }
    const customer = customers.find((c) => c.id === customerId)
    return customer ? `${customer.first_name || ''} ${customer.last_name || ''}`.trim() : 'N/A'
  }

  const getVehicleModelName = (modelId: string): string => {
    if (!Array.isArray(vehicleModels)) {
      return 'N/A'
    }
    const model = vehicleModels.find((m) => m.id === modelId)
    return model ? `${model.brand} ${model.model_name} (${model.year})` : 'N/A'
  }

  const searchFields = ['vin', 'license_plate']

  return (
    <>
      <GenericActionBar
        searchText={searchText}
        setSearchText={setSearchText}
        onReset={handleReset}
        onOpenModal={handleOpenModal}
        loading={loading}
        searchPlaceholder="Search by VIN or license plate..."
        addButtonText="Add Vehicle"
        allowCreate={true}
      />

      <GenericTable
        loading={loading}
        setLoading={setLoading}
        searchText={searchText}
        data={vehicles as (Record<string, unknown> & { id: string | number })[]}
        onOpenModal={handleOpenModal}
        onRefresh={handleReset}
        generateColumns={GenerateColumns}
        searchFields={searchFields}
        deleteEndpoint={API_ENDPOINTS.VEHICLES}
        deleteSuccessMessage="Vehicle deleted successfully"
        additionalProps={{ getCustomerName, getVehicleModelName }}
      />

      <VehicleModal
        loading={loading}
        setLoading={setLoading}
        onClose={handleReset}
        vehicle={updateVehicle ? (updateVehicle as unknown as Vehicle) : null}
        opened={isOpenModal}
        customers={customers}
        vehicleModels={vehicleModels}
        customersLoading={customersLoading}
        vehicleModelsLoading={vehicleModelsLoading}
        isUpdate={isUpdate}
      />
    </>
  )
}

export default VehicleManagement
