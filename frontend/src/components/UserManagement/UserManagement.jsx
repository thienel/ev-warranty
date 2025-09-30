import React, { useState, useEffect } from 'react'
import { message } from 'antd'
import api from '@services/api'
import { API_ENDPOINTS } from '@constants'
import UserModal from '@components/UserManagement/UserModal/UserModal.jsx'
import UserTable from '@components/UserManagement/UserTable/UserTable.jsx'
import UserActionBar from '@components/UserManagement/UserActionBar/UserActionBar.jsx'

const UserManagement = () => {
  const [users, setUsers] = useState([])
  const [offices, setOffices] = useState([])

  const [loading, setLoading] = useState(false)
  const [searchText, setSearchText] = useState('')

  const [updateUser, setUpdateUser] = useState(null)
  const [isOpenModal, setIsOpenModal] = useState(false)

  const [isResetTable, setIsResetTable] = useState(true)

  const handleOpenModal = (user = null) => {
    setUpdateUser(user)
    setIsOpenModal(true)
  }

  const handleCloseModal = () => {
    setIsOpenModal(false)
    setUpdateUser(null)
    fetchUsers()
  }

  const fetchUsers = async () => {
    setLoading(true)
    try {
      const response = await api.get(API_ENDPOINTS.USER)

      if (response.data.success) {
        const userData = response.data.data || []
        setUsers(userData)
        setIsResetTable(true)
      }
    } catch (error) {
      message.error(error.response?.data?.message || 'Failed to load users')
      console.error('Error fetching users:', error)
    } finally {
      setLoading(false)
    }
  }

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
    fetchUsers()
    fetchOffices()
  }, [])

  const handleReset = () => {
    setSearchText('')
    setIsResetTable(true)
  }

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
        isReset={isResetTable}
        setIsReset={setIsResetTable}
        users={users}
        offices={offices}
        handleOpenModal={handleOpenModal}
        onRefresh={handleReset}
      />

      <UserModal
        loading={loading}
        setLoading={setLoading}
        onClose={handleCloseModal}
        user={updateUser}
        opened={isOpenModal}
        offices={offices}
      />
    </>
  )
}

export default UserManagement
