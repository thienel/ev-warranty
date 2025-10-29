import React from 'react'
import { Layout } from 'antd'
import './LayoutContent.less'

const { Content } = Layout

interface LayoutContentProps {
  children: React.ReactNode
}

const LayoutContent: React.FC<LayoutContentProps> = ({ children }) => {
  return (
    <div className="app-content">
      <Content>{children}</Content>
    </div>
  )
}

export default LayoutContent
