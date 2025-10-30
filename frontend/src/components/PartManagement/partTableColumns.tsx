import { Button, Popconfirm, Space, Tag } from 'antd'
import {
  DeleteOutlined,
  EditOutlined,
  ToolOutlined,
  BarcodeOutlined,
  AppstoreOutlined,
  HomeOutlined,
} from '@ant-design/icons'
import { type Part } from '@/types/index.js'

type OnOpenModal = (
  item?: (Record<string, unknown> & { id: string | number }) | null,
  isUpdate?: boolean,
) => void
type OnDelete = (itemId: string | number) => Promise<void>

const PART_STATUS_LABELS = {
  Available: 'Available',
  Reserved: 'Reserved',
  Installed: 'Installed',
  Defective: 'Defective',
  Obsolete: 'Obsolete',
  Archived: 'Archived',
} as const

const PART_STATUS_COLORS = {
  Available: 'green',
  Reserved: 'blue',
  Installed: 'cyan',
  Defective: 'red',
  Obsolete: 'orange',
  Archived: 'default',
} as const

const GenerateColumns = (
  sortedInfo: Record<string, unknown>,
  _filteredInfo: Record<string, unknown>,
  onOpenModal: OnOpenModal,
  onDelete: OnDelete,
  additionalProps?: { getOfficeName?: (officeId: string) => string },
) => {
  const getOfficeName = additionalProps?.getOfficeName || (() => 'N/A')

  return [
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Serial Number</span>,
      dataIndex: 'serial_number',
      key: 'serial_number',
      width: '15%',
      sorter: (a: Record<string, unknown>, b: Record<string, unknown>) => {
        const aPart = a as unknown as Part
        const bPart = b as unknown as Part
        return (aPart.serial_number || '').localeCompare(bPart.serial_number || '')
      },
      sortOrder:
        sortedInfo.columnKey === 'serial_number'
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
          <BarcodeOutlined style={{ color: '#697565' }} />
          <span>{text || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Part Name</span>,
      dataIndex: 'part_name',
      key: 'part_name',
      width: '18%',
      sorter: (a: Record<string, unknown>, b: Record<string, unknown>) => {
        const aPart = a as unknown as Part
        const bPart = b as unknown as Part
        return (aPart.part_name || '').localeCompare(bPart.part_name || '')
      },
      sortOrder:
        sortedInfo.columnKey === 'part_name'
          ? (sortedInfo.order as 'ascend' | 'descend' | null)
          : null,
      render: (text: string) => (
        <Space style={{ padding: '0 14px', whiteSpace: 'normal', wordBreak: 'break-word' }}>
          <ToolOutlined style={{ color: '#697565' }} />
          <span>{text || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: 'Unit Price',
      dataIndex: 'unit_price',
      key: 'unit_price',
      align: 'right' as const,
      width: '12%',
      sorter: (a: Record<string, unknown>, b: Record<string, unknown>) => {
        const aPart = a as unknown as Part
        const bPart = b as unknown as Part
        return (aPart.unit_price || 0) - (bPart.unit_price || 0)
      },
      sortOrder:
        sortedInfo.columnKey === 'unit_price'
          ? (sortedInfo.order as 'ascend' | 'descend' | null)
          : null,
      render: (price: number) => (
        <Space style={{ padding: '0 14px' }}>
          <span>
            {price.toLocaleString('vi-VN', {
              style: 'currency',
              currency: 'VND',
            })}
          </span>
        </Space>
      ),
    },
    {
      title: 'Category',
      dataIndex: 'category_name',
      key: 'category_name',
      align: 'center' as const,
      width: '15%',
      render: (text: string) => (
        <Space>
          <AppstoreOutlined style={{ color: '#697565' }} />
          <span>{text || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: 'Office Location',
      dataIndex: 'office_location_id',
      key: 'office_location_id',
      align: 'center' as const,
      width: '13%',
      render: (officeId: string) => (
        <Space>
          <HomeOutlined style={{ color: '#697565' }} />
          <span>{officeId ? getOfficeName(officeId) : 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      align: 'center' as const,
      width: '12%',
      render: (status: string) => {
        const label = PART_STATUS_LABELS[status as keyof typeof PART_STATUS_LABELS] || status
        const color = PART_STATUS_COLORS[status as keyof typeof PART_STATUS_COLORS] || 'default'
        return <Tag color={color}>{label}</Tag>
      },
    },
    {
      title: 'Actions',
      key: 'actions',
      align: 'center' as const,
      width: '15%',
      render: (_: unknown, record: Record<string, unknown>) => {
        const part = record as unknown as Part
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
              title="Delete Part"
              description="Are you sure you want to delete this part?"
              onConfirm={() => onDelete(part.id)}
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
