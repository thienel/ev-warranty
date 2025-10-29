import React from 'react'
import { API_ENDPOINTS } from '@constants/common-constants'
import { type VehicleModel } from '@/types/index'
import VehicleModelModal from '@components/VehicleModelManagement/VehicleModelModal/VehicleModelModal'
import useManagement from '@/hooks/useManagement'
import GenericActionBar from '@components/common/GenericActionBar/GenericActionBar'
import GenericTable from '@components/common/GenericTable/GenericTable'
import GenerateColumns from './vehicleModelTableColumns'

const VehicleModelManagement: React.FC = () => {
  const {
    items: vehicleModels,
    loading,
    setLoading,
    searchText,
    setSearchText,
    updateItem: updateVehicleModel,
    isUpdate,
    isOpenModal,
    handleOpenModal,
    handleReset,
  } = useManagement(API_ENDPOINTS.VEHICLE_MODELS)

  const searchFields = ['brand', 'model_name']
  const searchFieldsWithYear = [
    ...searchFields,
    (vehicleModel: Record<string, unknown> & { id: string | number }) => {
      const modelRecord = vehicleModel as unknown as VehicleModel
      return modelRecord.year?.toString() || ''
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
        searchPlaceholder="Search by brand, model or year..."
        addButtonText="Add Model"
        allowCreate={true}
      />

      <GenericTable
        loading={loading}
        setLoading={setLoading}
        searchText={searchText}
        data={vehicleModels as (Record<string, unknown> & { id: string | number })[]}
        onOpenModal={handleOpenModal}
        onRefresh={handleReset}
        generateColumns={GenerateColumns}
        searchFields={searchFieldsWithYear}
        deleteEndpoint={API_ENDPOINTS.VEHICLE_MODELS}
        deleteSuccessMessage="Vehicle model deleted successfully"
      />

      <VehicleModelModal
        loading={loading}
        setLoading={setLoading}
        onClose={handleReset}
        vehicleModel={updateVehicleModel ? (updateVehicleModel as unknown as VehicleModel) : null}
        opened={isOpenModal}
        isUpdate={isUpdate}
      />
    </>
  )
}

export default VehicleModelManagement
