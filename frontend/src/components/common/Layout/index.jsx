import React, { useState } from 'react'
import { Layout } from 'antd'
import Sidebar from './Sidebar'
import AppHeader from './AppHeader'
import AppContent from './AppContent'
import '@styles/main.less'

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
