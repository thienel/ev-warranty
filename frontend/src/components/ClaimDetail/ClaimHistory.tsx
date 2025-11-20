import React from 'react'
import { Card, Timeline, Tag, Typography, Empty } from 'antd'
import { ClockCircleOutlined, HistoryOutlined, UserOutlined } from '@ant-design/icons'
import { CLAIM_STATUS_LABELS } from '@constants/common-constants'
import type { ClaimHistory as ClaimHistoryType, User } from '@/types/index'

const { Title, Text } = Typography

interface ClaimHistoryProps {
  history: ClaimHistoryType[]
  users: User[]
  loading: boolean
}

const ClaimHistory: React.FC<ClaimHistoryProps> = ({ history, users, loading }) => {
  // Get user name by ID
  const getUserName = (userId: string): string => {
    const user = users.find((u) => u.id === userId)
    return user?.name || `User ${userId.slice(0, 8)}...`
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
                  <Tag style={{ fontSize: '13px' }}>
                    {CLAIM_STATUS_LABELS[record.status] || record.status}
                  </Tag>
                </div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                  <UserOutlined style={{ fontSize: '12px', color: '#999' }} />
                  <Text type="secondary" style={{ fontSize: '12px' }}>
                    Changed by: <Text strong>{getUserName(record.changed_by)}</Text>
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
