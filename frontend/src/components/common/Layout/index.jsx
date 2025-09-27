import React, { useState } from 'react'
import { Layout } from 'antd'
import Sidebar from './Sidebar/index.jsx'
import AppHeader from './AppHeader/index.jsx'
import AppContent from '@components/common/Layout/AppContent/index.jsx'
import './layout.less'

const AppLayout = ({ children }) => {
  const [collapsed, setCollapsed] = useState(false)

  const handleToggleCollapse = () => {
    setCollapsed(!collapsed)
  }

  return (
    <Layout className="app-layout" style={{ height: '100vh', overflow: 'hidden' }}>
      <Sidebar collapsed={collapsed} />

      <Layout style={{ height: '100vh', overflow: 'hidden' }}>
        <AppHeader collapsed={collapsed} onToggleCollapse={handleToggleCollapse} />
        <AppContent>{children}</AppContent>
      </Layout>
    </Layout>
  )
}

export default AppLayout
