import { Button, Popconfirm, Space } from 'antd'
import { DeleteOutlined, EditOutlined, SafetyOutlined, EyeOutlined } from '@ant-design/icons'
import { type WarrantyPolicy } from '@/types/index'
import { useNavigate } from 'react-router-dom'

type OnOpenModal = (
  item?: (Record<string, unknown> & { id: string | number }) | null,
  isUpdate?: boolean,
) => void
type OnDelete = (itemId: string | number) => Promise<void>

const GenerateColumns = (
  sortedInfo: Record<string, unknown>,
  _filteredInfo: Record<string, unknown>,
  onOpenModal: OnOpenModal,
  onDelete: OnDelete,
) => {
  const navigate = useNavigate()

  const handleViewDetails = (id: string) => {
    navigate(`/admin/policies/${id}`)
  }

  return [
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Policy Name</span>,
      dataIndex: 'policy_name',
      key: 'policy_name',
      width: '25%',
      sorter: (a: Record<string, unknown>, b: Record<string, unknown>) => {
        const aPolicy = a as unknown as WarrantyPolicy
        const bPolicy = b as unknown as WarrantyPolicy
        return (aPolicy.policy_name || '').localeCompare(bPolicy.policy_name || '')
      },
      sortOrder:
        sortedInfo.columnKey === 'policy_name'
          ? (sortedInfo.order as 'ascend' | 'descend' | null)
          : null,
      render: (text: string) => (
        <Space
          style={{
            padding: '0 14px',
            whiteSpace: 'normal',
            wordBreak: 'break-word',
          }}
        >
          <SafetyOutlined style={{ color: '#697565' }} />
          <span>{text || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Vehicle Models</span>,
      dataIndex: 'vehicle_models',
      key: 'vehicle_models',
      width: '25%',
      render: (_: unknown, record: Record<string, unknown>) => {
        const policy = record as unknown as WarrantyPolicy
        if (!policy.vehicle_models || policy.vehicle_models.length === 0) {
          return <span style={{ color: '#999', padding: '0 14px' }}>No models assigned</span>
        }
        const modelNames = policy.vehicle_models
          .map((model) => `${model.brand} ${model.model_name} (${model.year})`)
          .join(', ')
        return (
          <span
            style={{
              padding: '0 14px',
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
      align: 'center' as const,
      width: '12%',
      sorter: (a: Record<string, unknown>, b: Record<string, unknown>) => {
        const aPolicy = a as unknown as WarrantyPolicy
        const bPolicy = b as unknown as WarrantyPolicy
        return (aPolicy.warranty_duration_months || 0) - (bPolicy.warranty_duration_months || 0)
      },
      sortOrder:
        sortedInfo.columnKey === 'warranty_duration_months'
          ? (sortedInfo.order as 'ascend' | 'descend' | null)
          : null,
      render: (months: number) => <span>{months || 'N/A'}</span>,
    },
    {
      title: 'Kilometer Limit',
      dataIndex: 'kilometer_limit',
      key: 'kilometer_limit',
      align: 'center' as const,
      width: '12%',
      sorter: (a: Record<string, unknown>, b: Record<string, unknown>) => {
        const aPolicy = a as unknown as WarrantyPolicy
        const bPolicy = b as unknown as WarrantyPolicy
        return (aPolicy.kilometer_limit || 0) - (bPolicy.kilometer_limit || 0)
      },
      sortOrder:
        sortedInfo.columnKey === 'kilometer_limit'
          ? (sortedInfo.order as 'ascend' | 'descend' | null)
          : null,
      render: (limit?: number) => (limit ? `${limit.toLocaleString()} km` : 'N/A'),
    },
    {
      title: 'Details',
      key: 'details',
      align: 'center' as const,
      width: '8%',
      render: (_: unknown, record: Record<string, unknown>) => {
        const policy = record as unknown as WarrantyPolicy
        return (
          <Space size="middle">
            <Button
              type="text"
              icon={<EyeOutlined />}
              onClick={() => handleViewDetails(policy.id)}
              title="View Details"
              style={{ color: '#1890ff' }}
            />
          </Space>
        )
      },
    },
    {
      title: 'Actions',
      key: 'actions',
      align: 'center' as const,
      width: '18%',
      render: (_: unknown, record: Record<string, unknown>) => {
        const policy = record as unknown as WarrantyPolicy
        return (
          <Space size="small">
            <Button
              type="text"
              size="small"
              icon={<EditOutlined />}
              onClick={() =>
                onOpenModal(record as Record<string, unknown> & { id: string | number }, true)
              }
              style={{ color: '#1890ff' }}
            >
              Edit
            </Button>
            <Popconfirm
              title="Delete Policy"
              description="Are you sure you want to delete this policy?"
              onConfirm={() => onDelete(policy.id)}
              okText="Yes"
              cancelText="No"
              placement="topRight"
            >
              <Button type="text" size="small" icon={<DeleteOutlined />} danger>
                Delete
              </Button>
            </Popconfirm>
          </Space>
        )
      },
    },
  ]
}

export default GenerateColumns
