import React from 'react'
import { Layout, Typography } from 'antd'
import './AppContent.less'

const { Content } = Layout
const { Text } = Typography

const Index = ({ children }) => {
  return (
    <div className="app-content">
      <Content>{children}</Content>
    </div>
  )
}

export default Index
