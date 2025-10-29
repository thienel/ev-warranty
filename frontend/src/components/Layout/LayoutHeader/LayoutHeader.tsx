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
import useHandleApiError from '@/hooks/useHandleApiError'
import type { ErrorResponse } from '@/constants/error-messages'

const { Header } = Layout
const { Text } = Typography

interface LayoutHeaderProps {
  collapsed: boolean
  onToggleCollapse: () => void
  title?: string
}

const LayoutHeader: React.FC<LayoutHeaderProps> = ({ collapsed, onToggleCollapse, title }) => {
  const dispatch = useDispatch()
  const handleError = useHandleApiError()

  const handleLogout = async () => {
    try {
      await api.post(API_ENDPOINTS.AUTH.LOGOUT, {}, { withCredentials: true })
      dispatch(logout())
      await persistor.purge()
      message.success('Logout successful!')
    } catch (error) {
      handleError(error as ErrorResponse)
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
