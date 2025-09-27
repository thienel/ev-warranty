import React from 'react'
import { Typography } from 'antd'
import './Logo.less'

const { Text } = Typography

const Index = ({ collapsed }) => {
  return (
    <div className={`logo-section ${collapsed ? 'collapsed' : ''}`}>
      <div className="logo-icon">A</div>
      {!collapsed && <Text className="logo-text">Admin Panel</Text>}
    </div>
  )
}

export default Index
