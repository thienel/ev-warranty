import React, { useEffect, useState } from 'react'
import { Modal, Table, Typography, Tag, Spin, Alert } from 'antd'
import { SafetyOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { policyCoveragePartsApi } from '@/services/policyCoveragePartsApi'
import type { PolicyCoveragePart, PartCategory } from '@/types/index'

const { Title, Text, Paragraph } = Typography

interface PolicyCoverageModalProps {
  visible: boolean
  onCancel: () => void
  policyId: string | null
  categoryId: string | null
  categoryName: string
  partCategories: PartCategory[]
}

const PolicyCoverageModal: React.FC<PolicyCoverageModalProps> = ({
  visible,
  onCancel,
  policyId,
  categoryId,
  categoryName,
  partCategories,
}) => {
  const [coverageParts, setCoverageParts] = useState<PolicyCoveragePart[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Get category name by ID
  const getCategoryName = (catId: string): string => {
    const category = partCategories.find((cat) => cat.id === catId)
    return category?.category_name || `Category ${catId.slice(0, 8)}...`
  }

  // Fetch policy coverage parts
  const fetchCoverageParts = async () => {
    if (!policyId) return

    try {
      setLoading(true)
      setError(null)
      const response = await policyCoveragePartsApi.getAll(policyId)
      let coverageData = response.data

      // Handle nested data structure
      if (coverageData && typeof coverageData === 'object' && 'data' in coverageData) {
        coverageData = (coverageData as { data: unknown }).data as PolicyCoveragePart[]
      }

      // Filter by category if specified
      let filteredData = coverageData as PolicyCoveragePart[]
      if (categoryId) {
        filteredData = filteredData.filter((item) => item.part_category_id === categoryId)
      }

      setCoverageParts(filteredData)
    } catch (err) {
      console.error('Error fetching policy coverage parts:', err)
      setError('Failed to load policy coverage information')
      setCoverageParts([])
    } finally {
      setLoading(false)
    }
  }

  // Fetch data when modal opens
  useEffect(() => {
    if (visible && policyId) {
      fetchCoverageParts()
    }
  }, [visible, policyId, categoryId]) // eslint-disable-line react-hooks/exhaustive-deps

  // Reset state when modal closes
  useEffect(() => {
    if (!visible) {
      setCoverageParts([])
      setError(null)
    }
  }, [visible])

  const columns: ColumnsType<PolicyCoveragePart> = [
    {
      title: 'Part Category',
      dataIndex: 'part_category_id',
      key: 'part_category_id',
      width: '30%',
      render: (catId: string) => (
        <Tag color="blue">{getCategoryName(catId)}</Tag>
      ),
    },
    {
      title: 'Coverage Conditions',
      dataIndex: 'coverage_conditions',
      key: 'coverage_conditions',
      width: '70%',
      render: (conditions: string) => (
        <Paragraph 
          style={{ 
            margin: 0,
            whiteSpace: 'pre-wrap',
            wordBreak: 'break-word'
          }}
        >
          {conditions || 'No specific conditions specified'}
        </Paragraph>
      ),
    },
  ]

  const modalTitle = categoryId 
    ? `Policy Coverage for ${categoryName}`
    : 'Policy Coverage Parts'

  return (
    <Modal
      title={
        <Title level={4} style={{ margin: 0 }}>
          <SafetyOutlined /> {modalTitle}
        </Title>
      }
      open={visible}
      onCancel={onCancel}
      footer={null}
      width={800}
      destroyOnClose
    >
      {loading ? (
        <div style={{ textAlign: 'center', padding: '40px' }}>
          <Spin size="large" />
          <div style={{ marginTop: '16px' }}>
            <Text>Loading policy coverage information...</Text>
          </div>
        </div>
      ) : error ? (
        <Alert
          message="Error"
          description={error}
          type="error"
          showIcon
          style={{ marginBottom: '16px' }}
        />
      ) : coverageParts.length > 0 ? (
        <>
          {categoryId && (
            <Alert
              message={`Showing coverage conditions for ${categoryName} category`}
              type="info"
              showIcon
              style={{ marginBottom: '16px' }}
            />
          )}
          <Table
            dataSource={coverageParts}
            columns={columns}
            rowKey="id"
            pagination={false}
            locale={{ emptyText: 'No coverage parts found' }}
            size="middle"
          />
        </>
      ) : (
        <div style={{ textAlign: 'center', padding: '40px' }}>
          <Text type="secondary">
            {categoryId 
              ? `No coverage conditions found for ${categoryName} category`
              : 'No policy coverage parts found'
            }
          </Text>
        </div>
      )}
    </Modal>
  )
}

export default PolicyCoverageModal