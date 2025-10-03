import React from 'react'
import { Layout } from 'antd'
import './LayoutContent.less'

const { Content } = Layout

const LayoutContent = ({ children }) => {
  return (
    <div className="app-content">
      <Content>{children}</Content>
    </div>
  )
}

export default LayoutContent
