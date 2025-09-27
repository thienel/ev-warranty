import React from 'react'
import { Button, Layout, Typography, Avatar, Space } from 'antd'
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  SettingOutlined,
  LogoutOutlined,
  UserOutlined,
} from '@ant-design/icons'
import './AppHeader.less'

const { Header } = Layout
const { Text } = Typography

const Index = ({ collapsed, onToggleCollapse }) => {
  return (
    <Header className="app-header">
      <div className="header-left">
        <Button
          type="text"
          icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
          onClick={onToggleCollapse}
          className="toggle-button"
        />

        <Text className="header-title">Dashboard</Text>
      </div>

      <div className="header-right">
        <Space className="header-actions">
          <Button type="text" icon={<SettingOutlined />} />
          <Button type="text" icon={<LogoutOutlined />} />
          <Avatar className="ant-avatar">
            <UserOutlined />
          </Avatar>
        </Space>
      </div>
    </Header>
  )
}

export default Index
