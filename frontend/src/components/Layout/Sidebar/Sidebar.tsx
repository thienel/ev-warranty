import React, { useCallback, useEffect, useState } from 'react'
import { Layout, Menu } from 'antd'
import type { MenuProps } from 'antd'
import {
  UserOutlined,
  ThunderboltFilled,
  FileTextOutlined,
  SafetyOutlined,
  CarOutlined,
  InboxOutlined,
  AppstoreOutlined,
  ToolOutlined,
  TeamOutlined,
  HomeOutlined,
} from '@ant-design/icons'
import './Sidebar.less'
import { useNavigate, useLocation } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { USER_ROLES, type UserRole } from '@constants/common-constants.js'
import type { RootState } from '@redux/store'

const { Sider } = Layout

interface SidebarProps {
  collapsed: boolean
}

type MenuItem = Required<MenuProps>['items'][number]

interface MenuItemData {
  key: string
  icon: React.ReactNode
  label: string
  path: string
}

const Sidebar: React.FC<SidebarProps> = ({ collapsed }) => {
  const navigate = useNavigate()
  const location = useLocation()
  const { user } = useSelector((state: RootState) => state.auth)

  const [menuItems, setMenuItems] = useState<MenuItem[]>([])

  useEffect(() => {
    if (user?.role) {
      const items = MENU_ITEMS[user.role as UserRole] || []
      setMenuItems(
        items.map((item) => ({
          key: item.key,
          icon: item.icon,
          label: item.label,
        })),
      )
    }
  }, [user?.role])

  useEffect(() => {
    if (user?.role) {
      const menuData = MENU_ITEMS[user.role as UserRole] || []
      if (menuData.length > 0) {
        const currentItem = menuData.find(
          (item) =>
            location.pathname === item.path || location.pathname.startsWith(item.path + '/'),
        )

        if (!currentItem) {
          navigate(menuData[0].path, { replace: true })
        }
      }
    }
  }, [user?.role, location.pathname, navigate])

  const getMenuData = (): MenuItemData[] => {
    return user?.role ? MENU_ITEMS[user.role as UserRole] || [] : []
  }

  const selectedKey = (() => {
    const menuData = getMenuData()
    return (
      menuData.find(
        (item) => location.pathname === item.path || location.pathname.startsWith(item.path + '/'),
      )?.key || menuData[0]?.key
    )
  })()

  const handleMenuClick: MenuProps['onClick'] = ({ key }) => {
    const menuData = getMenuData()
    const menuItem = menuData.find((item) => item.key === key)
    if (menuItem) {
      navigate(menuItem.path)
    }
  }

  const handleHomeClick = useCallback(() => {
    navigate('/')
  }, [navigate])

  return (
    <Sider
      trigger={null}
      collapsible
      collapsed={collapsed}
      width={260}
      collapsedWidth={80}
      className="sidebar"
    >
      <div className="sidebar-header" onClick={handleHomeClick}>
        <ThunderboltFilled />
        <div className={`sidebar-title ${collapsed ? 'collapsed' : 'expanded'}`}>
          EV Warranty System
        </div>
      </div>
      <Menu
        theme="dark"
        mode="inline"
        selectedKeys={selectedKey ? [selectedKey] : []}
        items={menuItems}
        onClick={handleMenuClick}
      />
    </Sider>
  )
}

const MENU_ITEMS: Record<UserRole, MenuItemData[]> = {
  [USER_ROLES.ADMIN]: [
    {
      key: 'users',
      icon: <UserOutlined />,
      label: 'Users',
      path: '/admin/users',
    },
    {
      key: 'offices',
      icon: <HomeOutlined />,
      label: 'Offices',
      path: '/admin/offices',
    },
    {
      key: 'policies',
      icon: <SafetyOutlined />,
      label: 'Policies',
      path: '/admin/policies',
    },
    {
      key: 'models',
      icon: <AppstoreOutlined />,
      label: 'Models',
      path: '/admin/vehicle-models',
    },
    {
      key: 'part-categories',
      icon: <ToolOutlined />,
      label: 'Part Categories',
      path: '/admin/part-categories',
    },
  ],
  [USER_ROLES.EVM_STAFF]: [
    {
      key: 'claims',
      icon: <FileTextOutlined />,
      label: 'Claims',
      path: '/evm-staff/claims',
    },
    {
      key: 'vehicles',
      icon: <CarOutlined />,
      label: 'Vehicles',
      path: '/evm-staff/vehicles',
    },
    {
      key: 'inventories',
      icon: <InboxOutlined />,
      label: 'Inventories',
      path: '/evm-staff/inventories',
    },
  ],
  [USER_ROLES.SC_STAFF]: [
    {
      key: 'claims',
      icon: <FileTextOutlined />,
      label: 'Claims',
      path: '/sc-staff/claims',
    },
    {
      key: 'customers',
      icon: <TeamOutlined />,
      label: 'Customers',
      path: '/sc-staff/customers',
    },
  ],
  [USER_ROLES.SC_TECHNICIAN]: [
    {
      key: 'claims',
      icon: <FileTextOutlined />,
      label: 'Claims',
      path: '/sc-technician/claims',
    },
  ],
}

export default Sidebar
