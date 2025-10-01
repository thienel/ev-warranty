import React from 'react'
import GenericTable from '@components/common/GenericTable/GenericTable.jsx'
import { API_ENDPOINTS } from '@constants'
import GenerateColumns from '@components/OfficeManagement/OfficeTable/officeTableColumns.jsx'

const OfficeTable = ({ loading, setLoading, searchText, offices, onOpenModal, onRefresh }) => {
  const searchFields = ['office_name', 'office_type', 'address']

  return (
    <GenericTable
      loading={loading}
      setLoading={setLoading}
      searchText={searchText}
      data={offices}
      onOpenModal={onOpenModal}
      onRefresh={onRefresh}
      generateColumns={GenerateColumns}
      searchFields={searchFields}
      deleteEndpoint={API_ENDPOINTS.OFFICE}
      deleteSuccessMessage="Office deleted successfully"
      deleteErrorMessage="Failed to delete office"
      additionalProps={{}}
    />
  )
}

export default OfficeTable
