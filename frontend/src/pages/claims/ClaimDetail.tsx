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
  type PartCategory,
  type Part,
} from '@/types/index'
import {
  claims as claimsApi,
  claimItems as claimItemsApi,
  claimAttachments as claimAttachmentsApi,
  customersApi,
  vehiclesApi,
  vehicleModelsApi,
  partCategoriesApi,
  partsApi,
} from '@services/index'
import useHandleApiError from '@/hooks/useHandleApiError'
import {
  CLAIM_STATUS_LABELS,
  CLAIM_ITEM_STATUS_LABELS,
  ATTACHMENT_TYPE_LABELS,
  CLAIM_ITEM_TYPE_LABELS,
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
  const [partCategories, setPartCategories] = useState<PartCategory[]>([])
  const [parts, setParts] = useState<Part[]>([])

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

  // Fetch parts for claim items
  const fetchPartsForClaimItems = useCallback(
    async (claimItems: ClaimItem[]) => {
      if (claimItems.length === 0) return

      try {
        // Get unique part IDs from claim items
        const uniquePartIds = Array.from(new Set(claimItems.map((item) => item.faulty_part_id)))

        // Fetch each part individually
        const partPromises = uniquePartIds.map((partId) => partsApi.getById(partId))
        const partResponses = await Promise.allSettled(partPromises)

        const fetchedParts: Part[] = []
        let hasFailures = false

        partResponses.forEach((response, index) => {
          if (response.status === 'fulfilled') {
            let partData = response.value.data
            if (partData && typeof partData === 'object' && 'data' in partData) {
              partData = (partData as { data: unknown }).data as Part
            }
            fetchedParts.push(partData as Part)
          } else {
            console.warn(`Part ${uniquePartIds[index]} not found, it may have been deleted`)
            hasFailures = true
          }
        })

        // If we couldn't fetch some parts individually, fall back to fetching all parts
        if (hasFailures && fetchedParts.length < uniquePartIds.length) {
          console.warn('Some parts not found individually, falling back to fetch all parts')
          try {
            const response = await partsApi.getAll()
            let partsData = response.data
            if (partsData && typeof partsData === 'object' && 'data' in partsData) {
              partsData = (partsData as { data: unknown }).data as Part[]
            }
            setParts(partsData as Part[])
            return
          } catch (fallbackError) {
            console.error('Failed to fetch all parts as fallback:', fallbackError)
          }
        }

        setParts(fetchedParts)
      } catch (error) {
        handleError(error as Error)
        setParts([])
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
      let items: ClaimItem[] = []
      if (itemsData && typeof itemsData === 'object' && 'data' in itemsData) {
        const nestedData = (itemsData as { data: unknown }).data
        if (Array.isArray(nestedData)) {
          items = nestedData
          setClaimItems(nestedData)
        } else {
          setClaimItems([])
        }
      } else if (Array.isArray(itemsData)) {
        items = itemsData
        setClaimItems(itemsData)
      } else {
        setClaimItems([])
      }

      // Fetch parts for the claim items
      if (items.length > 0) {
        await fetchPartsForClaimItems(items)
      }
    } catch (error) {
      handleError(error as Error)
      setClaimItems([])
    } finally {
      setItemsLoading(false)
    }
  }, [id, handleError, fetchPartsForClaimItems])

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

  // Fetch part categories
  const fetchPartCategories = useCallback(async () => {
    try {
      const response = await partCategoriesApi.getAll()
      let categoriesData = response.data
      if (categoriesData && typeof categoriesData === 'object' && 'data' in categoriesData) {
        categoriesData = (categoriesData as { data: unknown }).data as PartCategory[]
      }
      setPartCategories(categoriesData as PartCategory[])
    } catch (error) {
      handleError(error as Error)
      setPartCategories([])
    }
  }, [handleError])

  // Initial fetch claim
  useEffect(() => {
    if (id) {
      fetchClaim()
      fetchClaimItems()
      fetchAttachments()
    }
    // Fetch part categories for displaying category names
    fetchPartCategories()
  }, [id, fetchClaim, fetchClaimItems, fetchAttachments, fetchPartCategories])

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

  // Get part category name by ID
  const getPartCategoryName = (categoryId: string): string => {
    const category = partCategories.find((cat) => cat.id === categoryId)
    return category?.category_name || `Category ${categoryId.slice(0, 8)}...`
  }

  // Get part serial number by ID
  const getPartSerialNumber = (partId: string): string => {
    const part = parts.find((part) => part.id === partId)
    if (part?.serial_number) {
      return part.serial_number
    }
    // If part not found, show a more user-friendly message
    return 'N/A (Part not found)'
  }

  // Get part name by ID
  const getPartName = (partId: string): string => {
    const part = parts.find((part) => part.id === partId)
    if (part?.part_name) {
      return part.part_name
    }
    // If part not found, show a more user-friendly message
    return 'N/A (Part not found)'
  }

  // Claim items table columns
  const claimItemColumns: ColumnsType<ClaimItem> = [
    {
      title: 'Serial no.',
      dataIndex: 'faulty_part_id',
      key: 'serial_number',
      width: '15%',
      render: (_, record: ClaimItem) => {
        const serialNumber = getPartSerialNumber(record.faulty_part_id)
        const isNotFound = serialNumber.includes('not found')
        return (
          <Text type={isNotFound ? 'secondary' : undefined} italic={isNotFound}>
            {serialNumber}
          </Text>
        )
      },
    },
    {
      title: 'Part Category',
      dataIndex: 'part_category_id',
      key: 'part_category_id',
      width: '15%',
      render: (_, record: ClaimItem) => <Text>{getPartCategoryName(record.part_category_id)}</Text>,
    },
    {
      title: 'Faulty Part',
      dataIndex: 'faulty_part_id',
      key: 'faulty_part_id',
      width: '12%',
      render: (_, record: ClaimItem) => {
        const partName = getPartName(record.faulty_part_id)
        const isNotFound = partName.includes('not found')
        return (
          <Text type={isNotFound ? 'secondary' : undefined} italic={isNotFound}>
            {partName}
          </Text>
        )
      },
    },
    {
      title: 'Issue Description',
      dataIndex: 'issue_description',
      key: 'issue_description',
      width: '25%',
      render: (text: string) => (
        <div
          style={{
            whiteSpace: 'normal',
            wordBreak: 'break-word',
          }}
        >
          {text}
        </div>
      ),
    },
    {
      title: 'Type',
      dataIndex: 'type',
      key: 'type',
      width: '10%',
      render: (type: string) => (
        <>{CLAIM_ITEM_TYPE_LABELS[type as keyof typeof CLAIM_ITEM_TYPE_LABELS] || type}</>
      ),
    },
    {
      title: 'Cost',
      dataIndex: 'cost',
      key: 'cost',
      width: '10%',
      align: 'right',
      render: (cost: number) => (
        <Text strong style={{ color: '#52c41a' }}>
          {cost.toLocaleString('vi-VN', {
            style: 'currency',
            currency: 'VND',
          })}
        </Text>
      ),
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      width: '13%',
      render: (status: string) => (
        <>{CLAIM_ITEM_STATUS_LABELS[status as keyof typeof CLAIM_ITEM_STATUS_LABELS] || status}</>
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
          <Space style={{ width: '100%', justifyContent: 'space-between' }}>
            <Button icon={<ArrowLeftOutlined />} onClick={handleBack}>
              Back to Claims
            </Button>
          </Space>

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
                <Text>{claim.id}</Text>
              </Descriptions.Item>
              <Descriptions.Item label="Description" span={2}>
                <Paragraph style={{ marginBottom: 0, whiteSpace: 'pre-wrap' }}>
                  {claim.description}
                </Paragraph>
              </Descriptions.Item>
              <Descriptions.Item label="Total Cost">
                <Space>
                  <Text strong style={{ fontSize: '16px', color: '#52c41a' }}>
                    {claim.total_cost.toLocaleString('vi-VN', {
                      style: 'currency',
                      currency: 'VND',
                    })}
                  </Text>
                </Space>
              </Descriptions.Item>
              <Descriptions.Item label="Status">
                <Text color={getStatusColor(claim.status)}>
                  {CLAIM_STATUS_LABELS[claim.status] || claim.status}
                </Text>
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
                  <Text>{claim.approved_by}</Text>
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
