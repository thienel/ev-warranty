import React, { useEffect, useState } from 'react'
import { Layout, Menu } from 'antd'
import { UserOutlined, ThunderboltFilled, BankOutlined } from '@ant-design/icons'
import './Sidebar.less'
import { useNavigate, useLocation } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { USER_ROLES } from '@constants'

const { Sider } = Layout

const Sidebar = ({ collapsed }) => {
  const navigate = useNavigate()
  const location = useLocation()
  const { user } = useSelector((state) => state.auth)

  const [menuItems, setMenuItems] = useState([])

  useEffect(() => {
    switch (user.role) {
      case USER_ROLES.ADMIN:
        setMenuItems([
          { key: 'users', icon: <UserOutlined />, label: 'Users', path: '/admin/users' },
          { key: 'offices', icon: <BankOutlined />, label: 'Offices', path: '/admin/offices' },
        ])
        break
      default:
        setMenuItems([])
    }
  }, [user.role])

  useEffect(() => {
    if (menuItems.length > 0) {
      const currentItem = menuItems.find((item) => item.path === location.pathname)

      if (!currentItem) {
        navigate(menuItems[0].path, { replace: true })
      }
    }
  }, [menuItems, location.pathname, navigate])

  const selectedKey =
    menuItems.find((item) => item.path === location.pathname)?.key || menuItems[0]?.key

  const handleMenuClick = ({ key }) => {
    const menuItem = menuItems.find((item) => item.key === key)
    if (menuItem) {
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
