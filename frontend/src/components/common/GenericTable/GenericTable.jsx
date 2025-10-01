import React, { useEffect, useMemo, useState } from 'react'
import { message, Table } from 'antd'
import api from '@services/api.js'

const GenericTable = ({
  loading,
  setLoading,
  searchText,
  data,
  onOpenModal,
  onRefresh,
  generateColumns,
  searchFields = [],
  deleteEndpoint,
  deleteSuccessMessage = 'Item deleted successfully',
  deleteErrorMessage = 'Failed to delete item',
  additionalProps = {}, // For passing extra props like offices to columns
}) => {
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
      total: data.length,
    }))
  }, [data])

  const filteredData = useMemo(() => {
    if (!searchText) return data

    const searchLower = searchText.trim().toLowerCase()
    return data.filter((item) =>
      searchFields.some((field) => {
        let value
        if (typeof field === 'function') {
          value = field(item)
        } else if (field.includes('.')) {
          value = field.split('.').reduce((obj, key) => obj?.[key], item)
        } else {
          value = item[field]
        }
        return value?.toString().toLowerCase().includes(searchLower)
      })
    )
  }, [data, searchText, searchFields])

  const handleTableChange = (newPagination, filters, sorter) => {
    setPagination(newPagination)
    setFilteredInfo(filters)
    setSortedInfo(sorter)
  }

  const handleDelete = async (itemId) => {
    setLoading(true)
    try {
      await api.delete(`${deleteEndpoint}${itemId}`)
      message.success(deleteSuccessMessage)
      onRefresh()
    } catch (error) {
      message.error(error.response?.data?.message || deleteErrorMessage)
      console.error('Error deleting item:', error)
    } finally {
      setLoading(false)
    }
  }

  const columns = generateColumns(
    sortedInfo,
    filteredInfo,
    onOpenModal,
    handleDelete,
    additionalProps
  )

  return (
    <Table
      size={'middle'}
      columns={columns}
      rowKey="id"
      loading={loading}
      dataSource={filteredData}
      showSorterTooltip={false}
      pagination={{
        ...pagination,
        total: filteredData.length,
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

export default GenericTable
