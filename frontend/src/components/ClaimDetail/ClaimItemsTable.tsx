import React from 'react'
import { Table, Button, Typography, Divider, Card } from 'antd'
import { ToolOutlined, PlusOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { CLAIM_ITEM_STATUS_LABELS, CLAIM_ITEM_TYPE_LABELS } from '@constants/common-constants'
import type { ClaimItem, PartCategory, Part } from '@/types/index'

const { Title, Text } = Typography

interface ClaimItemsTableProps {
  claimItems: ClaimItem[]
  partCategories: PartCategory[]
  parts: Part[]
  loading: boolean
  canAddItems: boolean
  onAddItem: () => void
}

const ClaimItemsTable: React.FC<ClaimItemsTableProps> = ({
  claimItems,
  partCategories,
  parts,
  loading,
  canAddItems,
  onAddItem,
}) => {
  // Get part category name by ID
  const getPartCategoryName = (categoryId: string): string => {
    const category = partCategories.find((cat) => cat.id === categoryId)
    return category?.category_name || `Category ${categoryId.slice(0, 8)}...`
  }

  // Get part serial number by ID
  const getPartSerialNumber = (partId: string): string => {
    const part = parts.find((part) => part.id === partId)
    if (part?.serial_number) {
      return part.serial_number
    }
    // If part not found, show a more user-friendly message
    return 'N/A (Part not found)'
  }

  // Get part name by ID
  const getPartName = (partId: string): string => {
    const part = parts.find((part) => part.id === partId)
    if (part?.part_name) {
      return part.part_name
    }
    // If part not found, show a more user-friendly message
    return 'N/A (Part not found)'
  }

  // Claim items table columns
  const claimItemColumns: ColumnsType<ClaimItem> = [
    {
      title: 'Serial no.',
      dataIndex: 'faulty_part_id',
      key: 'serial_number',
      width: '15%',
      render: (_, record: ClaimItem) => {
        const serialNumber = getPartSerialNumber(record.faulty_part_id)
        const isNotFound = serialNumber.includes('not found')
        return (
          <Text type={isNotFound ? 'secondary' : undefined} italic={isNotFound}>
            {serialNumber}
          </Text>
        )
      },
    },
    {
      title: 'Part Category',
      dataIndex: 'part_category_id',
      key: 'part_category_id',
      width: '15%',
      render: (_, record: ClaimItem) => <Text>{getPartCategoryName(record.part_category_id)}</Text>,
    },
    {
      title: 'Faulty Part',
      dataIndex: 'faulty_part_id',
      key: 'faulty_part_id',
      width: '12%',
      render: (_, record: ClaimItem) => {
        const partName = getPartName(record.faulty_part_id)
        const isNotFound = partName.includes('not found')
        return (
          <Text type={isNotFound ? 'secondary' : undefined} italic={isNotFound}>
            {partName}
          </Text>
        )
      },
    },
    {
      title: 'Issue Description',
      dataIndex: 'issue_description',
      key: 'issue_description',
      width: '25%',
      render: (text: string) => (
        <div
          style={{
            whiteSpace: 'normal',
            wordBreak: 'break-word',
          }}
        >
          {text}
        </div>
      ),
    },
    {
      title: 'Type',
      dataIndex: 'type',
      key: 'type',
      width: '10%',
      render: (type: string) => (
        <>{CLAIM_ITEM_TYPE_LABELS[type as keyof typeof CLAIM_ITEM_TYPE_LABELS] || type}</>
      ),
    },
    {
      title: 'Cost',
      dataIndex: 'cost',
      key: 'cost',
      width: '10%',
      align: 'right',
      render: (cost: number) => (
        <Text strong style={{ color: '#52c41a' }}>
          {cost.toLocaleString('vi-VN', {
            style: 'currency',
            currency: 'VND',
          })}
        </Text>
      ),
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      width: '13%',
      render: (status: string) => (
        <>{CLAIM_ITEM_STATUS_LABELS[status as keyof typeof CLAIM_ITEM_STATUS_LABELS] || status}</>
      ),
    },
  ]

  return (
    <Card
      title={
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Title level={4}>
            <ToolOutlined /> Claim Items
          </Title>
          {canAddItems && (
            <Button type="primary" icon={<PlusOutlined />} onClick={onAddItem}>
              Add Item
            </Button>
          )}
        </div>
      }
      loading={loading}
    >
      <Table
        dataSource={claimItems}
        columns={claimItemColumns}
        rowKey="id"
        pagination={false}
        locale={{ emptyText: 'No claim items found' }}
      />
      {claimItems.length > 0 && <Divider />}
    </Card>
  )
}

export default ClaimItemsTable
