import React, { useState, useEffect, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import {
  Card,
  Descriptions,
  Table,
  Button,
  Space,
  Spin,
  Typography,
  Divider,
  message,
  Popconfirm,
} from 'antd'
import {
  ArrowLeftOutlined,
  EditOutlined,
  DeleteOutlined,
  PlusOutlined,
  AppstoreOutlined,
} from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import AppLayout from '@components/Layout/Layout.tsx'
import { type WarrantyPolicy, type PolicyCoveragePart, type PartCategory } from '@/types/index'
import { warrantyPoliciesApi, policyCoveragePartsApi, partCategoriesApi } from '@services/index'
import useHandleApiError from '@/hooks/useHandleApiError'
import api from '@services/api'
import { API_ENDPOINTS } from '@constants/common-constants'
import CoveredPartModal from '@components/WarrantyPolicyManagement/CoveredPartModal/CoveredPartModal'

const { Title, Paragraph } = Typography

const PolicyDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const handleError = useHandleApiError()

  const [policy, setPolicy] = useState<WarrantyPolicy | null>(null)
  const [coveredParts, setCoveredParts] = useState<PolicyCoveragePart[]>([])
  const [partCategories, setPartCategories] = useState<PartCategory[]>([])

  const [policyLoading, setPolicyLoading] = useState(false)
  const [coveredPartsLoading, setCoveredPartsLoading] = useState(false)
  const [partCategoriesLoading, setPartCategoriesLoading] = useState(false)
  const [modalLoading, setModalLoading] = useState(false)

  const [isCoveredPartModalOpen, setIsCoveredPartModalOpen] = useState(false)
  const [editingCoveredPart, setEditingCoveredPart] = useState<PolicyCoveragePart | null>(null)
  const [isUpdate, setIsUpdate] = useState(false)

  const fetchPolicy = useCallback(async () => {
    if (!id) return

    try {
      setPolicyLoading(true)
      const response = await warrantyPoliciesApi.getById(id)
      let policyData = response.data
      if (policyData && typeof policyData === 'object' && 'data' in policyData) {
        policyData = (policyData as { data: unknown }).data as WarrantyPolicy
      }
      setPolicy(policyData as WarrantyPolicy)
    } catch (error) {
      handleError(error as Error)
      message.error('Failed to load policy details')
    } finally {
      setPolicyLoading(false)
    }
  }, [id, handleError])

  const fetchCoveredParts = useCallback(async () => {
    if (!id) return

    try {
      setCoveredPartsLoading(true)
      const response = await policyCoveragePartsApi.getAll(id)
      let partsData = response.data
      if (partsData && typeof partsData === 'object' && 'data' in partsData) {
        partsData = (partsData as { data: unknown }).data as PolicyCoveragePart[]
      }
      if (Array.isArray(partsData)) {
        setCoveredParts(partsData)
      } else {
        setCoveredParts([])
      }
    } catch (error) {
      handleError(error as Error)
      setCoveredParts([])
    } finally {
      setCoveredPartsLoading(false)
    }
  }, [id, handleError])

  const fetchPartCategories = useCallback(async () => {
    try {
      setPartCategoriesLoading(true)
      const response = await partCategoriesApi.getAll()
      let categoriesData = response.data
      if (categoriesData && typeof categoriesData === 'object' && 'data' in categoriesData) {
        categoriesData = (categoriesData as { data: unknown }).data as PartCategory[]
      }
      if (Array.isArray(categoriesData)) {
        setPartCategories(categoriesData)
      } else {
        setPartCategories([])
      }
    } catch (error) {
      handleError(error as Error)
      setPartCategories([])
    } finally {
      setPartCategoriesLoading(false)
    }
  }, [handleError])

  useEffect(() => {
    fetchPolicy()
    fetchCoveredParts()
    fetchPartCategories()
  }, [fetchPolicy, fetchCoveredParts, fetchPartCategories])

  const handleOpenCoveredPartModal = (coveredPart?: PolicyCoveragePart) => {
    if (coveredPart) {
      setEditingCoveredPart(coveredPart)
      setIsUpdate(true)
    } else {
      setEditingCoveredPart(null)
      setIsUpdate(false)
    }
    setIsCoveredPartModalOpen(true)
  }

  const handleCloseCoveredPartModal = () => {
    setIsCoveredPartModalOpen(false)
    setEditingCoveredPart(null)
    fetchCoveredParts()
  }

  const handleDeleteCoveredPart = async (partId: string) => {
    try {
      await api.delete(`${API_ENDPOINTS.POLICY_COVERAGE_PARTS}/${partId}`)
      message.success('Covered part removed successfully')
      fetchCoveredParts()
    } catch (error) {
      handleError(error as Error)
    }
  }

  const getCategoryName = (categoryId: string): string => {
    const category = partCategories.find((cat) => cat.id === categoryId)
    return category ? category.category_name : 'N/A'
  }

  const coveredPartsColumns: ColumnsType<PolicyCoveragePart> = [
    {
      title: 'Part Category',
      dataIndex: 'part_category_id',
      key: 'part_category_id',
      render: (categoryId: string) => (
        <Space>
          <AppstoreOutlined style={{ color: '#697565' }} />
          <span>{getCategoryName(categoryId)}</span>
        </Space>
      ),
    },
    {
      title: 'Coverage Conditions',
      dataIndex: 'coverage_conditions',
      key: 'coverage_conditions',
      ellipsis: true,
      render: (text: string) => <span>{text || 'N/A'}</span>,
    },
    {
      title: 'Actions',
      key: 'actions',
      align: 'center',
      width: 150,
      render: (_: unknown, record: PolicyCoveragePart) => (
        <Space size="small">
          <Button
            type="text"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleOpenCoveredPartModal(record)}
            style={{ color: '#1890ff' }}
          >
            Edit
          </Button>
          <Popconfirm
            title="Remove Covered Part"
            description="Are you sure you want to remove this covered part?"
            onConfirm={() => handleDeleteCoveredPart(record.id)}
            okText="Yes"
            cancelText="No"
          >
            <Button type="text" size="small" icon={<DeleteOutlined />} danger>
              Remove
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  if (policyLoading) {
    return (
      <AppLayout title="Policy Details">
        <div style={{ textAlign: 'center', padding: '50px' }}>
          <Spin size="large" />
        </div>
      </AppLayout>
    )
  }

  if (!policy) {
    return (
      <AppLayout title="Policy Details">
        <Card>
          <div style={{ textAlign: 'center', padding: '50px' }}>
            <Title level={4}>Policy not found</Title>
            <Button type="primary" onClick={() => navigate('/evm-staff/policies')}>
              Back to Policies
            </Button>
          </div>
        </Card>
      </AppLayout>
    )
  }

  return (
    <AppLayout title="Policy Details">
      <div style={{ marginBottom: 16 }}>
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/evm-staff/policies')}>
          Back to Policies
        </Button>
      </div>

      <Card
        title={<Title level={4}>Warranty Policy Information</Title>}
        style={{ marginBottom: 24 }}
      >
        <Descriptions column={2} bordered>
          <Descriptions.Item label="Policy Name" span={2}>
            {policy.policy_name}
          </Descriptions.Item>
          <Descriptions.Item label="Warranty Duration">
            {policy.warranty_duration_months} months
          </Descriptions.Item>
          <Descriptions.Item label="Kilometer Limit">
            {policy.kilometer_limit ? `${policy.kilometer_limit.toLocaleString()} km` : 'N/A'}
          </Descriptions.Item>
          <Descriptions.Item label="Created At">
            {new Date(policy.created_at).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="Updated At">
            {policy.updated_at ? new Date(policy.updated_at).toLocaleString() : 'N/A'}
          </Descriptions.Item>
        </Descriptions>

        <Divider />

        <div>
          <Title level={5}>Terms and Conditions</Title>
          <Paragraph style={{ whiteSpace: 'pre-wrap' }}>{policy.terms_and_conditions}</Paragraph>
        </div>
      </Card>

      <Card
        title={<Title level={4}>Covered Parts</Title>}
        extra={
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => handleOpenCoveredPartModal()}
          >
            Add Covered Part
          </Button>
        }
      >
        <Table
          dataSource={coveredParts}
          columns={coveredPartsColumns}
          rowKey="id"
          loading={coveredPartsLoading}
          pagination={{ pageSize: 10 }}
        />
      </Card>

      {id && (
        <CoveredPartModal
          loading={modalLoading}
          setLoading={setModalLoading}
          onClose={handleCloseCoveredPartModal}
          coveragePart={editingCoveredPart}
          opened={isCoveredPartModalOpen}
          isUpdate={isUpdate}
          policyId={id}
          partCategories={partCategories}
          partCategoriesLoading={partCategoriesLoading}
        />
      )}
    </AppLayout>
  )
}

export default PolicyDetail
