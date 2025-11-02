import { Button, Popconfirm, Space } from 'antd'
import {
  DeleteOutlined,
  EditOutlined,
  CarOutlined,
  CalendarOutlined,
  TagOutlined,
} from '@ant-design/icons'
import { type VehicleModel } from '@/types/index'

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
  //   _additionalProps: Record<string, unknown>
) => {
  return [
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Brand</span>,
      dataIndex: 'brand',
      key: 'brand',
      width: '25%',
      sorter: (a: Record<string, unknown>, b: Record<string, unknown>) => {
        const aModel = a as unknown as VehicleModel
        const bModel = b as unknown as VehicleModel
        return (aModel.brand || '').localeCompare(bModel.brand || '')
      },
      sortOrder:
        sortedInfo.columnKey === 'brand' ? (sortedInfo.order as 'ascend' | 'descend' | null) : null,
      render: (text: string) => (
        <Space
          style={{
            padding: '0 14px',
            whiteSpace: 'normal',
            wordBreak: 'break-word',
          }}
        >
          <TagOutlined style={{ color: '#697565' }} />
          <span>{text || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Model Name</span>,
      dataIndex: 'model_name',
      key: 'model_name',
      width: '30%',
      sorter: (a: Record<string, unknown>, b: Record<string, unknown>) => {
        const aModel = a as unknown as VehicleModel
        const bModel = b as unknown as VehicleModel
        return (aModel.model_name || '').localeCompare(bModel.model_name || '')
      },
      sortOrder:
        sortedInfo.columnKey === 'model_name'
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
          <CarOutlined style={{ color: '#697565' }} />
          <span>{text || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: 'Year',
      dataIndex: 'year',
      key: 'year',
      width: '15%',
      align: 'center' as const,
      sorter: (a: Record<string, unknown>, b: Record<string, unknown>) => {
        const aModel = a as unknown as VehicleModel
        const bModel = b as unknown as VehicleModel
        return (aModel.year || 0) - (bModel.year || 0)
      },
      sortOrder:
        sortedInfo.columnKey === 'year' ? (sortedInfo.order as 'ascend' | 'descend' | null) : null,
      render: (year: number) => (
        <Space>
          <CalendarOutlined style={{ color: '#697565' }} />
          <span>{year || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: 'Created Date',
      dataIndex: 'created_at',
      key: 'created_at',
      align: 'center' as const,
      width: '20%',
      sorter: (a: Record<string, unknown>, b: Record<string, unknown>) => {
        const aModel = a as unknown as VehicleModel
        const bModel = b as unknown as VehicleModel
        const aDate = new Date(aModel.created_at || '')
        const bDate = new Date(bModel.created_at || '')
        return aDate.getTime() - bDate.getTime()
      },
      sortOrder:
        sortedInfo.columnKey === 'created_at'
          ? (sortedInfo.order as 'ascend' | 'descend' | null)
          : null,
      render: (date: string) => (
        <Space style={{ padding: '0 14px' }}>
          <CalendarOutlined style={{ color: '#697565' }} />
          <span>{date ? new Date(date).toLocaleDateString() : 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: 'Actions',
      key: 'action',
      fixed: 'right' as const,
      align: 'center' as const,
      width: '20%',
      render: (_: unknown, record: Record<string, unknown>) => {
        const vehicleModel = record as unknown as VehicleModel
        return (
          <Space size="small">
            <Button
              type="text"
              icon={<EditOutlined />}
              onClick={() =>
                onOpenModal(record as Record<string, unknown> & { id: string | number }, true)
              }
            >
              Edit
            </Button>
            <Popconfirm
              title="Delete vehicle model"
              description="Are you sure you want to delete this vehicle model?"
              onConfirm={() => onDelete(vehicleModel.id)}
              okText="Delete"
              cancelText="Cancel"
              okButtonProps={{ danger: true }}
            >
              <Button type="text" danger icon={<DeleteOutlined />}>
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
