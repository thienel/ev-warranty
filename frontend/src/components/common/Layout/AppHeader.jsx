import React from 'react'
import { Button, Layout, Typography, Avatar, Space } from 'antd'
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  SettingOutlined,
  LogoutOutlined,
} from '@ant-design/icons'

const { Header } = Layout
const { Text } = Typography

const AppHeader = ({ collapsed, onToggleCollapse }) => {
  return (
    <Header className="app-header">
      <div className="header-left">
        <Button
          type="text"
          icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
          onClick={onToggleCollapse}
          className="toggle-button"
        />

        <div className="header-title">
          <Text className="title">Dashboard</Text>
          <Text className="subtitle">Welcome back, manage your application</Text>
        </div>
      </div>

      <div className="header-right">
        <Space className="header-actions">
          <Button type="text" icon={<SettingOutlined />} />
          <Button type="text" icon={<LogoutOutlined />} />
          <Avatar>JD</Avatar>
        </Space>
      </div>
    </Header>
  )
}

export default AppHeader
