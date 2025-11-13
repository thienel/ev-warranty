import React from 'react'
import { Card, Descriptions, Typography } from 'antd'
import { TeamOutlined } from '@ant-design/icons'
import type { User } from '@/types/index'

const { Title, Text } = Typography

interface TechnicianInfoProps {
  technician: User | null
  loading: boolean
}

const TechnicianInfo: React.FC<TechnicianInfoProps> = ({ technician, loading }) => {
  console.log(technician)
  return (
    <Card
      title={
        <Title level={4}>
          <TeamOutlined /> Assigned Technician
        </Title>
      }
      loading={loading}
    >
      {technician ? (
        <Descriptions bordered column={1}>
          <Descriptions.Item label="Name">
            <Text strong>{technician.name}</Text>
          </Descriptions.Item>
          <Descriptions.Item label="Email">
            <Text>{technician.email}</Text>
          </Descriptions.Item>
        </Descriptions>
      ) : !loading ? (
        <Text type="secondary">Technician information not available</Text>
      ) : null}
    </Card>
  )
}

export default TechnicianInfo
