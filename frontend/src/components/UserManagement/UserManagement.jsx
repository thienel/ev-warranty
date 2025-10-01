import React, { useEffect, useState } from 'react'
import { message } from 'antd'
import api from '@services/api'
import { API_ENDPOINTS } from '@constants'
import UserModal from '@components/UserManagement/UserModal/UserModal.jsx'
import UserTable from '@components/UserManagement/UserTable/UserTable.jsx'
import UserActionBar from '@components/UserManagement/UserActionBar/UserActionBar.jsx'
import useManagement from '@/hooks/useManagement.js'

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

  const fetchOffices = async () => {
    try {
      const response = await api.get(API_ENDPOINTS.OFFICE)
      if (response.data.success) {
        setOffices(response.data.data || [])
      }
    } catch (error) {
      message.error(error.response?.data?.message || 'Failed to load offices')
      console.error('Error fetching offices:', error)
    }
  }

  useEffect(() => {
    fetchOffices()
  }, [])

  return (
    <>
      <UserActionBar
        searchText={searchText}
        setSearchText={setSearchText}
        onReset={handleReset}
        onOpenModal={handleOpenModal}
        loading={loading}
      />

      <UserTable
        loading={loading}
        setLoading={setLoading}
        users={users}
        offices={offices}
        searchText={searchText}
        onOpenModal={handleOpenModal}
        onRefresh={handleReset}
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
