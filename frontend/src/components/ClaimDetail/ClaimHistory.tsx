import React from 'react'
import { Card, Timeline, Tag, Typography, Empty } from 'antd'
import { ClockCircleOutlined, HistoryOutlined, UserOutlined } from '@ant-design/icons'
import { CLAIM_STATUS_LABELS } from '@constants/common-constants'
import type { ClaimHistory as ClaimHistoryType } from '@/types/index'

const { Title, Text } = Typography

interface ClaimHistoryProps {
  history: ClaimHistoryType[]
  loading: boolean
}

const ClaimHistory: React.FC<ClaimHistoryProps> = ({ history, loading }) => {
  // Get status color for tags
  const getStatusColor = (status: string): string => {
    const colors: Record<string, string> = {
      DRAFT: 'default',
      SUBMITTED: 'blue',
      REVIEWING: 'cyan',
      REQUEST_INFO: 'orange',
      APPROVED: 'green',
      PARTIALLY_APPROVED: 'lime',
      REJECTED: 'red',
      CANCELLED: 'gray',
      COMPLETED: 'purple',
    }
    return colors[status] || 'default'
  }

  // Sort history by date (most recent first)
  const sortedHistory = [...history].sort(
    (a, b) => new Date(b.changed_at).getTime() - new Date(a.changed_at).getTime(),
  )

  return (
    <Card
      title={
        <Title level={4} style={{ margin: 0 }}>
          <HistoryOutlined /> Claim History
        </Title>
      }
      loading={loading}
    >
      {!loading && history.length === 0 ? (
        <Empty
          description="No history records found"
          image={Empty.PRESENTED_IMAGE_SIMPLE}
          style={{ margin: '20px 0' }}
        />
      ) : (
        <Timeline
          mode="left"
          items={sortedHistory.map((record) => ({
            key: record.id,
            dot: <ClockCircleOutlined style={{ fontSize: '16px' }} />,
            color: getStatusColor(record.status),
            label: (
              <Text type="secondary" style={{ fontSize: '12px' }}>
                {record.changed_at
                  ? new Date(record.changed_at).toLocaleString('en-US', {
                      year: 'numeric',
                      month: 'short',
                      day: 'numeric',
                      hour: '2-digit',
                      minute: '2-digit',
                    })
                  : 'Unknown time'}
              </Text>
            ),
            children: (
              <div style={{ marginTop: '-4px' }}>
                <div style={{ marginBottom: '8px' }}>
                  <Tag color={getStatusColor(record.status)} style={{ fontSize: '13px' }}>
                    {CLAIM_STATUS_LABELS[record.status] || record.status}
                  </Tag>
                </div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                  <UserOutlined style={{ fontSize: '12px', color: '#999' }} />
                  <Text type="secondary" style={{ fontSize: '12px' }}>
                    Changed by: <Text code>{record.changed_by.slice(0, 8)}...</Text>
                  </Text>
                </div>
              </div>
            ),
          }))}
        />
      )}
    </Card>
  )
}

export default ClaimHistory