import React from 'react'
import { Card, Row, Col, Image, Button, Tag, Typography, Space } from 'antd'
import { PaperClipOutlined, FileTextOutlined } from '@ant-design/icons'
import type { ClaimAttachment } from '@/types/index'
import { ATTACHMENT_TYPE_LABELS } from '@/constants/common-constants'

const { Title, Text } = Typography

interface ClaimAttachmentsProps {
  attachments: ClaimAttachment[]
  loading: boolean
}

const ClaimAttachments: React.FC<ClaimAttachmentsProps> = ({ attachments, loading }) => {
  return (
    <Card
      title={
        <Title level={4}>
          <PaperClipOutlined /> Attachments
        </Title>
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
                        View Full
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
