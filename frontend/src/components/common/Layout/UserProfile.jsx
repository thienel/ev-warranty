import React from 'react'
import { Space, Typography, Avatar } from 'antd'

const { Text } = Typography

const UserProfile = ({ collapsed }) => {
  if (collapsed) return null

  return (
    <div className="user-profile">
      <Space>
        <div className="user-avatar">
          <Avatar size="small" style={{ backgroundColor: '#87d068' }}>
            U
          </Avatar>
        </div>
        <div className="user-info">
          <Text className="user-name">John Doe</Text>
          <Text className="user-role">Administrator</Text>
        </div>
      </Space>
    </div>
  )
}

export default UserProfile
