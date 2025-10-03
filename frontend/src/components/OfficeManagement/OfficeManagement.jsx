import React from 'react'
import { API_ENDPOINTS, USER_ROLES } from '@constants'
import OfficeModal from '@components/OfficeManagement/OfficeModal/OfficeModal.jsx'
import useManagement from '@/hooks/useManagement.js'
import GenericActionBar from '@components/common/GenericActionBar/GenericActionBar.jsx'
import GenericTable from '@components/common/GenericTable/GenericTable.jsx'
import GenerateColumns from './officeTableColumns.jsx'
import useCheckRole from '@/hooks/useCheckRole.js'

const OfficeManagement = () => {
  useCheckRole([USER_ROLES.ADMIN])

  const {
    items: offices,
    loading,
    setLoading,
    searchText,
    setSearchText,
    updateItem: updateOffice,
    isUpdate,
    isOpenModal,
    handleOpenModal,
    handleReset,
  } = useManagement(API_ENDPOINTS.OFFICE, 'office')

  const searchFields = ['office_name', 'office_type', 'address']

  return (
    <>
      <GenericActionBar
        searchText={searchText}
        setSearchText={setSearchText}
        onReset={handleReset}
        onOpenModal={handleOpenModal}
        loading={loading}
        searchPlaceholder="Search by office name, type or address..."
        addButtonText="Add Office"
      />

      <GenericTable
        loading={loading}
        setLoading={setLoading}
        searchText={searchText}
        data={offices}
        onOpenModal={handleOpenModal}
        onRefresh={handleReset}
        generateColumns={GenerateColumns}
        searchFields={searchFields}
        deleteEndpoint={API_ENDPOINTS.OFFICE}
        deleteSuccessMessage="Office deleted successfully"
        deleteErrorMessage="Failed to delete office"
      />

      <OfficeModal
        loading={loading}
        setLoading={setLoading}
        onClose={handleReset}
        office={updateOffice}
        opened={isOpenModal}
        isUpdate={isUpdate}
      />
    </>
  )
}

export default OfficeManagement
