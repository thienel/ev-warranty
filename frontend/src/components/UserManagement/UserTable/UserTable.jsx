import React, { useEffect, useMemo, useState } from 'react'
import { message, Table } from 'antd'
import api from '@services/api.js'
import { API_ENDPOINTS, ROLE_LABELS } from '@constants'
import GenerateColumns from '@components/UserManagement/UserTable/userTableColumns.jsx'

const UserTable = ({
  loading,
  setLoading,
  isReset,
  setIsReset,
  searchText,
  users,
  offices,
  handleOpenModal,
  onRefresh,
}) => {
  const [filteredInfo, setFilteredInfo] = useState({})
  const [sortedInfo, setSortedInfo] = useState({})
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  })

  useEffect(() => {
    if (isReset) {
      setFilteredInfo({})
      setSortedInfo({})
      setPagination((prev) => ({
        ...prev,
        total: users.length,
      }))
    }
  }, [isReset])

  const filteredUsers = useMemo(() => {
    if (!searchText) return users

    setIsReset(false)
    const searchLower = searchText.toLowerCase()
    return users.filter(
      (user) =>
        user.name?.toLowerCase().includes(searchLower) ||
        user.email?.toLowerCase().includes(searchLower) ||
        ROLE_LABELS[user.role]?.toLowerCase().includes(searchLower)
    )
  }, [users, searchText])

  const handleTableChange = (newPagination, filters, sorter) => {
    setPagination(newPagination)
    setFilteredInfo(filters)
    setSortedInfo(sorter)
  }
  const handleDelete = async (userId) => {
    setLoading(true)
    try {
      const response = await api.delete(`${API_ENDPOINTS.USER}${userId}`)

      if (response.data.success) {
        message.success('User deleted successfully')
        onRefresh()
      }
    } catch (error) {
      message.error(error.response?.data?.message || 'Failed to delete user')
      console.error('Error deleting user:', error)
    } finally {
      setLoading(false)
    }
  }

  const getOfficeName = (officeId) => {
    const office = offices.find((o) => o.id === officeId)
    return office ? office.name : 'N/A'
  }

  const columns = GenerateColumns(
    sortedInfo,
    filteredInfo,
    handleOpenModal,
    handleDelete,
    getOfficeName
  )

  return (
    <Table
      columns={columns}
      rowKey="id"
      loading={loading}
      dataSource={filteredUsers}
      pagination={{
        ...pagination,
        total: filteredUsers.length,
        showTotal: (total, range) => `${range[0]}-${range[1]} of ${total} users`,
        showSizeChanger: true,
        showQuickJumper: true,
        pageSizeOptions: ['10', '20', '50', '100'],
      }}
      onChange={handleTableChange}
      scroll={{ x: 1000 }}
      bordered
    />
  )
}

export default UserTable
