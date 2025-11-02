import React from 'react'
import { Card, Descriptions, Typography } from 'antd'
import { UserOutlined } from '@ant-design/icons'
import type { Customer } from '@/types/index'

const { Title, Text } = Typography

interface CustomerInfoProps {
  customer: Customer | null
  loading: boolean
}

const CustomerInfo: React.FC<CustomerInfoProps> = ({ customer, loading }) => {
  return (
    <Card
      title={
        <Title level={4}>
          <UserOutlined /> Customer Information
        </Title>
      }
      loading={loading}
    >
      {customer ? (
        <Descriptions bordered column={1}>
          <Descriptions.Item label="Name">
            <Text strong>
              {customer.full_name || `${customer.first_name} ${customer.last_name}`}
            </Text>
          </Descriptions.Item>
          {customer.email && <Descriptions.Item label="Email">{customer.email}</Descriptions.Item>}
          {customer.phone_number && (
            <Descriptions.Item label="Phone">{customer.phone_number}</Descriptions.Item>
          )}
          {customer.address && (
            <Descriptions.Item label="Address">{customer.address}</Descriptions.Item>
          )}
        </Descriptions>
      ) : !loading ? (
        <Text type="secondary">Customer information not available</Text>
      ) : null}
    </Card>
  )
}

export default CustomerInfo
