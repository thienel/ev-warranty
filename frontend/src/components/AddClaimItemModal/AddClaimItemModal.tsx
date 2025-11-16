import React, { useState, useEffect } from 'react'
import { Modal, Form, Input, Select, Button, Space, message } from 'antd'
import { PlusOutlined, ToolOutlined } from '@ant-design/icons'
import useHandleApiError from '@/hooks/useHandleApiError'
import { claimItems as claimItemsApi } from '@services/index'
import { CLAIM_ITEM_TYPES, CLAIM_ITEM_TYPE_LABELS } from '@constants/common-constants'
import type { CreateClaimItemRequest } from '@/types/index'

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
  faulty_part_serial: string
  issue_description: string
  type: keyof typeof CLAIM_ITEM_TYPES
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

  const handleError = useHandleApiError()

  // Reset form when modal opens/closes
  useEffect(() => {
    if (visible) {
      form.resetFields()
    }
  }, [visible, form])

  const handleSubmit = async (values: ClaimItemFormData) => {
    try {
      setLoading(true)

      const claimItemData: CreateClaimItemRequest = {
        part_category_id: values.part_category_id,
        faulty_part_serial: values.faulty_part_serial.trim(),
        issue_description: values.issue_description.trim(),
        type: CLAIM_ITEM_TYPES[values.type],
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

  return (
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
      width={700}
      destroyOnClose
    >
      <Form form={form} layout="vertical" onFinish={handleSubmit} autoComplete="off">
        <Form.Item
          label="Part Category"
          name="part_category_id"
          rules={[{ required: true, message: 'Please select a part category' }]}
        >
          <Select
            placeholder="Search and select a part category"
            showSearch
            filterOption={(input, option) => {
              const label = option?.children as unknown as string
              return label?.toLowerCase().includes(input.toLowerCase()) ?? false
            }}
            size="large"
          >
            {partCategories.map((category) => (
              <Option key={category.id} value={category.id}>
                {category.category_name}
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label="Faulty Part Serial Number"
          name="faulty_part_serial"
          rules={[
            { required: true, message: 'Please enter the faulty part serial number' },
            { min: 3, message: 'Serial number must be at least 3 characters' },
            { max: 100, message: 'Serial number cannot exceed 100 characters' },
          ]}
        >
          <Input placeholder="Enter faulty part serial number" size="large" />
        </Form.Item>

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

        <Form.Item style={{ marginBottom: 0, textAlign: 'right', marginTop: '24px' }}>
          <Space>
            <Button size="large" onClick={onCancel}>
              Cancel
            </Button>
            <Button type="primary" htmlType="submit" loading={loading} size="large">
              Add Claim Item
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  )
}

export default AddClaimItemModal
