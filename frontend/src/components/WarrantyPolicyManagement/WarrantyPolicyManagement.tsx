import React, { useState, useCallback, useEffect } from 'react'
import { Button, Space, Popconfirm, message, Input, Table } from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
  SafetyOutlined,
  SearchOutlined,
} from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useNavigate } from 'react-router-dom'
import { API_ENDPOINTS } from '@constants/common-constants'
import { type WarrantyPolicy } from '@/types/index'
import { warrantyPoliciesApi } from '@services/index'
import WarrantyPolicyModal from './WarrantyPolicyModal/WarrantyPolicyModal'
import useHandleApiError from '@/hooks/useHandleApiError'
import api from '@services/api'

const WarrantyPolicyManagement: React.FC = () => {
  const navigate = useNavigate()
  const [policies, setPolicies] = useState<WarrantyPolicy[]>([])
  const [filteredPolicies, setFilteredPolicies] = useState<WarrantyPolicy[]>([])
  const [searchText, setSearchText] = useState('')

  const [policiesLoading, setPoliciesLoading] = useState(false)
  const [modalLoading, setModalLoading] = useState(false)

  const [isPolicyModalOpen, setIsPolicyModalOpen] = useState(false)
  const [editingPolicy, setEditingPolicy] = useState<WarrantyPolicy | null>(null)
  const [isUpdate, setIsUpdate] = useState(false)

  const handleError = useHandleApiError()

  const fetchPolicies = useCallback(async () => {
    try {
      setPoliciesLoading(true)
      const response = await warrantyPoliciesApi.getAll()
      let policiesData = response.data
      if (policiesData && typeof policiesData === 'object' && 'data' in policiesData) {
        policiesData = (policiesData as { data: unknown }).data as WarrantyPolicy[]
      }
      if (Array.isArray(policiesData)) {
        setPolicies(policiesData)
        setFilteredPolicies(policiesData)
      } else {
        setPolicies([])
        setFilteredPolicies([])
      }
    } catch (error) {
      handleError(error as Error)
      setPolicies([])
      setFilteredPolicies([])
    } finally {
      setPoliciesLoading(false)
    }
  }, [handleError])

  useEffect(() => {
    fetchPolicies()
  }, [fetchPolicies])

  useEffect(() => {
    if (searchText.trim() === '') {
      setFilteredPolicies(policies)
    } else {
      const filtered = policies.filter((policy) =>
        policy.policy_name.toLowerCase().includes(searchText.toLowerCase()),
      )
      setFilteredPolicies(filtered)
    }
  }, [searchText, policies])

  const handleOpenPolicyModal = (policy?: WarrantyPolicy) => {
    if (policy) {
      setEditingPolicy(policy)
      setIsUpdate(true)
    } else {
      setEditingPolicy(null)
      setIsUpdate(false)
    }
    setIsPolicyModalOpen(true)
  }

  const handleClosePolicyModal = () => {
    setIsPolicyModalOpen(false)
    setEditingPolicy(null)
    fetchPolicies()
  }

  const handleDeletePolicy = async (id: string) => {
    try {
      await api.delete(`${API_ENDPOINTS.WARRANTY_POLICIES}/${id}`)
      message.success('Warranty policy deleted successfully')
      fetchPolicies()
    } catch (error) {
      handleError(error as Error)
    }
  }

  const handleViewDetails = (id: string) => {
    navigate(`/evm-staff/policies/${id}`)
  }

  const policyColumns: ColumnsType<WarrantyPolicy> = [
    {
      title: 'Policy Name',
      dataIndex: 'policy_name',
      key: 'policy_name',
      render: (text: string) => (
        <Space>
          <SafetyOutlined style={{ color: '#697565' }} />
          <span>{text}</span>
        </Space>
      ),
    },
    {
      title: 'Duration (Months)',
      dataIndex: 'warranty_duration_months',
      key: 'warranty_duration_months',
      align: 'center',
      width: 180,
      render: (months: number) => <span>{months}</span>,
    },
    {
      title: 'Kilometer Limit',
      dataIndex: 'kilometer_limit',
      key: 'kilometer_limit',
      align: 'right',
      width: 180,
      render: (limit?: number) => (limit ? `${limit.toLocaleString()} km` : 'N/A'),
    },
    {
      title: 'Actions',
      key: 'actions',
      align: 'center',
      width: 250,
      render: (_: unknown, record: WarrantyPolicy) => (
        <Space size="small">
          <Button
            type="default"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => handleViewDetails(record.id)}
            style={{ color: '#52c41a' }}
          >
            View Details
          </Button>
          <Button
            type="text"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleOpenPolicyModal(record)}
            style={{ color: '#1890ff' }}
          >
            Edit
          </Button>
          <Popconfirm
            title="Delete Policy"
            description="Are you sure you want to delete this policy?"
            onConfirm={() => handleDeletePolicy(record.id)}
            okText="Yes"
            cancelText="No"
          >
            <Button type="text" size="small" icon={<DeleteOutlined />} danger>
              Delete
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Input
          placeholder="Search by policy name..."
          prefix={<SearchOutlined />}
          value={searchText}
          onChange={(e) => setSearchText(e.target.value)}
          allowClear
          style={{ width: 300 }}
        />
        <Button type="primary" icon={<PlusOutlined />} onClick={() => handleOpenPolicyModal()}>
          Add Policy
        </Button>
      </div>

      <Table
        dataSource={filteredPolicies}
        columns={policyColumns}
        rowKey="id"
        loading={policiesLoading}
        pagination={{ pageSize: 10 }}
      />

      <WarrantyPolicyModal
        loading={modalLoading}
        setLoading={setModalLoading}
        onClose={handleClosePolicyModal}
        policy={editingPolicy}
        opened={isPolicyModalOpen}
        isUpdate={isUpdate}
      />
    </>
  )
}

export default WarrantyPolicyManagement
