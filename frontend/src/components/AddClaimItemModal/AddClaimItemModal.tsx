import React, { useState, useEffect } from 'react'
import {
  Modal,
  Form,
  Input,
  Select,
  Button,
  Space,
  Typography,
  Card,
  Descriptions,
  message,
  InputNumber,
} from 'antd'
import { PlusOutlined, ToolOutlined, DollarOutlined } from '@ant-design/icons'
import PartSearchModal from '@/components/PartSearchModal/PartSearchModal'
import useHandleApiError from '@/hooks/useHandleApiError'
import { claimItems as claimItemsApi } from '@services/index'
import { CLAIM_ITEM_TYPES, CLAIM_ITEM_TYPE_LABELS } from '@constants/common-constants'
import type { Part, CreateClaimItemRequest } from '@/types/index'

const { Text } = Typography
const { Option } = Select
const { TextArea } = Input

interface AddClaimItemModalProps {
  visible: boolean
  onCancel: () => void
  onSuccess: () => void
  claimId: string
  partCategories: Array<{ id: string; category_name: string }>
}

interface ClaimItemFormData {
  part_category_id: string
  faulty_part_id: string
  issue_description: string
  type: keyof typeof CLAIM_ITEM_TYPES
  cost: number
}

const AddClaimItemModal: React.FC<AddClaimItemModalProps> = ({
  visible,
  onCancel,
  onSuccess,
  claimId,
  partCategories,
}) => {
  const [form] = Form.useForm<ClaimItemFormData>()
  const [loading, setLoading] = useState(false)
  const [selectedPart, setSelectedPart] = useState<Part | null>(null)
  const [partSearchVisible, setPartSearchVisible] = useState(false)

  const handleError = useHandleApiError()

  // Reset form when modal opens/closes
  useEffect(() => {
    if (visible) {
      form.resetFields()
      setSelectedPart(null)
    }
  }, [visible, form])

  // Update form when a part is selected
  useEffect(() => {
    if (selectedPart) {
      form.setFieldsValue({
        part_category_id: selectedPart.category_id,
        faulty_part_id: selectedPart.id,
        cost: selectedPart.unit_price,
      })
    }
  }, [selectedPart, form])

  const handlePartSelect = (part: Part) => {
    setSelectedPart(part)
    setPartSearchVisible(false)
    message.success('Part selected successfully!')
  }

  const handleSubmit = async (values: ClaimItemFormData) => {
    if (!selectedPart) {
      message.error('Please select a part first')
      return
    }

    try {
      setLoading(true)

      const claimItemData: CreateClaimItemRequest = {
        part_category_id: values.part_category_id,
        faulty_part_id: values.faulty_part_id,
        issue_description: values.issue_description.trim(),
        type: CLAIM_ITEM_TYPES[values.type],
        cost: values.cost,
      }

      await claimItemsApi.create(claimId, claimItemData)

      message.success('Claim item added successfully!')
      onSuccess()
      onCancel() // Close modal
    } catch (error) {
      handleError(error as Error)
    } finally {
      setLoading(false)
    }
  }

  const getCategoryName = (categoryId: string): string => {
    const category = partCategories.find((cat) => cat.id === categoryId)
    return category?.category_name || 'Unknown Category'
  }

  const handleRemovePart = () => {
    setSelectedPart(null)
    form.setFieldsValue({
      part_category_id: undefined,
      faulty_part_id: undefined,
      cost: undefined,
    })
  }

  return (
    <>
      <Modal
        title={
          <Space>
            <PlusOutlined />
            Add New Claim Item
          </Space>
        }
        open={visible}
        onCancel={onCancel}
        footer={null}
        width={800}
        destroyOnClose
      >
        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          {/* Part Selection Section */}
          <Card size="small" title="Select Faulty Part">
            {selectedPart ? (
              <Space direction="vertical" size="small" style={{ width: '100%' }}>
                <Descriptions bordered column={2} size="small">
                  <Descriptions.Item label="Serial Number" span={1}>
                    <Text strong>{selectedPart.serial_number}</Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="Part Name" span={1}>
                    {selectedPart.part_name}
                  </Descriptions.Item>
                  <Descriptions.Item label="Category" span={1}>
                    {getCategoryName(selectedPart.category_id)}
                  </Descriptions.Item>
                  <Descriptions.Item label="Unit Price" span={1}>
                    <Text strong style={{ color: '#52c41a' }}>
                      {selectedPart.unit_price?.toLocaleString('vi-VN', {
                        style: 'currency',
                        currency: 'VND',
                      })}
                    </Text>
                  </Descriptions.Item>
                </Descriptions>
                <Space>
                  <Button
                    type="default"
                    onClick={() => setPartSearchVisible(true)}
                    icon={<ToolOutlined />}
                  >
                    Change Part
                  </Button>
                  <Button type="text" danger onClick={handleRemovePart}>
                    Remove Part
                  </Button>
                </Space>
              </Space>
            ) : (
              <Space direction="vertical" style={{ width: '100%' }}>
                <Text type="secondary">
                  No part selected. Click the button below to search and select a part.
                </Text>
                <Button
                  type="primary"
                  onClick={() => setPartSearchVisible(true)}
                  icon={<ToolOutlined />}
                >
                  Search & Select Part
                </Button>
              </Space>
            )}
          </Card>

          {/* Claim Item Form */}
          {selectedPart && (
            <Card size="small" title="Claim Item Details">
              <Form form={form} layout="vertical" onFinish={handleSubmit} autoComplete="off">
                <Form.Item
                  label="Issue Description"
                  name="issue_description"
                  rules={[
                    { required: true, message: 'Please describe the issue' },
                    { min: 10, message: 'Description must be at least 10 characters' },
                    { max: 1000, message: 'Description cannot exceed 1000 characters' },
                  ]}
                >
                  <TextArea placeholder="Describe the issue with this part..." rows={4} />
                </Form.Item>

                <Form.Item
                  label="Repair Type"
                  name="type"
                  rules={[{ required: true, message: 'Please select repair type' }]}
                >
                  <Select placeholder="Select repair type" size="large">
                    {Object.entries(CLAIM_ITEM_TYPES).map(([key, value]) => (
                      <Option key={key} value={value}>
                        <Space>
                          <ToolOutlined />
                          {CLAIM_ITEM_TYPE_LABELS[value]}
                        </Space>
                      </Option>
                    ))}
                  </Select>
                </Form.Item>

                <Form.Item
                  label="Cost"
                  name="cost"
                  rules={[
                    { required: true, message: 'Please enter the cost' },
                    { type: 'number', min: 0.01, message: 'Cost must be at least 0.01' },
                  ]}
                >
                  <InputNumber
                    placeholder="Enter cost"
                    prefix={<DollarOutlined />}
                    size="large"
                    style={{ width: '100%' }}
                    min={0.01}
                    step={0.01}
                    precision={2}
                    formatter={(value) => `${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')}
                  />
                </Form.Item>

                {/* Hidden fields for part info */}
                <Form.Item name="part_category_id" hidden>
                  <Input />
                </Form.Item>
                <Form.Item name="faulty_part_id" hidden>
                  <Input />
                </Form.Item>

                <Form.Item style={{ marginBottom: 0, textAlign: 'right', marginTop: '24px' }}>
                  <Space>
                    <Button size="large" onClick={onCancel}>
                      Cancel
                    </Button>
                    <Button
                      type="primary"
                      htmlType="submit"
                      loading={loading}
                      size="large"
                      disabled={!selectedPart}
                    >
                      Add Claim Item
                    </Button>
                  </Space>
                </Form.Item>
              </Form>
            </Card>
          )}
        </Space>
      </Modal>

      {/* Part Search Modal */}
      <PartSearchModal
        visible={partSearchVisible}
        onCancel={() => setPartSearchVisible(false)}
        onSelectPart={handlePartSelect}
        title="Search and Select Faulty Part"
      />
    </>
  )
}

export default AddClaimItemModal
