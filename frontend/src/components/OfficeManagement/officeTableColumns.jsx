import { Button, Popconfirm, Space, Tag } from 'antd'
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  DeleteOutlined,
  EditOutlined,
  EnvironmentOutlined,
  BankOutlined,
} from '@ant-design/icons'

const OFFICE_TYPE_LABELS = {
  evm: 'EVM',
  sc: 'Service Center',
}

const GenerateColumns = (sortedInfo, filteredInfo, onOpenModal, onDelete) => {
  return [
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Office Name</span>,
      dataIndex: 'office_name',
      key: 'office_name',
      width: '25%',
      sorter: (a, b) => (a.office_name || '').localeCompare(b.office_name || ''),
      sortOrder: sortedInfo.columnKey === 'office_name' ? sortedInfo.order : null,
      render: (text) => (
        <Space style={{ padding: '0 14px', whiteSpace: 'normal', wordBreak: 'break-word' }}>
          <BankOutlined style={{ color: '#697565' }} />
          <span>{text || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: 'Office Type',
      dataIndex: 'office_type',
      key: 'office_type',
      align: 'center',
      width: '15%',
      filters: [
        { text: 'EVM', value: 'evm' },
        { text: 'Service Center', value: 'sc' },
      ],
      filteredValue: filteredInfo.office_type || null,
      onFilter: (value, record) => record.office_type === value,
      render: (office_type) => {
        const label = OFFICE_TYPE_LABELS[office_type] || office_type
        const color = office_type === 'evm' ? 'blue' : 'green'
        return <Tag color={color}>{label}</Tag>
      },
    },
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Address</span>,
      dataIndex: 'address',
      key: 'address',
      width: '30%',
      sorter: (a, b) => (a.address || '').localeCompare(b.address || ''),
      sortOrder: sortedInfo.columnKey === 'address' ? sortedInfo.order : null,
      ellipsis: true,
      render: (text) => (
        <Space style={{ padding: '0 14px' }}>
          <EnvironmentOutlined style={{ color: '#697565' }} />
          <span>{text || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: 'Status',
      dataIndex: 'is_active',
      key: 'is_active',
      align: 'center',
      width: '12%',
      filters: [
        { text: 'Active', value: true },
        { text: 'Inactive', value: false },
      ],
      filteredValue: filteredInfo.is_active || null,
      onFilter: (value, record) => record.is_active === value,
      render: (is_active) => {
        const color = is_active ? 'success' : 'error'
        const icon = is_active ? <CheckCircleOutlined /> : <CloseCircleOutlined />
        const text = is_active ? 'Active' : 'Inactive'
        return (
          <Tag color={color} icon={icon}>
            {text}
          </Tag>
        )
      },
    },
    {
      title: 'Actions',
      key: 'actions',
      align: 'center',
      width: '18%',
      render: (_, record) => (
        <Space size="small">
          <Button
            type="text"
            size="small"
            icon={<EditOutlined />}
            onClick={() => onOpenModal(record, true)}
            style={{ color: '#1890ff' }}
          >
            Edit
          </Button>
          <Popconfirm
            title="Delete Office"
            description="Are you sure you want to delete this office?"
            onConfirm={() => onDelete(record.id)}
            okText="Yes"
            cancelText="No"
            placement="topRight"
          >
            <Button type="text" size="small" icon={<DeleteOutlined />} danger>
              Delete
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]
}

export default GenerateColumns
