import React from 'react'
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
import { useNavigate } from 'react-router-dom'

const { Sider } = Layout

const Sidebar = ({ collapsed }) => {
  const menuItems = [
    { key: '1', icon: <BarChartOutlined />, label: 'Reports' },
    { key: '2', icon: <UserOutlined />, label: 'Users' },
    { key: '3', icon: <BankOutlined />, label: 'Offices' },
    { key: '4', icon: <UserOutlined />, label: 'Customers' },
    { key: '5', icon: <CarOutlined />, label: 'Vehicles' },
    { key: '6', icon: <ContainerOutlined />, label: 'Warranty claims' },
  ]

  const navigate = useNavigate()

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
        defaultSelectedKeys={['1']}
        items={menuItems}
        onClick={({ key }) => {
          switch (key) {
            case '1':
              navigate('/reports')
              break
            case '2':
              navigate('/users')
              break
            case '3':
              navigate('/offices')
              break
            case '4':
              navigate('/customers')
              break
            case '5':
              navigate('/vehicles')
              break
            case '6':
              navigate('/claims')
              break
            default:
              break
          }
        }}
      />
    </Sider>
  )
}

export default Sidebar
