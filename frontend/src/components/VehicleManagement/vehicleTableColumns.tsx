import { Button, Popconfirm, Space } from 'antd'
import {
  DeleteOutlined,
  EditOutlined,
  CarOutlined,
  UserOutlined,
  FileTextOutlined,
  CalendarOutlined,
} from '@ant-design/icons'
import { type Vehicle } from '@/types/index'

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
  additionalProps: Record<string, unknown>,
) => {
  const { getCustomerName, getVehicleModelName } = additionalProps as {
    getCustomerName?: (customerId: string) => string
    getVehicleModelName?: (modelId: string) => string
  }

  return [
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>VIN</span>,
      dataIndex: 'vin',
      key: 'vin',
      width: '225px',
      sorter: (a: Record<string, unknown>, b: Record<string, unknown>) => {
        const aVehicle = a as unknown as Vehicle
        const bVehicle = b as unknown as Vehicle
        return (aVehicle.vin || '').localeCompare(bVehicle.vin || '')
      },
      sortOrder:
        sortedInfo.columnKey === 'vin' ? (sortedInfo.order as 'ascend' | 'descend' | null) : null,
      render: (text: string) => (
        <Space
          style={{
            padding: '0 14px',
            whiteSpace: 'normal',
            wordBreak: 'break-word',
          }}
        >
          <FileTextOutlined style={{ color: '#697565' }} />
          <span>{text || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>License Plate</span>,
      dataIndex: 'license_plate',
      key: 'license_plate',
      width: '225px',
      sorter: (a: Record<string, unknown>, b: Record<string, unknown>) => {
        const aVehicle = a as unknown as Vehicle
        const bVehicle = b as unknown as Vehicle
        return (aVehicle.license_plate || '').localeCompare(bVehicle.license_plate || '')
      },
      sortOrder:
        sortedInfo.columnKey === 'license_plate'
          ? (sortedInfo.order as 'ascend' | 'descend' | null)
          : null,
      render: (text: string) => (
        <Space style={{ padding: '0 14px' }}>
          <CarOutlined style={{ color: '#697565' }} />
          <span>{text || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: 'Customer',
      dataIndex: 'customer_id',
      key: 'customer_id',
      render: (customerId: string) => (
        <Space style={{ padding: '0 14px' }}>
          <UserOutlined style={{ color: '#697565' }} />
          <span>{getCustomerName?.(customerId) || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: 'Vehicle Model',
      dataIndex: 'model_id',
      key: 'model_id',
      render: (modelId: string) => (
        <Space style={{ padding: '0 14px' }}>
          <CarOutlined style={{ color: '#697565' }} />
          <span>{getVehicleModelName?.(modelId) || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: 'Purchase Date',
      dataIndex: 'purchase_date',
      key: 'purchase_date',
      width: '150px',
      align: 'center' as const,
      sorter: (a: Record<string, unknown>, b: Record<string, unknown>) => {
        const aVehicle = a as unknown as Vehicle
        const bVehicle = b as unknown as Vehicle
        const aDate = new Date(aVehicle.purchase_date || '')
        const bDate = new Date(bVehicle.purchase_date || '')
        return aDate.getTime() - bDate.getTime()
      },
      sortOrder:
        sortedInfo.columnKey === 'purchase_date'
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
      width: '100px',
      render: (_: unknown, record: Record<string, unknown>) => {
        const vehicle = record as unknown as Vehicle
        return (
          <Space size="small">
            <Button
              type="text"
              icon={<EditOutlined />}
              onClick={() =>
                onOpenModal(record as Record<string, unknown> & { id: string | number }, true)
              }
            />
            <Popconfirm
              title="Delete vehicle"
              description="Are you sure you want to delete this vehicle?"
              onConfirm={() => onDelete(vehicle.id)}
              okText="Delete"
              cancelText="Cancel"
              okButtonProps={{ danger: true }}
            >
              <Button type="text" danger icon={<DeleteOutlined />} />
            </Popconfirm>
          </Space>
        )
      },
    },
  ]
}

export default GenerateColumns
