import React, { useEffect, useMemo, useState } from 'react'
import { message, Table } from 'antd'
import api from '@services/api.js'
import { API_ENDPOINTS, ROLE_LABELS } from '@constants'
import GenerateColumns from '@components/UserManagement/UserTable/userTableColumns.jsx'

const UserTable = ({ loading, setLoading, searchText, users, offices, onOpenModal, onRefresh }) => {
  const [filteredInfo, setFilteredInfo] = useState({})
  const [sortedInfo, setSortedInfo] = useState({})
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  })

  useEffect(() => {
    setFilteredInfo({})
    setSortedInfo({})
    setPagination((prev) => ({
      ...prev,
      total: users.length,
    }))
  }, [users])

  const filteredUsers = useMemo(() => {
    if (!searchText) return users

    const searchLower = searchText.trim().toLowerCase()
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
      await api.delete(`${API_ENDPOINTS.USER}${userId}`)
      message.success('User deleted successfully')
      onRefresh()
    } catch (error) {
      message.error(error.response?.data?.message || 'Failed to delete user')
      console.error('Error deleting user:', error)
    } finally {
      setLoading(false)
    }
  }

  const getOfficeName = (officeId) => {
    const office = offices.find((o) => o.id === officeId)
    return office ? office.office_name : 'N/A'
  }

  const columns = GenerateColumns(
    sortedInfo,
    filteredInfo,
    onOpenModal,
    handleDelete,
    getOfficeName
  )

  return (
    <Table
      size={'middle'}
      columns={columns}
      rowKey="id"
      loading={loading}
      dataSource={filteredUsers}
      showSorterTooltip={false}
      pagination={{
        ...pagination,
        total: filteredUsers.length,
        showSizeChanger: true,
        showQuickJumper: true,
        pageSizeOptions: ['10', '20', '50', '100'],
        size: 'default',
      }}
      onChange={handleTableChange}
      scroll={{ x: 1000 }}
      bordered
    />
  )
}

export default UserTable
