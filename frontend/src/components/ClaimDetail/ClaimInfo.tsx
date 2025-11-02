import React from 'react'
import { Card, Descriptions, Tag, Typography } from 'antd'
import { FileTextOutlined } from '@ant-design/icons'
import type { Claim } from '@/types/index'
import { CLAIM_STATUSES } from '@/constants/common-constants'

const { Title, Text } = Typography

interface ClaimInfoProps {
  claim: Claim | null
  loading: boolean
}

const getStatusColor = (status: string) => {
  switch (status) {
    case CLAIM_STATUSES.DRAFT:
      return 'gray'
    case CLAIM_STATUSES.SUBMITTED:
      return 'blue'
    case CLAIM_STATUSES.REQUEST_INFO:
      return 'orange'
    case CLAIM_STATUSES.REVIEWING:
      return 'processing'
    case CLAIM_STATUSES.APPROVED:
      return 'green'
    case CLAIM_STATUSES.REJECTED:
      return 'red'
    default:
      return 'default'
  }
}

const ClaimInfo: React.FC<ClaimInfoProps> = ({ claim, loading }) => {
  return (
    <Card
      title={
        <Title level={4}>
          <FileTextOutlined /> Claim Information
        </Title>
      }
      loading={loading}
    >
      {claim ? (
        <Descriptions bordered column={1}>
          <Descriptions.Item label="Claim ID">
            <Text code strong>
              {claim.id}
            </Text>
          </Descriptions.Item>
          <Descriptions.Item label="Status">
            <Tag color={getStatusColor(claim.status)}>{claim.status}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="Description">
            <Text>{claim.description}</Text>
          </Descriptions.Item>
          {claim.created_at && (
            <Descriptions.Item label="Created At">
              {new Date(claim.created_at).toLocaleString()}
            </Descriptions.Item>
          )}
          {claim.updated_at && (
            <Descriptions.Item label="Last Updated">
              {new Date(claim.updated_at).toLocaleString()}
            </Descriptions.Item>
          )}
        </Descriptions>
      ) : !loading ? (
        <Text type="secondary">Claim information not available</Text>
      ) : null}
    </Card>
  )
}

export default ClaimInfo
