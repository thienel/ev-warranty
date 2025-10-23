import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import {
  Card,
  Descriptions,
  Button,
  Space,
  Tag,
  List,
  Avatar,
  Divider,
  Modal,
  message,
  Popconfirm,
} from 'antd'
import {
  EditOutlined,
  DeleteOutlined,
  HistoryOutlined,
  FileOutlined,
  EyeOutlined,
  ArrowRightOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons'
// import { useSelector } from 'react-redux'
import AppLayout from '@components/Layout/Layout.jsx'
import { CLAIM_STATUS_LABELS, USER_ROLES, API_ENDPOINTS } from '@constants/common-constants.js'
import useCheckRole from '@/hooks/useCheckRole.js'
// import api from '@services/api'

const ClaimDetail = () => {
  const { id } = useParams()
  const navigate = useNavigate()
  // const { user } = useSelector((state) => state.auth)

  const [claim, setClaim] = useState(null)
  const [attachments, setAttachments] = useState([])
  const [statusHistory, setStatusHistory] = useState([])
  const [loading, setLoading] = useState(true)
  const [historyModalVisible, setHistoryModalVisible] = useState(false)
  const [actionLoading, setActionLoading] = useState(false)

  // Role checks
  const isScStaff = useCheckRole(USER_ROLES.SC_STAFF)
  const isScTechnician = useCheckRole(USER_ROLES.SC_TECHNICIAN)
  const canEdit = isScStaff || isScTechnician
  const canDelete = isScStaff

  // Mock claim data for demonstration
  const mockClaim = {
    id: id || 'claim-001-abc-def',
    status: 'DRAFT',
    customer_id: 'cust-123-xyz',
    customer_name: 'John Smith',
    customer_email: 'john.smith@email.com',
    customer_phone: '+84 123 456 789',
    vehicle_id: 'veh-456-abc',
    vehicle_info: 'Tesla Model 3 2023',
    vehicle_vin: 'VIN123456789',
    description:
      'Battery replacement issue - customer reports significant decrease in battery capacity after 2 years of use. Vehicle shows error codes related to battery management system.',
    total_cost: 25000000, // 25 million VND
    approved_by: null,
    created_at: '2024-10-20T10:30:00Z',
    updated_at: '2024-10-20T10:30:00Z',
  }

  const mockAttachments = [
    {
      id: '1',
      name: 'battery_diagnostic_report.pdf',
      size: 2500000,
      uploaded_at: '2024-10-20T10:30:00Z',
      uploaded_by: 'John Smith',
    },
    {
      id: '2',
      name: 'vehicle_photos.zip',
      size: 15000000,
      uploaded_at: '2024-10-20T11:15:00Z',
      uploaded_by: 'John Smith',
    },
    {
      id: '3',
      name: 'warranty_certificate.pdf',
      size: 850000,
      uploaded_at: '2024-10-20T10:45:00Z',
      uploaded_by: 'John Smith',
    },
  ]

  const mockStatusHistory = [
    {
      id: '1',
      status: 'DRAFT',
      changed_by: 'John Smith',
      changed_at: '2024-10-20T10:30:00Z',
      notes: 'Initial claim created',
    },
  ]

  useEffect(() => {
    fetchClaimDetails()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id])

  const fetchClaimDetails = async () => {
    try {
      setLoading(true)
      // TODO: Replace with actual API calls when backend is ready
      // const claimResponse = await api.get(`${API_ENDPOINTS.CLAIM}/${id}`)
      // const attachmentsResponse = await api.get(`${API_ENDPOINTS.CLAIM}/${id}/attachments`)
      // const historyResponse = await api.get(`${API_ENDPOINTS.CLAIM}/${id}/history`)

      // Mock data for now
      setClaim(mockClaim)
      setAttachments(mockAttachments)
      setStatusHistory(mockStatusHistory)
    } catch (error) {
      message.error('Failed to fetch claim details')
      console.error('Error fetching claim details:', error)
    } finally {
      setLoading(false)
    }
  }

  const getValidTransitions = (currentStatus) => {
    const transitions = {
      DRAFT: [
        { status: 'SUBMITTED', label: 'Submit', type: 'primary' },
        { status: 'CANCELLED', label: 'Cancel', type: 'default' },
      ],
      SUBMITTED: [
        { status: 'REVIEWING', label: 'Start Review', type: 'primary' },
        { status: 'CANCELLED', label: 'Cancel', type: 'default' },
      ],
      REQUEST_INFO: [
        { status: 'SUBMITTED', label: 'Resubmit', type: 'primary' },
        { status: 'CANCELLED', label: 'Cancel', type: 'default' },
      ],
    }
    return transitions[currentStatus] || []
  }

  const handleStatusChange = async (newStatus) => {
    try {
      setActionLoading(true)
      // TODO: Replace with actual API call
      // await api.patch(`${API_ENDPOINTS.CLAIM}/${id}/status`, { status: newStatus })

      // Mock status change
      setClaim((prev) => ({ ...prev, status: newStatus, updated_at: new Date().toISOString() }))
      message.success(`Claim status changed to ${CLAIM_STATUS_LABELS[newStatus]}`)

      // Refresh claim details
      await fetchClaimDetails()
    } catch (error) {
      message.error('Failed to update claim status')
      console.error('Error updating status:', error)
    } finally {
      setActionLoading(false)
    }
  }

  const handleEdit = () => {
    navigate(`/claims/${id}/edit`)
  }

  const handleDelete = async () => {
    try {
      setActionLoading(true)
      // TODO: Replace with actual API call
      // await api.delete(`${API_ENDPOINTS.CLAIM}/${id}`)

      message.success('Claim deleted successfully')
      navigate('/claims')
    } catch (error) {
      message.error('Failed to delete claim')
      console.error('Error deleting claim:', error)
    } finally {
      setActionLoading(false)
    }
  }

  const handleViewAttachment = (attachment) => {
    message.info(`Viewing ${attachment.name}. File viewer will be implemented later.`)
  }

  const formatFileSize = (bytes) => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  if (loading) {
    return (
      <AppLayout title="Claim Details">
        <div style={{ textAlign: 'center', padding: '50px' }}>Loading...</div>
      </AppLayout>
    )
  }

  if (!claim) {
    return (
      <AppLayout title="Claim Details">
        <div style={{ textAlign: 'center', padding: '50px' }}>Claim not found</div>
      </AppLayout>
    )
  }

  const canEditClaim = canEdit && (claim.status === 'DRAFT' || claim.status === 'REQUEST_INFO')
  const canDeleteClaim = canDelete && claim.status === 'DRAFT'
  const validTransitions = getValidTransitions(claim.status)

  const textColor =
    claim.status === 'DRAFT' || claim.status === 'PARTIALLY_APPROVED' ? '#000' : '#fff'

  return (
    <AppLayout title="Claim Details">
      <div style={{ maxWidth: 1200, margin: '0 auto', padding: '20px' }}>
        {/* Header with Actions */}
        <Card>
          <div
            style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              marginBottom: 20,
            }}
          >
            <div>
              <h2 style={{ margin: 0 }}>Claim #{claim.id.slice(0, 8)}</h2>
              <Tag
                style={{
                  color: textColor,
                  fontWeight: '500',
                  padding: '4px 12px',
                  borderRadius: '6px',
                  fontSize: '14px',
                  marginTop: 8,
                }}
              >
                {CLAIM_STATUS_LABELS[claim.status]}
              </Tag>
            </div>

            <Space>
              <Button icon={<HistoryOutlined />} onClick={() => setHistoryModalVisible(true)}>
                Status History
              </Button>

              {canEditClaim && (
                <Button type="primary" icon={<EditOutlined />} onClick={handleEdit}>
                  Edit
                </Button>
              )}

              {canDeleteClaim && (
                <Popconfirm
                  title="Delete Claim"
                  description="Are you sure you want to delete this claim? This action cannot be undone."
                  icon={<ExclamationCircleOutlined style={{ color: 'red' }} />}
                  onConfirm={handleDelete}
                  okText="Yes, Delete"
                  cancelText="Cancel"
                  okButtonProps={{ loading: actionLoading, danger: true }}
                >
                  <Button danger icon={<DeleteOutlined />} loading={actionLoading}>
                    Delete
                  </Button>
                </Popconfirm>
              )}
            </Space>
          </div>

          {/* Status Transition Actions */}
          {validTransitions.length > 0 && (
            <div style={{ marginTop: 20 }}>
              <Divider orientation="left">Available Actions</Divider>
              <Space>
                {validTransitions.map((transition) => (
                  <Button
                    key={transition.status}
                    type={transition.type}
                    icon={<ArrowRightOutlined />}
                    loading={actionLoading}
                    onClick={() => handleStatusChange(transition.status)}
                  >
                    {transition.label}
                  </Button>
                ))}
              </Space>
            </div>
          )}
        </Card>

        {/* Claim Information */}
        <Card title="Claim Information" style={{ marginTop: 20 }}>
          <Descriptions column={2} bordered>
            <Descriptions.Item label="Claim ID">{claim.id}</Descriptions.Item>
            <Descriptions.Item label="Status">
              <Tag
                style={{
                  color: textColor,
                  fontWeight: '500',
                  padding: '2px 8px',
                  borderRadius: '6px',
                }}
              >
                {CLAIM_STATUS_LABELS[claim.status]}
              </Tag>
            </Descriptions.Item>
            <Descriptions.Item label="Customer Name">{claim.customer_name}</Descriptions.Item>
            <Descriptions.Item label="Customer Email">{claim.customer_email}</Descriptions.Item>
            <Descriptions.Item label="Customer Phone">{claim.customer_phone}</Descriptions.Item>
            <Descriptions.Item label="Vehicle">{claim.vehicle_info}</Descriptions.Item>
            <Descriptions.Item label="VIN">{claim.vehicle_vin}</Descriptions.Item>
            <Descriptions.Item label="Total Cost">
              {claim.total_cost ? `${claim.total_cost.toLocaleString('vi-VN')} VND` : '0 VND'}
            </Descriptions.Item>
            <Descriptions.Item label="Created At">
              {new Date(claim.created_at).toLocaleString()}
            </Descriptions.Item>
            <Descriptions.Item label="Updated At">
              {new Date(claim.updated_at).toLocaleString()}
            </Descriptions.Item>
            <Descriptions.Item label="Description" span={2}>
              {claim.description}
            </Descriptions.Item>
          </Descriptions>
        </Card>

        {/* Attachments */}
        <Card title="Attachments" style={{ marginTop: 20 }}>
          <List
            itemLayout="horizontal"
            dataSource={attachments}
            renderItem={(attachment) => (
              <List.Item
                actions={[
                  <Button
                    type="link"
                    icon={<EyeOutlined />}
                    onClick={() => handleViewAttachment(attachment)}
                  >
                    View
                  </Button>,
                ]}
              >
                <List.Item.Meta
                  avatar={<Avatar icon={<FileOutlined />} />}
                  title={attachment.name}
                  description={
                    <Space>
                      <span>{formatFileSize(attachment.size)}</span>
                      <span>•</span>
                      <span>Uploaded by {attachment.uploaded_by}</span>
                      <span>•</span>
                      <span>{new Date(attachment.uploaded_at).toLocaleString()}</span>
                    </Space>
                  }
                />
              </List.Item>
            )}
          />
        </Card>

        {/* Status History Modal */}
        <Modal
          title="Status History"
          open={historyModalVisible}
          onCancel={() => setHistoryModalVisible(false)}
          footer={null}
          width={600}
        >
          <List
            itemLayout="horizontal"
            dataSource={statusHistory}
            renderItem={(item) => (
              <List.Item>
                <List.Item.Meta
                  title={
                    <Space>
                      <Tag
                        style={{
                          color:
                            item.status === 'DRAFT' || item.status === 'PARTIALLY_APPROVED'
                              ? '#000'
                              : '#fff',
                          fontWeight: '500',
                          padding: '2px 8px',
                          borderRadius: '6px',
                        }}
                      >
                        {CLAIM_STATUS_LABELS[item.status]}
                      </Tag>
                      <span>by {item.changed_by}</span>
                    </Space>
                  }
                  description={
                    <div>
                      <div>{new Date(item.changed_at).toLocaleString()}</div>
                      {item.notes && (
                        <div style={{ marginTop: 4, fontStyle: 'italic' }}>{item.notes}</div>
                      )}
                    </div>
                  }
                />
              </List.Item>
            )}
          />
        </Modal>
      </div>
    </AppLayout>
  )
}

export default ClaimDetail
