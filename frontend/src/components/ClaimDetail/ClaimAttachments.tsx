import React from 'react'
import { Card, Row, Col, Image, Button, Tag, Typography, Space } from 'antd'
import {
  PaperClipOutlined,
  FileTextOutlined,
  PlusOutlined,
  PlayCircleOutlined,
} from '@ant-design/icons'
import type { ClaimAttachment } from '@/types/index'
import { ATTACHMENT_TYPE_LABELS } from '@/constants/common-constants'

const { Title, Text } = Typography

interface ClaimAttachmentsProps {
  attachments: ClaimAttachment[]
  loading: boolean
  canAddAttachments?: boolean
  onAddAttachment?: () => void
}

const ClaimAttachments: React.FC<ClaimAttachmentsProps> = ({
  attachments,
  loading,
  canAddAttachments = false,
  onAddAttachment,
}) => {
  return (
    <Card
      title={
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Title level={4}>
            <PaperClipOutlined /> Attachments
          </Title>
          {canAddAttachments && onAddAttachment && (
            <Button type="primary" icon={<PlusOutlined />} onClick={onAddAttachment}>
              Add Attachment
            </Button>
          )}
        </div>
      }
      loading={loading}
    >
      {attachments.length > 0 ? (
        <Row gutter={[16, 16]}>
          {attachments.map((attachment) => (
            <Col key={attachment.id} xs={24} sm={12} md={8} lg={6}>
              <Card
                hoverable
                cover={
                  attachment.type === 'image' ? (
                    <Image
                      src={attachment.url}
                      alt="Attachment"
                      style={{ height: 200, objectFit: 'cover' }}
                    />
                  ) : attachment.type === 'video' ? (
                    <div
                      style={{
                        height: 200,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        backgroundColor: '#f0f0f0',
                        position: 'relative',
                      }}
                    >
                      <PlayCircleOutlined style={{ fontSize: 48, color: '#1890ff' }} />
                      <div
                        style={{
                          position: 'absolute',
                          bottom: 8,
                          left: 8,
                          right: 8,
                          background: 'rgba(0, 0, 0, 0.7)',
                          color: 'white',
                          padding: '4px 8px',
                          borderRadius: '4px',
                          fontSize: '12px',
                          textAlign: 'center',
                        }}
                      >
                        Video File
                      </div>
                    </div>
                  ) : (
                    <div
                      style={{
                        height: 200,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        backgroundColor: '#f0f0f0',
                      }}
                    >
                      <FileTextOutlined style={{ fontSize: 48, color: '#999' }} />
                    </div>
                  )
                }
              >
                <Card.Meta
                  title={
                    <Tag color="blue">
                      {ATTACHMENT_TYPE_LABELS[
                        attachment.type as keyof typeof ATTACHMENT_TYPE_LABELS
                      ] || attachment.type}
                    </Tag>
                  }
                  description={
                    <Space direction="vertical" size={0}>
                      <Text type="secondary" style={{ fontSize: '12px' }}>
                        {attachment.created_at
                          ? new Date(attachment.created_at).toLocaleDateString()
                          : 'N/A'}
                      </Text>
                      <Button
                        type="link"
                        size="small"
                        href={attachment.url}
                        target="_blank"
                        style={{ padding: 0 }}
                      >
                        {attachment.type === 'video' ? 'Play Video' : 'View Full'}
                      </Button>
                    </Space>
                  }
                />
              </Card>
            </Col>
          ))}
        </Row>
      ) : (
        <Text type="secondary">No attachments found</Text>
      )}
    </Card>
  )
}

export default ClaimAttachments
