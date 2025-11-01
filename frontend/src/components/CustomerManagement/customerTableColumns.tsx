import { Button, Popconfirm, Space, Tag } from 'antd'
import {
  DeleteOutlined,
  EditOutlined,
  MailOutlined,
  PhoneOutlined,
  UserOutlined,
  HomeOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
} from '@ant-design/icons'
import { type Customer } from '@/types/index'

type OnOpenModal = (
  item?: (Record<string, unknown> & { id: string | number }) | null,
  isUpdate?: boolean,
) => void
type OnDelete = (itemId: string | number) => Promise<void>

const GenerateColumns = (
  sortedInfo: Record<string, unknown>,
  filteredInfo: Record<string, unknown>,
  onOpenModal: OnOpenModal,
  onDelete: OnDelete,
  //   _additionalProps: Record<string, unknown>
) => {
  return [
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Name</span>,
      dataIndex: 'first_name',
      key: 'name',
      width: '20%',
      sorter: (a: Record<string, unknown>, b: Record<string, unknown>) => {
        const aCustomer = a as unknown as Customer
        const bCustomer = b as unknown as Customer
        const aName = `${aCustomer.first_name || ''} ${aCustomer.last_name || ''}`
        const bName = `${bCustomer.first_name || ''} ${bCustomer.last_name || ''}`
        return aName.localeCompare(bName)
      },
      sortOrder:
        sortedInfo.columnKey === 'name' ? (sortedInfo.order as 'ascend' | 'descend' | null) : null,
      render: (_text: string, record: Record<string, unknown>) => {
        const customer = record as unknown as Customer
        const fullName = `${customer.first_name || ''} ${customer.last_name || ''}`.trim()
        return (
          <Space
            style={{
              padding: '0 14px',
              whiteSpace: 'normal',
              wordBreak: 'break-word',
            }}
          >
            <UserOutlined style={{ color: '#697565' }} />
            <span>{fullName || 'N/A'}</span>
          </Space>
        )
      },
    },
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Email</span>,
      dataIndex: 'email',
      key: 'email',
      width: '22%',
      sorter: (a: Record<string, unknown>, b: Record<string, unknown>) => {
        const aCustomer = a as unknown as Customer
        const bCustomer = b as unknown as Customer
        return (aCustomer.email || '').localeCompare(bCustomer.email || '')
      },
      sortOrder:
        sortedInfo.columnKey === 'email' ? (sortedInfo.order as 'ascend' | 'descend' | null) : null,
      ellipsis: true,
      render: (text: string) => (
        <Space style={{ padding: '0 14px' }}>
          <MailOutlined style={{ color: '#697565' }} />
          <span>{text || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Phone</span>,
      dataIndex: 'phone_number',
      key: 'phone_number',
      width: '18%',
      sorter: (a: Record<string, unknown>, b: Record<string, unknown>) => {
        const aCustomer = a as unknown as Customer
        const bCustomer = b as unknown as Customer
        return (aCustomer.phone_number || '').localeCompare(bCustomer.phone_number || '')
      },
      sortOrder:
        sortedInfo.columnKey === 'phone_number'
          ? (sortedInfo.order as 'ascend' | 'descend' | null)
          : null,
      render: (text: string) => (
        <Space style={{ padding: '0 14px' }}>
          <PhoneOutlined style={{ color: '#697565' }} />
          <span>{text || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: 'Address',
      dataIndex: 'address',
      key: 'address',
      width: '25%',
      ellipsis: true,
      render: (text: string) => (
        <Space style={{ padding: '0 14px' }}>
          <HomeOutlined style={{ color: '#697565' }} />
          <span>{text || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: 'Status',
      dataIndex: 'is_deleted',
      key: 'is_deleted',
      align: 'center' as const,
      width: '10%',
      filters: [
        { text: 'Active', value: false },
        { text: 'Deleted', value: true },
      ],
      filteredValue: (filteredInfo.is_deleted as React.Key[] | null) || null,
      onFilter: (value: boolean | React.Key, record: Record<string, unknown>) => {
        const customer = record as unknown as Customer
        return customer.is_deleted === value
      },
      render: (isDeleted: boolean) => (
        <Tag
          icon={!isDeleted ? <CheckCircleOutlined /> : <CloseCircleOutlined />}
          color={!isDeleted ? 'green' : 'red'}
        >
          {!isDeleted ? 'Active' : 'Deleted'}
        </Tag>
      ),
    },
    {
      title: 'Actions',
      key: 'action',
      fixed: 'right' as const,
      align: 'center' as const,
      width: '10%',
      render: (_: unknown, record: Record<string, unknown>) => {
        const customer = record as unknown as Customer
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
              title="Delete customer"
              description="Are you sure you want to delete this customer?"
              onConfirm={() => onDelete(customer.id)}
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
