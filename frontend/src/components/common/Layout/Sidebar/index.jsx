import React from 'react'
import { Layout, Menu } from 'antd'
import {
  DashboardOutlined,
  UserOutlined,
  VideoCameraOutlined,
  UploadOutlined,
  SettingOutlined,
  ThunderboltFilled,
} from '@ant-design/icons'
import './Sidebar.less'

const { Sider } = Layout

const Sidebar = ({ collapsed }) => {
  const menuItems = [
    {
      key: '1',
      icon: <DashboardOutlined />,
      label: 'Dashboard',
      style: { marginTop: '8px' },
    },
    { key: '2', icon: <UserOutlined />, label: 'Users' },
    { key: '3', icon: <VideoCameraOutlined />, label: 'Media' },
    { key: '4', icon: <UploadOutlined />, label: 'Uploads' },
    {
      key: '5',
      icon: <SettingOutlined />,
      label: 'Settings',
      style: { marginTop: 'auto' },
    },
  ]

  return (
    <Sider
      trigger={null}
      collapsible
      collapsed={collapsed}
      width={260}
      collapsedWidth={80}
      className="sidebar"
    >
      <div className="sidebar-header">
        <ThunderboltFilled />
        <div className={`sidebar-title ${collapsed ? 'collapsed' : 'expanded'}`}>
          EV Warranty System
        </div>
      </div>
      <Menu theme="dark" mode="inline" defaultSelectedKeys={['1']} items={menuItems} />
    </Sider>
  )
}

export default Sidebar
