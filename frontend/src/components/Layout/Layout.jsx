import React, { useState } from 'react'
import { Layout } from 'antd'
import Sidebar from './Sidebar/Sidebar.jsx'
import AppHeader from './AppHeader/AppHeader.jsx'
import AppContent from './AppContent/AppContent.jsx'
import './Layout.less'

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
