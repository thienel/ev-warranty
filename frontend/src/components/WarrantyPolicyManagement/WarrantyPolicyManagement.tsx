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
import { type WarrantyPolicy, type VehicleModel } from '@/types/index'
import { warrantyPoliciesApi, vehicleModelsApi } from '@services/index'
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

      // Fetch both policies and vehicle models in parallel
      const [policiesResponse, vehicleModelsResponse] = await Promise.all([
        warrantyPoliciesApi.getAll(),
        vehicleModelsApi.getAll(),
      ])

      let policiesData = policiesResponse.data
      if (policiesData && typeof policiesData === 'object' && 'data' in policiesData) {
        policiesData = (policiesData as { data: unknown }).data as WarrantyPolicy[]
      }

      let vehicleModelsData = vehicleModelsResponse.data
      if (
        vehicleModelsData &&
        typeof vehicleModelsData === 'object' &&
        'data' in vehicleModelsData
      ) {
        vehicleModelsData = (vehicleModelsData as { data: unknown }).data as VehicleModel[]
      }

      if (Array.isArray(policiesData) && Array.isArray(vehicleModelsData)) {
        // Map vehicle models to their policies
        const policiesWithModels = policiesData.map((policy) => {
          const relatedModels = vehicleModelsData
            .filter((model: VehicleModel & { policy_id?: string }) => model.policy_id === policy.id)
            .map((model: VehicleModel) => ({
              id: model.id,
              brand: model.brand,
              model_name: model.model_name,
              year: model.year,
            }))

          return {
            ...policy,
            vehicle_models: relatedModels,
          }
        })

        setPolicies(policiesWithModels)
        setFilteredPolicies(policiesWithModels)
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
      const searchLower = searchText.toLowerCase()
      const filtered = policies.filter((policy) => {
        // Search by policy name
        if (policy.policy_name.toLowerCase().includes(searchLower)) {
          return true
        }
        // Search by vehicle model names
        if (policy.vehicle_models && policy.vehicle_models.length > 0) {
          return policy.vehicle_models.some(
            (model) =>
              model.brand.toLowerCase().includes(searchLower) ||
              model.model_name.toLowerCase().includes(searchLower) ||
              `${model.brand} ${model.model_name}`.toLowerCase().includes(searchLower),
          )
        }
        return false
      })
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
      width: '25%',
      render: (text: string) => (
        <Space>
          <SafetyOutlined style={{ color: '#697565' }} />
          <span>{text}</span>
        </Space>
      ),
    },
    {
      title: 'Vehicle Models',
      dataIndex: 'vehicle_models',
      key: 'vehicle_models',
      width: '25%',
      render: (_: unknown, record: WarrantyPolicy) => {
        if (!record.vehicle_models || record.vehicle_models.length === 0) {
          return <span style={{ color: '#999' }}>No models assigned</span>
        }
        const modelNames = record.vehicle_models
          .map((model) => `${model.brand} ${model.model_name} (${model.year})`)
          .join(', ')
        return (
          <span
            style={{
              whiteSpace: 'normal',
              wordBreak: 'break-word',
            }}
          >
            {modelNames}
          </span>
        )
      },
    },
    {
      title: 'Duration (Months)',
      dataIndex: 'warranty_duration_months',
      key: 'warranty_duration_months',
      align: 'center',
      width: '12%',
      render: (months: number) => <span>{months}</span>,
    },
    {
      title: 'Kilometer Limit',
      dataIndex: 'kilometer_limit',
      key: 'kilometer_limit',
      align: 'right',
      width: '12%',
      render: (limit?: number) => (limit ? `${limit.toLocaleString()} km` : 'N/A'),
    },
    {
      title: 'Details',
      key: 'details',
      align: 'center',
      width: '8%',
      render: (_: unknown, record: WarrantyPolicy) => (
        <Space size="middle">
          <Button
            type="text"
            icon={<EyeOutlined />}
            onClick={() => handleViewDetails(record.id)}
            title="View Details"
            style={{ color: '#1890ff' }}
          />
        </Space>
      ),
    },
    {
      title: 'Actions',
      key: 'actions',
      align: 'center',
      width: '18%',
      render: (_: unknown, record: WarrantyPolicy) => (
        <Space size="small">
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
      <div
        style={{
          marginBottom: 16,
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        <Input
          placeholder="Search by policy name or vehicle model..."
          prefix={<SearchOutlined />}
          value={searchText}
          onChange={(e) => setSearchText(e.target.value)}
          allowClear
          style={{ width: 350 }}
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
