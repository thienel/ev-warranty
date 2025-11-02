import React from 'react'
import { Card, Descriptions, Typography, Tag } from 'antd'
import { SafetyOutlined } from '@ant-design/icons'
import type { WarrantyPolicy } from '@/types/index'

const { Title, Text, Paragraph } = Typography

interface WarrantyPolicyCardProps {
  warrantyPolicy: WarrantyPolicy | null
  loading: boolean
}

const WarrantyPolicyCard: React.FC<WarrantyPolicyCardProps> = ({ warrantyPolicy, loading }) => {
  const formatDuration = (months: number): string => {
    if (months >= 12) {
      const years = Math.floor(months / 12)
      const remainingMonths = months % 12
      if (remainingMonths === 0) {
        return `${years} ${years === 1 ? 'year' : 'years'}`
      }
      return `${years} ${years === 1 ? 'year' : 'years'} ${remainingMonths} ${remainingMonths === 1 ? 'month' : 'months'}`
    }
    return `${months} ${months === 1 ? 'month' : 'months'}`
  }

  const formatKilometerLimit = (limit?: number): string => {
    if (!limit) return 'No limit'
    return `${limit.toLocaleString()} km`
  }

  return (
    <Card
      title={
        <Title level={4}>
          <SafetyOutlined /> Warranty Policy Information
        </Title>
      }
      loading={loading}
      style={{ marginBottom: '24px' }}
    >
      {warrantyPolicy ? (
        <>
          <Descriptions bordered column={1} style={{ marginBottom: '16px' }}>
            <Descriptions.Item label="Policy Name">
              <Text strong>{warrantyPolicy.policy_name}</Text>
            </Descriptions.Item>
            <Descriptions.Item label="Warranty Duration">
              <Tag color="blue">{formatDuration(warrantyPolicy.warranty_duration_months)}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="Kilometer Limit">
              <Tag color={warrantyPolicy.kilometer_limit ? 'orange' : 'green'}>
                {formatKilometerLimit(warrantyPolicy.kilometer_limit)}
              </Tag>
            </Descriptions.Item>
            {warrantyPolicy.vehicle_models && warrantyPolicy.vehicle_models.length > 0 && (
              <Descriptions.Item label="Applicable Models">
                <div style={{ display: 'flex', flexWrap: 'wrap', gap: '4px' }}>
                  {warrantyPolicy.vehicle_models.map((model) => (
                    <Tag key={model.id} color="purple">
                      {model.brand} {model.model_name} ({model.year})
                    </Tag>
                  ))}
                </div>
              </Descriptions.Item>
            )}
          </Descriptions>

          <div>
            <Text strong style={{ fontSize: '14px', marginBottom: '8px', display: 'block' }}>
              Terms and Conditions:
            </Text>
            <Paragraph
              style={{
                backgroundColor: '#f5f5f5',
                padding: '12px',
                borderRadius: '4px',
                margin: 0,
                whiteSpace: 'pre-wrap',
              }}
            >
              {warrantyPolicy.terms_and_conditions}
            </Paragraph>
          </div>
        </>
      ) : (
        <Text type="secondary">No warranty policy information available</Text>
      )}
    </Card>
  )
}

export default WarrantyPolicyCard
