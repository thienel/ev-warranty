import React, { useState, useEffect, useMemo } from 'react'
import { Button, Input, Space, message, Row, Col, Card, Statistic } from 'antd'
import {
  PlusOutlined,
  SearchOutlined,
  ReloadOutlined,
  UserOutlined,
  TeamOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
} from '@ant-design/icons'
import api from '@services/api'
import { USER_ROLES, API_ENDPOINTS } from '@constants'
import UserModal from '@components/UserManagement/UserModal/UserModal.jsx'
import UserTable from '@components/UserManagement/UserTable/UserTable.jsx'

const UserManagement = () => {
  const [users, setUsers] = useState([])
  const [offices, setOffices] = useState([])

  const [loading, setLoading] = useState(false)
  const [searchText, setSearchText] = useState('')

  const [updateUser, setUpdateUser] = useState(null)
  const [isOpenModal, setIsOpenModal] = useState(false)

  const [isResetTable, setIsResetTable] = useState(true)

  const handleOpenModal = (user = null) => {
    setUpdateUser(user)
    setIsOpenModal(true)
  }

  const handleCloseModal = () => {
    setIsOpenModal(false)
    setUpdateUser(null)
    fetchUsers()
  }

  const fetchUsers = async () => {
    setLoading(true)
    try {
      const response = await api.get(API_ENDPOINTS.USER)

      if (response.data.success) {
        const userData = response.data.data || []
        setUsers(userData)
        setIsResetTable(true)
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

  const statistics = useMemo(() => {
    return {
      total: users.length,
      active: users.filter((u) => u.is_active).length,
      inactive: users.filter((u) => !u.is_active).length,
      admins: users.filter((u) => u.role === USER_ROLES.ADMIN).length,
    }
  }, [users])

  const handleReset = () => {
    setSearchText('')
    setIsResetTable(true)
  }

  return (
    <div style={{ padding: '24px' }}>
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
            <Button type="primary" icon={<PlusOutlined />} onClick={handleOpenModal}>
              Add User
            </Button>
          </Space>
        </Col>
      </Row>

      <UserTable
        loading={loading}
        setLoading={setLoading}
        isReset={isResetTable}
        setIsReset={setIsResetTable}
        users={users}
        offices={offices}
        handleOpenModal={handleOpenModal}
        onRefresh={handleReset}
      />

      <UserModal
        loading={loading}
        setLoading={setLoading}
        onClose={handleCloseModal}
        user={updateUser}
        opened={isOpenModal}
        offices={offices}
      />
    </div>
  )
}

export default UserManagement
