import React, { useEffect, useState } from 'react'
import { Layout, Menu } from 'antd'
import {
  UserOutlined,
  ThunderboltFilled,
  BankOutlined,
  CarOutlined,
  ContainerOutlined,
  BarChartOutlined,
} from '@ant-design/icons'
import './Sidebar.less'
import { useNavigate, useLocation } from 'react-router-dom'

const { Sider } = Layout

const Sidebar = ({ collapsed }) => {
  const navigate = useNavigate()
  const location = useLocation()

  const menuItems = [
    { key: 'reports', icon: <BarChartOutlined />, label: 'Reports', path: '/reports' },
    { key: 'users', icon: <UserOutlined />, label: 'Users', path: '/admin/users' },
    { key: 'offices', icon: <BankOutlined />, label: 'Offices', path: '/admin/offices' },
    { key: 'customers', icon: <UserOutlined />, label: 'Customers', path: '/customers' },
    { key: 'vehicles', icon: <CarOutlined />, label: 'Vehicles', path: '/vehicles' },
    { key: 'claims', icon: <ContainerOutlined />, label: 'Warranty claims', path: '/claims' },
  ]

  const getCurrentKey = () => {
    const currentItem = menuItems.find((item) => item.path === location.pathname)
    return currentItem ? currentItem.key : 'users'
  }

  const [selectedKey, setSelectedKey] = useState(getCurrentKey())

  useEffect(() => {
    setSelectedKey(getCurrentKey())
  }, [location.pathname])

  const handleMenuClick = ({ key }) => {
    const menuItem = menuItems.find((item) => item.key === key)
    if (menuItem) {
      setSelectedKey(key)
      navigate(menuItem.path)
    }
  }

  return (
    <Sider
      trigger={null}
      collapsible
      collapsed={collapsed}
      width={260}
      collapsedWidth={80}
      className="sidebar"
    >
      <div className="sidebar-header" onClick={() => navigate('/')}>
        <ThunderboltFilled />
        <div className={`sidebar-title ${collapsed ? 'collapsed' : 'expanded'}`}>
          EV Warranty System
        </div>
      </div>
      <Menu
        theme="dark"
        mode="inline"
        selectedKeys={[selectedKey]}
        items={menuItems}
        onClick={handleMenuClick}
      />
    </Sider>
  )
}

export default Sidebar
