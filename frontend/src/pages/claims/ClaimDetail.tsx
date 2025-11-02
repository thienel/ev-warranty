import React, { useState, useEffect, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import {
  Card,
  Descriptions,
  Table,
  Button,
  Space,
  Spin,
  Typography,
  Divider,
  message,
  Tag,
  Image,
  Row,
  Col,
} from 'antd'
import {
  ArrowLeftOutlined,
  UserOutlined,
  CarOutlined,
  FileTextOutlined,
  DollarOutlined,
  CalendarOutlined,
  PaperClipOutlined,
  ToolOutlined,
} from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import AppLayout from '@components/Layout/Layout'
import {
  type ClaimDetail as ClaimDetailType,
  type ClaimItem,
  type ClaimAttachment,
  type Customer,
  type VehicleDetail,
  type VehicleModel,
} from '@/types/index'
import {
  claims as claimsApi,
  claimItems as claimItemsApi,
  claimAttachments as claimAttachmentsApi,
  customersApi,
  vehiclesApi,
  vehicleModelsApi,
} from '@services/index'
import useHandleApiError from '@/hooks/useHandleApiError'
import {
  CLAIM_STATUS_LABELS,
  CLAIM_ITEM_STATUS_LABELS,
  ATTACHMENT_TYPE_LABELS,
} from '@constants/common-constants'
import { getClaimsBasePath } from '@/utils/navigationHelpers'
import './ClaimDetail.less'

const { Title, Text, Paragraph } = Typography

const ClaimDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const handleError = useHandleApiError()

  const [claim, setClaim] = useState<ClaimDetailType | null>(null)
  const [customer, setCustomer] = useState<Customer | null>(null)
  const [vehicle, setVehicle] = useState<VehicleDetail | null>(null)
  const [claimItems, setClaimItems] = useState<ClaimItem[]>([])
  const [attachments, setAttachments] = useState<ClaimAttachment[]>([])

  const [claimLoading, setClaimLoading] = useState(false)
  const [customerLoading, setCustomerLoading] = useState(false)
  const [vehicleLoading, setVehicleLoading] = useState(false)
  const [itemsLoading, setItemsLoading] = useState(false)
  const [attachmentsLoading, setAttachmentsLoading] = useState(false)

  // Fetch claim details
  const fetchClaim = useCallback(async () => {
    if (!id) return

    try {
      setClaimLoading(true)
      const response = await claimsApi.getById(id)
      let claimData = response.data

      if (claimData && typeof claimData === 'object' && 'data' in claimData) {
        claimData = (claimData as { data: unknown }).data as ClaimDetailType
      }

      setClaim(claimData as ClaimDetailType)
    } catch (error) {
      handleError(error as Error)
      message.error('Failed to load claim details')
    } finally {
      setClaimLoading(false)
    }
  }, [id, handleError])

  // Fetch customer details
  const fetchCustomer = useCallback(
    async (customerId: string) => {
      try {
        setCustomerLoading(true)
        const customerResponse = await customersApi.getById(customerId)
        let customerData = customerResponse.data
        if (customerData && typeof customerData === 'object' && 'data' in customerData) {
          customerData = (customerData as { data: unknown }).data as Customer
        }
        setCustomer(customerData as Customer)
      } catch (error) {
        handleError(error as Error)
      } finally {
        setCustomerLoading(false)
      }
    },
    [handleError],
  )

  // Fetch vehicle details with model info
  const fetchVehicle = useCallback(
    async (vehicleId: string) => {
      try {
        setVehicleLoading(true)
        const vehicleResponse = await vehiclesApi.getById(vehicleId)
        let vehicleData = vehicleResponse.data
        if (vehicleData && typeof vehicleData === 'object' && 'data' in vehicleData) {
          vehicleData = (vehicleData as { data: unknown }).data as VehicleDetail
        }

        const vehicleInfo = vehicleData as VehicleDetail

        // Fetch vehicle model if model_id exists
        if (vehicleInfo.model_id) {
          try {
            const modelResponse = await vehicleModelsApi.getById(vehicleInfo.model_id)
            let modelData: unknown = modelResponse.data
            if (modelData && typeof modelData === 'object' && 'data' in modelData) {
              modelData = (modelData as { data: VehicleModel }).data
            }
            vehicleInfo.model = modelData as VehicleModel
          } catch (error) {
            console.error('Failed to fetch vehicle model:', error)
          }
        }

        setVehicle(vehicleInfo)
      } catch (error) {
        handleError(error as Error)
      } finally {
        setVehicleLoading(false)
      }
    },
    [handleError],
  )

  // Fetch claim items
  const fetchClaimItems = useCallback(async () => {
    if (!id) return

    try {
      setItemsLoading(true)
      const response = await claimItemsApi.getByClaimId(id)
      const itemsData: unknown = response.data

      // Backend returns: { data: ClaimItem[] } not { data: { items: ClaimItem[] } }
      if (itemsData && typeof itemsData === 'object' && 'data' in itemsData) {
        const nestedData = (itemsData as { data: unknown }).data
        if (Array.isArray(nestedData)) {
          setClaimItems(nestedData)
        } else {
          setClaimItems([])
        }
      } else if (Array.isArray(itemsData)) {
        setClaimItems(itemsData)
      } else {
        setClaimItems([])
      }
    } catch (error) {
      handleError(error as Error)
      setClaimItems([])
    } finally {
      setItemsLoading(false)
    }
  }, [id, handleError])

  // Fetch claim attachments
  const fetchAttachments = useCallback(async () => {
    if (!id) return

    try {
      setAttachmentsLoading(true)
      const response = await claimAttachmentsApi.getByClaimId(id)
      const attachmentsData: unknown = response.data

      // Backend returns: { data: ClaimAttachment[] } not { data: { attachments: ClaimAttachment[] } }
      if (attachmentsData && typeof attachmentsData === 'object' && 'data' in attachmentsData) {
        const nestedData = (attachmentsData as { data: unknown }).data
        if (Array.isArray(nestedData)) {
          setAttachments(nestedData)
        } else {
          setAttachments([])
        }
      } else if (Array.isArray(attachmentsData)) {
        setAttachments(attachmentsData)
      } else {
        setAttachments([])
      }
    } catch (error) {
      handleError(error as Error)
      setAttachments([])
    } finally {
      setAttachmentsLoading(false)
    }
  }, [id, handleError])

  // Initial fetch claim
  useEffect(() => {
    if (id) {
      fetchClaim()
      fetchClaimItems()
      fetchAttachments()
    }
  }, [id, fetchClaim, fetchClaimItems, fetchAttachments])

  // Fetch customer and vehicle when claim is loaded
  useEffect(() => {
    if (claim) {
      if (claim.customer_id) {
        fetchCustomer(claim.customer_id)
      }
      if (claim.vehicle_id) {
        fetchVehicle(claim.vehicle_id)
      }
    }
  }, [claim, fetchCustomer, fetchVehicle])

  const handleBack = () => {
    const location = window.location.pathname
    navigate(getClaimsBasePath(location))
  }

  // Get status color
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

  const getItemStatusColor = (status: string): string => {
    const colors: Record<string, string> = {
      PENDING: 'default',
      APPROVED: 'green',
      REJECTED: 'red',
      COMPLETED: 'purple',
    }
    return colors[status] || 'default'
  }

  // Claim items table columns
  const claimItemColumns: ColumnsType<ClaimItem> = [
    {
      title: 'Item ID',
      dataIndex: 'id',
      key: 'id',
      width: '15%',
      render: (id: string) => <Text code>{id.slice(0, 8)}...</Text>,
    },
    {
      title: 'Part Category',
      dataIndex: 'part_category_id',
      key: 'part_category_id',
      width: '15%',
    },
    {
      title: 'Faulty Part',
      dataIndex: 'faulty_part_id',
      key: 'faulty_part_id',
      width: '12%',
      render: (id: string) => <Text code>{id.slice(0, 8)}...</Text>,
    },
    {
      title: 'Issue Description',
      dataIndex: 'issue_description',
      key: 'issue_description',
      width: '25%',
      ellipsis: true,
    },
    {
      title: 'Type',
      dataIndex: 'type',
      key: 'type',
      width: '10%',
      render: (type: string) => <Tag color="blue">{type}</Tag>,
    },
    {
      title: 'Cost',
      dataIndex: 'cost',
      key: 'cost',
      width: '10%',
      align: 'right',
      render: (cost: number) => (
        <Text strong style={{ color: '#52c41a' }}>
          ${cost.toLocaleString()}
        </Text>
      ),
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      width: '13%',
      render: (status: string) => (
        <Tag color={getItemStatusColor(status)}>
          {CLAIM_ITEM_STATUS_LABELS[status as keyof typeof CLAIM_ITEM_STATUS_LABELS] || status}
        </Tag>
      ),
    },
  ]

  if (claimLoading) {
    return (
      <AppLayout title="Claim Details">
        <div style={{ textAlign: 'center', padding: '50px' }}>
          <Spin size="large" />
        </div>
      </AppLayout>
    )
  }

  if (!claim) {
    return (
      <AppLayout title="Claim Details">
        <Card>
          <Text>Claim not found</Text>
        </Card>
      </AppLayout>
    )
  }

  return (
    <AppLayout title="Claim Details">
      <div className="claim-detail-page">
        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          {/* Header */}
          <Card>
            <Space style={{ width: '100%', justifyContent: 'space-between' }}>
              <Button icon={<ArrowLeftOutlined />} onClick={handleBack}>
                Back to Claims
              </Button>
              <Space>
                <Tag
                  color={getStatusColor(claim.status)}
                  style={{ fontSize: '14px', padding: '4px 12px' }}
                >
                  {CLAIM_STATUS_LABELS[claim.status] || claim.status}
                </Tag>
              </Space>
            </Space>
          </Card>

          {/* Claim Information */}
          <Card
            title={
              <Title level={4}>
                <FileTextOutlined /> Claim Information
              </Title>
            }
          >
            <Descriptions bordered column={2}>
              <Descriptions.Item label="Claim ID" span={2}>
                <Text code>{claim.id}</Text>
              </Descriptions.Item>
              <Descriptions.Item label="Description" span={2}>
                <Paragraph style={{ marginBottom: 0, whiteSpace: 'pre-wrap' }}>
                  {claim.description}
                </Paragraph>
              </Descriptions.Item>
              <Descriptions.Item label="Total Cost">
                <Space>
                  <DollarOutlined style={{ color: '#52c41a' }} />
                  <Text strong style={{ fontSize: '16px', color: '#52c41a' }}>
                    ${claim.total_cost.toLocaleString()}
                  </Text>
                </Space>
              </Descriptions.Item>
              <Descriptions.Item label="Status">
                <Tag color={getStatusColor(claim.status)}>
                  {CLAIM_STATUS_LABELS[claim.status] || claim.status}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="Created At">
                <Space>
                  <CalendarOutlined />
                  {claim.created_at ? new Date(claim.created_at).toLocaleString() : 'N/A'}
                </Space>
              </Descriptions.Item>
              <Descriptions.Item label="Updated At">
                <Space>
                  <CalendarOutlined />
                  {claim.updated_at ? new Date(claim.updated_at).toLocaleString() : 'N/A'}
                </Space>
              </Descriptions.Item>
              {claim.approved_by && (
                <Descriptions.Item label="Approved By" span={2}>
                  <Text code>{claim.approved_by}</Text>
                </Descriptions.Item>
              )}
            </Descriptions>
          </Card>

          <Row gutter={16}>
            {/* Customer Information */}
            <Col xs={24} lg={12}>
              <Card
                title={
                  <Title level={4}>
                    <UserOutlined /> Customer Information
                  </Title>
                }
                loading={customerLoading}
              >
                {customer ? (
                  <Descriptions bordered column={1}>
                    <Descriptions.Item label="Customer ID">
                      <Text code>{customer.id}</Text>
                    </Descriptions.Item>
                    <Descriptions.Item label="Name">
                      <Text strong>
                        {customer.full_name || `${customer.first_name} ${customer.last_name}`}
                      </Text>
                    </Descriptions.Item>
                    {customer.email && (
                      <Descriptions.Item label="Email">{customer.email}</Descriptions.Item>
                    )}
                    {customer.phone_number && (
                      <Descriptions.Item label="Phone">{customer.phone_number}</Descriptions.Item>
                    )}
                    {customer.address && (
                      <Descriptions.Item label="Address">{customer.address}</Descriptions.Item>
                    )}
                  </Descriptions>
                ) : !customerLoading ? (
                  <Text type="secondary">Customer information not available</Text>
                ) : null}
              </Card>
            </Col>

            {/* Vehicle Information */}
            <Col xs={24} lg={12}>
              <Card
                title={
                  <Title level={4}>
                    <CarOutlined /> Vehicle Information
                  </Title>
                }
                loading={vehicleLoading}
              >
                {vehicle ? (
                  <Descriptions bordered column={1}>
                    <Descriptions.Item label="Vehicle ID">
                      <Text code>{vehicle.id}</Text>
                    </Descriptions.Item>
                    <Descriptions.Item label="VIN">
                      <Text strong>{vehicle.vin}</Text>
                    </Descriptions.Item>
                    {vehicle.license_plate && (
                      <Descriptions.Item label="License Plate">
                        {vehicle.license_plate}
                      </Descriptions.Item>
                    )}
                    {vehicle.model && (
                      <>
                        <Descriptions.Item label="Brand">{vehicle.model.brand}</Descriptions.Item>
                        <Descriptions.Item label="Model">
                          {vehicle.model.model_name}
                        </Descriptions.Item>
                        <Descriptions.Item label="Year">{vehicle.model.year}</Descriptions.Item>
                      </>
                    )}
                    {vehicle.purchase_date && (
                      <Descriptions.Item label="Purchase Date">
                        {new Date(vehicle.purchase_date).toLocaleDateString()}
                      </Descriptions.Item>
                    )}
                  </Descriptions>
                ) : !vehicleLoading ? (
                  <Text type="secondary">Vehicle information not available</Text>
                ) : null}
              </Card>
            </Col>
          </Row>

          {/* Claim Items */}
          <Card
            title={
              <Title level={4}>
                <ToolOutlined /> Claim Items
              </Title>
            }
            loading={itemsLoading}
          >
            <Table
              dataSource={claimItems}
              columns={claimItemColumns}
              rowKey="id"
              pagination={false}
              locale={{ emptyText: 'No claim items found' }}
            />
            {claimItems.length > 0 && <Divider />}
          </Card>

          {/* Attachments */}
          <Card
            title={
              <Title level={4}>
                <PaperClipOutlined /> Attachments
              </Title>
            }
            loading={attachmentsLoading}
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
        </Space>
      </div>
    </AppLayout>
  )
}

export default ClaimDetail
