import { ROLE_LABELS, USER_ROLES } from '@constants'
import { Button, Popconfirm, Space, Tag, Tooltip } from 'antd'
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  DeleteOutlined,
  EditOutlined,
  HomeOutlined,
  MailOutlined,
  UserOutlined,
} from '@ant-design/icons'

const GenerateColumns = (sortedInfo, filteredInfo, onOpenModal, onDelete, getOfficeName) => {
  return [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
      sorter: (a, b) => (a.name || '').localeCompare(b.name || ''),
      sortOrder: sortedInfo.columnKey === 'name' ? sortedInfo.order : null,
      render: (text) => (
        <Space style={{ padding: '0 14px' }}>
          <UserOutlined style={{ color: '#697565' }} />
          <span>{text || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: 'Email',
      dataIndex: 'email',
      key: 'email',
      sorter: (a, b) => (a.email || '').localeCompare(b.email || ''),
      sortOrder: sortedInfo.columnKey === 'email' ? sortedInfo.order : null,
      render: (text) => (
        <Space style={{ padding: '0 14px' }}>
          <MailOutlined style={{ color: '#697565' }} />
          <span>{text || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: 'Role',
      dataIndex: 'role',
      key: 'role',
      align: 'center',
      width: 150,
      filters: Object.values(USER_ROLES).map((role) => ({
        text: ROLE_LABELS[role],
        value: role,
      })),
      filteredValue: filteredInfo.role || null,
      onFilter: (value, record) => record.role === value,
      render: (role) => {
        return <Space>{ROLE_LABELS[role] || role}</Space>
      },
    },
    {
      title: 'Office',
      dataIndex: 'office_id',
      key: 'office_id',
      align: 'center',
      render: (officeId) => (
        <Space>
          <HomeOutlined style={{ color: '#697565' }} />
          <span>{getOfficeName(officeId)}</span>
        </Space>
      ),
    },
    {
      title: 'Status',
      dataIndex: 'is_active',
      key: 'is_active',
      align: 'center',
      width: 125,
      filters: [
        { text: 'Active', value: true },
        { text: 'Inactive', value: false },
      ],
      filteredValue: filteredInfo.is_active || null,
      onFilter: (value, record) => record.is_active === value,
      render: (isActive) => (
        <Tag
          icon={isActive ? <CheckCircleOutlined /> : <CloseCircleOutlined />}
          color={isActive ? 'green' : 'red'}
        >
          {isActive ? 'Active' : 'Inactive'}
        </Tag>
      ),
    },
    {
      title: 'Actions',
      key: 'action',
      fixed: 'right',
      align: 'center',
      width: 100,
      render: (_, record) => (
        <Space size="small">
          <Tooltip title="Edit">
            <Button type="text" icon={<EditOutlined />} onClick={() => onOpenModal(record, true)} />
          </Tooltip>
          <Popconfirm
            title="Delete user"
            description="Are you sure you want to delete this user?"
            onConfirm={() => onDelete(record.id)}
            okText="Delete"
            cancelText="Cancel"
            okButtonProps={{ danger: true }}
          >
            <Tooltip title="Delete">
              <Button type="text" danger icon={<DeleteOutlined />} />
            </Tooltip>
          </Popconfirm>
        </Space>
      ),
    },
  ]
}

export default GenerateColumns
