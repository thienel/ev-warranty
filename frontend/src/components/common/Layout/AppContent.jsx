import React from 'react'
import { Layout, Typography } from 'antd'

const { Content } = Layout
const { Text } = Typography

const AppContent = ({ children }) => {
  return (
    <div className="app-content">
      <Content>
        {children || (
          <>
            {/* Content Header */}
            <div className="content-header">
              <Text className="content-title">Main Content Area</Text>
              <Text className="content-subtitle">
                This is your main content area where you can display your application content
              </Text>
            </div>

            {/* Demo Content */}
            <div className="demo-content">
              <Text className="demo-text">Your content goes here...</Text>
            </div>
          </>
        )}
      </Content>
    </div>
  )
}

export default AppContent
