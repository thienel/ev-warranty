import React, { useState, useEffect, useMemo } from 'react'
import {
  Table,
  Button,
  Modal,
  Form,
  Input,
  Select,
  Switch,
  Space,
  message,
  Popconfirm,
  Tag,
  Row,
  Col,
  Tooltip,
  Card,
  Statistic,
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SearchOutlined,
  ReloadOutlined,
  UserOutlined,
  MailOutlined,
  HomeOutlined,
  LockOutlined,
  TeamOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
} from '@ant-design/icons'
import api from '@services/api'
import { USER_ROLES, API_ENDPOINTS } from '@constants'

const { Option } = Select

const UserManagement = () => {
  const [users, setUsers] = useState([])
  const [offices, setOffices] = useState([])
  const [loading, setLoading] = useState(false)
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [editingUser, setEditingUser] = useState(null)
  const [searchText, setSearchText] = useState('')
  const [filteredInfo, setFilteredInfo] = useState({})
  const [sortedInfo, setSortedInfo] = useState({})
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  })
  const [form] = Form.useForm()

  const roleLabels = {
    [USER_ROLES.ADMIN]: 'Admin',
    [USER_ROLES.SC_STAFF]: 'SC Staff',
    [USER_ROLES.SC_TECHNICIAN]: 'SC Technician',
    [USER_ROLES.EVM_STAFF]: 'EVM Staff',
  }

  const fetchUsers = async () => {
    setLoading(true)
    try {
      const response = await api.get(API_ENDPOINTS.USER)

      if (response.data.success) {
        const userData = response.data.data || []
        setUsers(userData)
        setPagination((prev) => ({
          ...prev,
          total: userData.length,
        }))
      }
    } catch (error) {
      message.error(error.response?.data?.message || 'Failed to load users')
      console.error('Error fetching users:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchOffices = async () => {
    try {
      const response = await api.get(API_ENDPOINTS.OFFICE)

      if (response.data.success) {
        setOffices(response.data.data || [])
      }
    } catch (error) {
      message.error(error.response?.data?.message || 'Failed to load offices')
      console.error('Error fetching offices:', error)
    }
  }

  useEffect(() => {
    fetchUsers()
    fetchOffices()
  }, [])

  const filteredUsers = useMemo(() => {
    if (!searchText) return users

    const searchLower = searchText.toLowerCase()
    return users.filter(
      (user) =>
        user.name?.toLowerCase().includes(searchLower) ||
        user.email?.toLowerCase().includes(searchLower) ||
        roleLabels[user.role]?.toLowerCase().includes(searchLower)
    )
  }, [users, searchText])

  // Calculate statistics
  const statistics = useMemo(() => {
    return {
      total: users.length,
      active: users.filter((u) => u.is_active).length,
      inactive: users.filter((u) => !u.is_active).length,
      admins: users.filter((u) => u.role === USER_ROLES.ADMIN).length,
    }
  }, [users])

  // Handle form submission
  const handleSubmit = async (values) => {
    setLoading(true)
    try {
      let response

      // Prepare payload
      const payload = { ...values }

      if (!editingUser && !payload.password) {
        message.warning('Password is required for new user')
        setLoading(false)
        return
      }

      // Remove password field if editing and password is empty
      if (editingUser && !payload.password) {
        delete payload.password
      }

      if (editingUser) {
        response = await api.put(`${API_ENDPOINTS.USER.UPDATE}${editingUser.id}`, payload)

        if (response.data.success) {
          message.success('User updated successfully')
        }
      } else {
        response = await api.post(API_ENDPOINTS.USER.CREATE, payload)

        if (response.data.success) {
          message.success('User created successfully')
        }
      }

      setIsModalVisible(false)
      form.resetFields()
      setEditingUser(null)
      fetchUsers()
    } catch (error) {
      const errorMsg =
        error.response?.data?.message ||
        (editingUser ? 'Failed to update user' : 'Failed to create user')
      message.error(errorMsg)
      console.error('Error saving user:', error)
    } finally {
      setLoading(false)
    }
  }

  // Handle delete
  const handleDelete = async (userId) => {
    setLoading(true)
    try {
      const response = await api.delete(`${API_ENDPOINTS.USER.DELETE}${userId}`)

      if (response.data.success) {
        message.success('User deleted successfully')
        fetchUsers()
      }
    } catch (error) {
      message.error(error.response?.data?.message || 'Failed to delete user')
      console.error('Error deleting user:', error)
    } finally {
      setLoading(false)
    }
  }

  // Open modal for create/edit
  const openModal = (user = null) => {
    setEditingUser(user)
    if (user) {
      form.setFieldsValue({
        ...user,
        password: '', // Don't show password
      })
    } else {
      form.resetFields()
      form.setFieldsValue({ is_active: true })
    }
    setIsModalVisible(true)
  }

  // Get office name by ID
  const getOfficeName = (officeId) => {
    const office = offices.find((o) => o.id === officeId)
    return office ? office.name : 'N/A'
  }

  // Handle table change (pagination, filters, sorter)
  const handleTableChange = (newPagination, filters, sorter) => {
    setPagination(newPagination)
    setFilteredInfo(filters)
    setSortedInfo(sorter)
  }

  // Reset filters and sorters
  const handleReset = () => {
    setSearchText('')
    setFilteredInfo({})
    setSortedInfo({})
    setPagination({
      current: 1,
      pageSize: 10,
      total: users.length,
    })
  }

  const columns = [
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
        text: roleLabels[role],
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
        return <Tag color={colors[role] || 'default'}>{roleLabels[role] || role}</Tag>
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
            <Button type="text" icon={<EditOutlined />} onClick={() => openModal(record)} />
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

  return (
    <div style={{ padding: '24px' }}>
      {/* Statistics Cards */}
      <Row gutter={[16, 16]} style={{ marginBottom: '24px' }}>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="Total Users"
              value={statistics.total}
              prefix={<TeamOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="Active Users"
              value={statistics.active}
              prefix={<CheckCircleOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="Inactive Users"
              value={statistics.inactive}
              prefix={<CloseCircleOutlined />}
              valueStyle={{ color: '#8c8c8c' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="Admins"
              value={statistics.admins}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#f5222d' }}
            />
          </Card>
        </Col>
      </Row>

      {/* Search and Actions */}
      <Row gutter={[16, 16]} style={{ marginBottom: '16px' }}>
        <Col xs={24} sm={24} md={12} lg={16}>
          <Input
            placeholder="Search by name, email or role..."
            prefix={<SearchOutlined />}
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
            allowClear
            size="large"
          />
        </Col>
        <Col xs={24} sm={24} md={12} lg={8} style={{ textAlign: 'right' }}>
          <Space wrap>
            <Button icon={<ReloadOutlined />} onClick={handleReset}>
              Reset
            </Button>
            <Button icon={<ReloadOutlined />} onClick={fetchUsers} loading={loading}>
              Refresh
            </Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={() => openModal()}>
              Add User
            </Button>
          </Space>
        </Col>
      </Row>

      {/* Users Table */}
      <Table
        columns={columns}
        rowKey="id"
        loading={loading}
        dataSource={filteredUsers}
        pagination={{
          ...pagination,
          total: filteredUsers.length,
          showTotal: (total, range) => `${range[0]}-${range[1]} of ${total} users`,
          showSizeChanger: true,
          showQuickJumper: true,
          pageSizeOptions: ['10', '20', '50', '100'],
        }}
        onChange={handleTableChange}
        scroll={{ x: 1000 }}
        bordered
      />

      {/* Create/Edit Modal */}
      <Modal
        title={
          <Space>
            <UserOutlined />
            {editingUser ? 'Edit User' : 'Add New User'}
          </Space>
        }
        open={isModalVisible}
        onCancel={() => {
          setIsModalVisible(false)
          form.resetFields()
          setEditingUser(null)
        }}
        footer={null}
        width={600}
        destroyOnHidden
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit} autoComplete="off">
          <Form.Item
            label="Full Name"
            name="name"
            rules={[
              { required: true, message: 'Please enter full name' },
              { min: 2, message: 'Name must be at least 2 characters' },
              { max: 100, message: 'Name cannot exceed 100 characters' },
              {
                pattern: /^[a-zA-ZÀ-ỹ\s]+$/,
                message: 'Name can only contain letters and spaces',
              },
            ]}
          >
            <Input placeholder="Enter full name" prefix={<UserOutlined />} size="large" />
          </Form.Item>

          <Form.Item
            label="Email"
            name="email"
            rules={[
              { required: true, message: 'Please enter email' },
              { type: 'email', message: 'Invalid email format' },
              { max: 100, message: 'Email cannot exceed 100 characters' },
            ]}
          >
            <Input
              placeholder="Enter email"
              prefix={<MailOutlined />}
              size="large"
              disabled={!!editingUser} // Don't allow email change when editing
            />
          </Form.Item>

          <Form.Item
            label={editingUser ? 'Password (leave blank to keep current)' : 'Password'}
            name="password"
            rules={
              editingUser
                ? [
                    { min: 6, message: 'Password must be at least 6 characters' },
                    { max: 50, message: 'Password cannot exceed 50 characters' },
                  ]
                : [
                    { required: true, message: 'Please enter password' },
                    { min: 6, message: 'Password must be at least 6 characters' },
                    { max: 50, message: 'Password cannot exceed 50 characters' },
                  ]
            }
          >
            <Input.Password
              placeholder={editingUser ? 'Enter new password (optional)' : 'Enter password'}
              prefix={<LockOutlined />}
              size="large"
            />
          </Form.Item>

          <Form.Item
            label="Role"
            name="role"
            rules={[{ required: true, message: 'Please select a role' }]}
          >
            <Select placeholder="Select role" size="large">
              {Object.entries(USER_ROLES).map(([key, value]) => (
                <Option key={key} value={value}>
                  <Space>
                    <UserOutlined />
                    {roleLabels[value]}
                  </Space>
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            label="Office"
            name="office_id"
            rules={[{ required: true, message: 'Please select an office' }]}
          >
            <Select
              placeholder="Select office"
              size="large"
              showSearch
              optionFilterProp="children"
              filterOption={(input, option) =>
                option.children.toLowerCase().includes(input.toLowerCase())
              }
            >
              {offices.map((office) => (
                <Option key={office.id} value={office.id}>
                  <Space>
                    <HomeOutlined />
                    {office.name}
                  </Space>
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item label="Status" name="is_active" valuePropName="checked">
            <Switch checkedChildren="Active" unCheckedChildren="Inactive" size="default" />
          </Form.Item>

          <Form.Item style={{ marginBottom: 0, textAlign: 'right', marginTop: '24px' }}>
            <Space>
              <Button
                size="large"
                onClick={() => {
                  setIsModalVisible(false)
                  form.resetFields()
                  setEditingUser(null)
                }}
              >
                Cancel
              </Button>
              <Button type="primary" htmlType="submit" loading={loading} size="large">
                {editingUser ? 'Update User' : 'Create User'}
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default UserManagement
