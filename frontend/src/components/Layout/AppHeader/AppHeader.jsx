import React from 'react'
import { Button, Layout, Typography, Avatar, Space, message } from 'antd'
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  SettingOutlined,
  LogoutOutlined,
  UserOutlined,
} from '@ant-design/icons'
import './AppHeader.less'
import api from '@services/api.js'
import { API_ENDPOINTS } from '@constants'
import { useDispatch } from 'react-redux'
import { logout } from '@redux/authSlice.js'
import { persistor } from '@redux/store.js'

const { Header } = Layout
const { Text } = Typography

const AppHeader = ({ collapsed, onToggleCollapse, title }) => {
  const dispatch = useDispatch()
  const handleLogout = async () => {
    try {
      const _ = await api.post(API_ENDPOINTS.AUTH.LOGOUT, {}, { withCredentials: true })
      dispatch(logout())
      await persistor.purge()
      message.success('Logout successful!')
    } catch (error) {
      console.log(error)
      message.error('Logout failed')
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

export default AppHeader
