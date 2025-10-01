import React from 'react'
import GenericTable from '@components/common/GenericTable/GenericTable.jsx'
import { API_ENDPOINTS, ROLE_LABELS } from '@constants'
import GenerateColumns from '@components/UserManagement/UserTable/userTableColumns.jsx'

const UserTable = ({ loading, setLoading, searchText, users, offices, onOpenModal, onRefresh }) => {
  const getOfficeName = (officeId) => {
    const office = offices.find((o) => o.id === officeId)
    return office ? office.office_name : 'N/A'
  }

  const searchFields = ['name', 'email']
  const searchFieldsWithRole = [...searchFields, (user) => ROLE_LABELS[user.role]]

  return (
    <GenericTable
      loading={loading}
      setLoading={setLoading}
      searchText={searchText}
      data={users}
      onOpenModal={onOpenModal}
      onRefresh={onRefresh}
      generateColumns={GenerateColumns}
      searchFields={searchFieldsWithRole}
      deleteEndpoint={API_ENDPOINTS.USER}
      deleteSuccessMessage="User deleted successfully"
      deleteErrorMessage="Failed to delete user"
      additionalProps={{ getOfficeName }}
    />
  )
}

export default UserTable
