import React from 'react'
import { Button, Layout, Typography, Avatar, Space, message } from 'antd'
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  SettingOutlined,
  LogoutOutlined,
  UserOutlined,
} from '@ant-design/icons'
import './LayoutHeader.less'
import api from '@services/api'
import { API_ENDPOINTS } from '@constants/common-constants.js'
import { useDispatch } from 'react-redux'
import { logout } from '@redux/authSlice'
import { persistor } from '@redux/store'

const { Header } = Layout
const { Text } = Typography

interface LayoutHeaderProps {
  collapsed: boolean
  onToggleCollapse: () => void
  title?: string
}

const LayoutHeader: React.FC<LayoutHeaderProps> = ({ collapsed, onToggleCollapse, title }) => {
  const dispatch = useDispatch()

  const handleLogout = async () => {
    try {
      // Call logout API to invalidate refresh token on server
      await api.post(API_ENDPOINTS.AUTH.LOGOUT, {}, { withCredentials: true })
    } catch (error) {
      // Ignore errors from logout API (e.g., if token already expired)
      console.warn('Logout API error (ignored):', error)
    } finally {
      // Always clear local state regardless of API result
      dispatch(logout())
      await persistor.purge()
      message.success('Logout successful!')
    }
  }

  return (
    <Header className="app-header">
      <div className="header-left">
        <Button
          type="text"
          icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
          onClick={onToggleCollapse}
          className="toggle-button"
        />

        <Text className="header-title">{title}</Text>
      </div>

      <div className="header-right">
        <Space className="header-actions">
          <Button className="action-button" type="text" icon={<SettingOutlined />} />
          <Button
            type="text"
            className="action-button"
            onClick={handleLogout}
            icon={<LogoutOutlined />}
          />
          <Avatar className="user-avatar">
            <UserOutlined />
          </Avatar>
        </Space>
      </div>
    </Header>
  )
}

export default LayoutHeader
