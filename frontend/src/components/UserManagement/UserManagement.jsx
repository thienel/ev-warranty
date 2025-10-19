import React, { useEffect, useState } from 'react'
import api from '@services/api'
import { API_ENDPOINTS, ROLE_LABELS } from '@constants/common-constants.js'
import UserModal from '@components/UserManagement/UserModal/UserModal.jsx'
import useManagement from '@/hooks/useManagement.js'
import GenericActionBar from '@components/common/GenericActionBar/GenericActionBar.jsx'
import GenericTable from '@components/common/GenericTable/GenericTable.jsx'
import GenerateColumns from './userTableColumns.jsx'
import useHandleApiError from '@/hooks/useHandleApiError.js'

const UserManagement = () => {
  const {
    items: users,
    loading,
    setLoading,
    searchText,
    setSearchText,
    updateItem: updateUser,
    isUpdate,
    isOpenModal,
    handleOpenModal,
    handleReset,
  } = useManagement(API_ENDPOINTS.USER, 'user')

  const [offices, setOffices] = useState([])
  const handleError = useHandleApiError()

  const fetchOffices = async () => {
    try {
      const response = await api.get(API_ENDPOINTS.OFFICE)
      setOffices(response.data.data || [])
    } catch (error) {
      handleError(error)
    }
  }

  useEffect(() => {
    fetchOffices()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const getOfficeName = (officeId) => {
    const office = offices.find((o) => o.id === officeId)
    return office ? office.office_name : 'N/A'
  }

  const searchFields = ['name', 'email']
  const searchFieldsWithRole = [...searchFields, (user) => ROLE_LABELS[user.role]]

  return (
    <>
      <GenericActionBar
        searchText={searchText}
        setSearchText={setSearchText}
        onReset={handleReset}
        onOpenModal={handleOpenModal}
        loading={loading}
        searchPlaceholder="Search by name, email or role..."
        addButtonText="Add User"
      />

      <GenericTable
        loading={loading}
        setLoading={setLoading}
        searchText={searchText}
        data={users}
        onOpenModal={handleOpenModal}
        onRefresh={handleReset}
        generateColumns={GenerateColumns}
        searchFields={searchFieldsWithRole}
        deleteEndpoint={API_ENDPOINTS.USER}
        deleteSuccessMessage="User deleted successfully"
        deleteErrorMessage="Failed to delete user"
        additionalProps={{ getOfficeName }}
      />

      <UserModal
        loading={loading}
        setLoading={setLoading}
        onClose={handleReset}
        user={updateUser}
        opened={isOpenModal}
        offices={offices}
        isUpdate={isUpdate}
      />
    </>
  )
}

export default UserManagement
