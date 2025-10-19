import React, { useState } from 'react'
import { Layout } from 'antd'
import Sidebar from './Sidebar/Sidebar.jsx'
import LayoutHeader from './LayoutHeader/LayoutHeader.jsx'
import LayoutContent from './LayoutContent/LayoutContent.jsx'
import './Layout.less'

const AppLayout = ({ children, title }) => {
  const [collapsed, setCollapsed] = useState(false)

  const handleToggleCollapse = () => {
    setCollapsed(!collapsed)
  }

  return (
    <Layout className="app-layout" style={{ height: '100vh', overflow: 'hidden' }}>
      <Sidebar collapsed={collapsed} />

      <Layout style={{ height: '100vh', overflow: 'hidden' }}>
        <LayoutHeader collapsed={collapsed} onToggleCollapse={handleToggleCollapse} title={title} />
        <LayoutContent>{children}</LayoutContent>
      </Layout>
    </Layout>
  )
}

export default AppLayout
