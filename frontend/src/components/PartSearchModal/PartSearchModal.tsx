import React, { useState, useEffect, useCallback } from 'react'
import {
  Modal,
  Input,
  Table,
  Button,
  Space,
  Typography,
  Select,
  Row,
  Col,
  message,
  Card,
  Descriptions,
} from 'antd'
import { SearchOutlined, ToolOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { partsApi, partCategoriesApi } from '@services/index'
import useHandleApiError from '@/hooks/useHandleApiError'
import type { Part, PartCategory } from '@/types/index'

const { Text } = Typography
const { Option } = Select

interface PartSearchModalProps {
  visible: boolean
  onCancel: () => void
  onSelectPart: (part: Part) => void
  title?: string
}

const PartSearchModal: React.FC<PartSearchModalProps> = ({
  visible,
  onCancel,
  onSelectPart,
  title = 'Search and Select Part',
}) => {
  const [parts, setParts] = useState<Part[]>([])
  const [filteredParts, setFilteredParts] = useState<Part[]>([])
  const [partCategories, setPartCategories] = useState<PartCategory[]>([])
  const [selectedPart, setSelectedPart] = useState<Part | null>(null)

  const [loading, setLoading] = useState(false)
  const [searchText, setSearchText] = useState('')
  const [selectedCategory, setSelectedCategory] = useState<string | undefined>(undefined)

  const handleError = useHandleApiError()

  // Fetch all parts
  const fetchParts = useCallback(async () => {
    try {
      setLoading(true)
      const response = await partsApi.getAll()
      let partsData = response.data
      if (partsData && typeof partsData === 'object' && 'data' in partsData) {
        partsData = (partsData as { data: unknown }).data as Part[]
      }

      const partsArray = Array.isArray(partsData) ? partsData : []
      setParts(partsArray)
      setFilteredParts(partsArray)
    } catch (error) {
      handleError(error as Error)
      setParts([])
      setFilteredParts([])
    } finally {
      setLoading(false)
    }
  }, [handleError])

  // Fetch part categories
  const fetchPartCategories = useCallback(async () => {
    try {
      const response = await partCategoriesApi.getAll()
      let categoriesData = response.data
      if (categoriesData && typeof categoriesData === 'object' && 'data' in categoriesData) {
        categoriesData = (categoriesData as { data: unknown }).data as PartCategory[]
      }

      const categoriesArray = Array.isArray(categoriesData) ? categoriesData : []
      setPartCategories(categoriesArray)
    } catch (error) {
      handleError(error as Error)
      setPartCategories([])
    }
  }, [handleError])

  // Filter parts based on search criteria
  const filterParts = useCallback(() => {
    let filtered = [...parts]

    // Filter by search text (serial number or part name)
    if (searchText.trim()) {
      const searchLower = searchText.toLowerCase().trim()
      filtered = filtered.filter(
        (part) =>
          part.serial_number?.toLowerCase().includes(searchLower) ||
          part.part_name?.toLowerCase().includes(searchLower),
      )
    }

    // Filter by category
    if (selectedCategory) {
      filtered = filtered.filter((part) => part.category_id === selectedCategory)
    }

    setFilteredParts(filtered)
  }, [parts, searchText, selectedCategory])

  // Load data when modal opens
  useEffect(() => {
    if (visible) {
      fetchParts()
      fetchPartCategories()
    }
  }, [visible, fetchParts, fetchPartCategories])

  // Apply filters when search criteria change
  useEffect(() => {
    filterParts()
  }, [filterParts])

  // Reset modal state when it closes
  useEffect(() => {
    if (!visible) {
      setSearchText('')
      setSelectedCategory(undefined)
      setSelectedPart(null)
    }
  }, [visible])

  const handleSelectPart = (part: Part) => {
    setSelectedPart(part)
  }

  const handleConfirmSelection = () => {
    if (selectedPart) {
      onSelectPart(selectedPart)
      onCancel() // Close modal
    } else {
      message.warning('Please select a part first')
    }
  }

  const getCategoryName = (categoryId: string): string => {
    const category = partCategories.find((cat) => cat.id === categoryId)
    return category?.category_name || 'Unknown Category'
  }

  const columns: ColumnsType<Part> = [
    {
      title: 'Serial Number',
      dataIndex: 'serial_number',
      key: 'serial_number',
      width: '20%',
      render: (text: string) => <Text strong>{text}</Text>,
    },
    {
      title: 'Part Name',
      dataIndex: 'part_name',
      key: 'part_name',
      width: '25%',
    },
    {
      title: 'Category',
      dataIndex: 'category_id',
      key: 'category_id',
      width: '20%',
      render: (categoryId: string) => getCategoryName(categoryId),
    },
    {
      title: 'Cost',
      dataIndex: 'unit_price',
      key: 'unit_price',
      width: '15%',
      align: 'right',
      render: (cost: number) => (
        <Text strong style={{ color: '#52c41a' }}>
          {cost?.toLocaleString('vi-VN', {
            style: 'currency',
            currency: 'VND',
          }) || 'N/A'}
        </Text>
      ),
    },
    {
      title: 'Action',
      key: 'action',
      width: '20%',
      align: 'center',
      render: (_, record) => (
        <Button
          type={selectedPart?.id === record.id ? 'primary' : 'default'}
          size="small"
          onClick={() => handleSelectPart(record)}
        >
          {selectedPart?.id === record.id ? 'Selected' : 'Select'}
        </Button>
      ),
    },
  ]

  return (
    <Modal
      title={
        <Space>
          <ToolOutlined />
          {title}
        </Space>
      }
      open={visible}
      onCancel={onCancel}
      width={1000}
      footer={[
        <Button key="cancel" onClick={onCancel}>
          Cancel
        </Button>,
        <Button
          key="confirm"
          type="primary"
          onClick={handleConfirmSelection}
          disabled={!selectedPart}
        >
          Confirm Selection
        </Button>,
      ]}
    >
      <Space direction="vertical" size="large" style={{ width: '100%' }}>
        {/* Search Filters */}
        <Card size="small">
          <Row gutter={16}>
            <Col span={12}>
              <Input
                placeholder="Search by serial number or part name..."
                prefix={<SearchOutlined />}
                value={searchText}
                onChange={(e) => setSearchText(e.target.value)}
                allowClear
              />
            </Col>
            <Col span={12}>
              <Select
                placeholder="Filter by category"
                style={{ width: '100%' }}
                value={selectedCategory}
                onChange={setSelectedCategory}
                allowClear
              >
                {partCategories.map((category) => (
                  <Option key={category.id} value={category.id}>
                    {category.category_name}
                  </Option>
                ))}
              </Select>
            </Col>
          </Row>
        </Card>

        {/* Selected Part Details */}
        {selectedPart && (
          <Card size="small" title="Selected Part Details">
            <Descriptions bordered column={2} size="small">
              <Descriptions.Item label="Serial Number" span={1}>
                <Text strong>{selectedPart.serial_number}</Text>
              </Descriptions.Item>
              <Descriptions.Item label="Part Name" span={1}>
                {selectedPart.part_name}
              </Descriptions.Item>
              <Descriptions.Item label="Category" span={1}>
                {getCategoryName(selectedPart.category_id)}
              </Descriptions.Item>
              <Descriptions.Item label="Cost" span={1}>
                <Text strong style={{ color: '#52c41a' }}>
                  {selectedPart.unit_price?.toLocaleString('vi-VN', {
                    style: 'currency',
                    currency: 'VND',
                  }) || 'N/A'}
                </Text>
              </Descriptions.Item>
            </Descriptions>
          </Card>
        )}

        {/* Parts Table */}
        <Table
          dataSource={filteredParts}
          columns={columns}
          rowKey="id"
          loading={loading}
          pagination={{
            pageSize: 8,
            showSizeChanger: false,
            showQuickJumper: true,
            showTotal: (total, range) => `${range[0]}-${range[1]} of ${total} parts`,
          }}
          scroll={{ y: 400 }}
          locale={{ emptyText: 'No parts found' }}
          rowSelection={{
            type: 'radio',
            selectedRowKeys: selectedPart ? [selectedPart.id] : [],
            onSelect: (record) => handleSelectPart(record),
          }}
        />
      </Space>
    </Modal>
  )
}

export default PartSearchModal
