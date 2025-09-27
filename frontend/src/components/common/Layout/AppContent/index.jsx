import React from 'react'
import { Layout, Typography } from 'antd'
import './AppContent.less'

const { Content } = Layout
const { Text } = Typography

const Index = ({ children }) => {
  return (
    <div className="app-content">
      <Content>
        {children || (
          <>
            <div className="content-header">
              <Text className="content-title">Main Content Area</Text>
            </div>

            <div className="demo-content">
              <Text className="demo-text">Your content goes here...</Text>
            </div>
          </>
        )}
      </Content>
    </div>
  )
}

export default Index
