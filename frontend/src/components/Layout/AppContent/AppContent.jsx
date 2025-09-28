import React from 'react'
import { Layout } from 'antd'
import './AppContent.less'

const { Content } = Layout

const AppContent = ({ children }) => {
  return (
    <div className="app-content">
      <Content>{children}</Content>
    </div>
  )
}

export default AppContent
