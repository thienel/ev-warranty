import React, { useEffect, useMemo, useState } from 'react'
import { message, Table } from 'antd'
import api from '@services/api.js'
import { API_ENDPOINTS } from '@constants'
import GenerateColumns from '@components/OfficeManagement/OfficeTable/officeTableColumns.jsx'

const OfficeTable = ({ loading, setLoading, searchText, offices, onOpenModal, onRefresh }) => {
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
      total: offices.length,
    }))
  }, [offices])

  const filteredOffices = useMemo(() => {
    if (!searchText) return offices

    const searchLower = searchText.trim().toLowerCase()
    return offices.filter(
      (office) =>
        office.office_name?.toLowerCase().includes(searchLower) ||
        office.office_type?.toLowerCase().includes(searchLower) ||
        office.address?.toLowerCase().includes(searchLower)
    )
  }, [offices, searchText])

  const handleTableChange = (newPagination, filters, sorter) => {
    setPagination(newPagination)
    setFilteredInfo(filters)
    setSortedInfo(sorter)
  }

  const handleDelete = async (officeId) => {
    setLoading(true)
    try {
      await api.delete(`${API_ENDPOINTS.OFFICE}${officeId}`)
      message.success('Office deleted successfully')
      onRefresh()
    } catch (error) {
      message.error(error.response?.data?.message || 'Failed to delete office')
      console.error('Error deleting office:', error)
    } finally {
      setLoading(false)
    }
  }

  const columns = GenerateColumns(sortedInfo, filteredInfo, onOpenModal, handleDelete)

  return (
    <Table
      size={'middle'}
      columns={columns}
      rowKey="id"
      loading={loading}
      dataSource={filteredOffices}
      showSorterTooltip={false}
      pagination={{
        ...pagination,
        total: filteredOffices.length,
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

export default OfficeTable
