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

const GenerateColumns = (
  sortedInfo,
  filteredInfo,
  handleOpenModal,
  handleDelete,
  getOfficeName
) => {
  return [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
      sorter: (a, b) => (a.name || '').localeCompare(b.name || ''),
      sortOrder: sortedInfo.columnKey === 'name' ? sortedInfo.order : null,
      render: (text) => (
        <Space>
          <UserOutlined style={{ color: '#1890ff' }} />
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
        <Space>
          <MailOutlined style={{ color: '#52c41a' }} />
          <span>{text || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: 'Role',
      dataIndex: 'role',
      key: 'role',
      filters: Object.values(USER_ROLES).map((role) => ({
        text: ROLE_LABELS[role],
        value: role,
      })),
      filteredValue: filteredInfo.role || null,
      onFilter: (value, record) => record.role === value,
      render: (role) => {
        const colors = {
          [USER_ROLES.ADMIN]: 'red',
          [USER_ROLES.SC_STAFF]: 'blue',
          [USER_ROLES.SC_TECHNICIAN]: 'green',
          [USER_ROLES.EVM_STAFF]: 'orange',
        }
        return <Tag color={colors[role] || 'default'}>{ROLE_LABELS[role] || role}</Tag>
      },
    },
    {
      title: 'Office',
      dataIndex: 'office_id',
      key: 'office_id',
      render: (officeId) => (
        <Space>
          <HomeOutlined style={{ color: '#722ed1' }} />
          <span>{getOfficeName(officeId)}</span>
        </Space>
      ),
    },
    {
      title: 'Status',
      dataIndex: 'is_active',
      key: 'is_active',
      filters: [
        { text: 'Active', value: true },
        { text: 'Inactive', value: false },
      ],
      filteredValue: filteredInfo.is_active || null,
      onFilter: (value, record) => record.is_active === value,
      render: (isActive) => (
        <Tag
          icon={isActive ? <CheckCircleOutlined /> : <CloseCircleOutlined />}
          color={isActive ? 'success' : 'default'}
        >
          {isActive ? 'Active' : 'Inactive'}
        </Tag>
      ),
    },
    {
      title: 'Actions',
      key: 'action',
      fixed: 'right',
      width: 120,
      render: (_, record) => (
        <Space size="small">
          <Tooltip title="Edit">
            <Button
              type="text"
              icon={<EditOutlined />}
              onClick={() => handleOpenModal(record, true)}
            />
          </Tooltip>
          <Popconfirm
            title="Delete user"
            description="Are you sure you want to delete this user?"
            onConfirm={() => handleDelete(record.id)}
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
