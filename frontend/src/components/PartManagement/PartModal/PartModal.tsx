import React, { useEffect } from 'react'
import { Modal, Button, Form, Input, message, Select, Space, InputNumber } from 'antd'
import { ToolOutlined, BarcodeOutlined, DollarOutlined } from '@ant-design/icons'
import { API_ENDPOINTS } from '@constants/common-constants.js'
import { type PartModalProps, type PartFormData } from '@/types/index.js'
import api from '@services/api.js'
import useHandleApiError from '@/hooks/useHandleApiError.js'

const PartModal: React.FC<PartModalProps> = ({
  loading,
  setLoading,
  onClose,
  part = null,
  opened = false,
  isUpdate,
  partCategories,
  offices,
  partCategoriesLoading,
  officesLoading,
}) => {
  const [form] = Form.useForm<PartFormData>()
  const handleError = useHandleApiError()

  // Populate form when part prop changes or modal opens
  useEffect(() => {
    if (opened) {
      if (part && isUpdate) {
        // When editing, populate form with part data
        const formData: PartFormData = {
          serial_number: part.serial_number,
          part_name: part.part_name,
          unit_price: part.unit_price,
          category_id: part.category_id,
          office_location_id: part.office_location_id,
        }
        form.setFieldsValue(formData)
      } else {
        // When creating new, reset to default values
        form.resetFields()
      }
    }
  }, [form, part, isUpdate, opened])

  // Clear form when modal closes
  useEffect(() => {
    if (!opened) {
      form.resetFields()
    }
  }, [form, opened])

  const handleSubmit = async (values: PartFormData): Promise<void> => {
    setLoading(true)
    console.log('Submitting part data:', values)
    try {
      const payload = {
        ...values,
        office_location_id: values.office_location_id || undefined,
      }

      if (isUpdate) {
        // For update, only send part_name, unit_price, and office_location_id
        const updatePayload = {
          part_name: payload.part_name,
          unit_price: payload.unit_price,
          office_location_id: payload.office_location_id,
        }
        await api.put(`${API_ENDPOINTS.PARTS}/${part?.id}`, updatePayload)
        message.success('Part updated successfully')
      } else {
        await api.post(API_ENDPOINTS.PARTS, payload)
        message.success('Part created successfully')
      }

      onClose()
    } catch (error) {
      handleError(error as Error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Modal
      title={<Space style={{ margin: '14px 0' }}>{isUpdate ? 'Edit Part' : 'Add New Part'}</Space>}
      open={opened}
      onCancel={onClose}
      style={{ margin: 'auto' }}
      footer={null}
      width={500}
      destroyOnHidden
    >
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
        autoComplete="off"
        key={part?.id || 'new'}
      >
        <Form.Item
          label="Serial Number"
          name="serial_number"
          validateFirst
          rules={[
            { required: true, message: 'Please enter serial number' },
            { min: 1, message: 'Serial number must be at least 1 character' },
            { max: 255, message: 'Serial number cannot exceed 255 characters' },
          ]}
        >
          <Input
            placeholder="Enter serial number"
            prefix={<BarcodeOutlined />}
            size="large"
            disabled={isUpdate} // Serial number cannot be changed after creation
          />
        </Form.Item>

        <Form.Item
          label="Part Name"
          name="part_name"
          validateFirst
          rules={[
            { required: true, message: 'Please enter part name' },
            { min: 1, message: 'Part name must be at least 1 character' },
            { max: 255, message: 'Part name cannot exceed 255 characters' },
          ]}
        >
          <Input placeholder="Enter part name" prefix={<ToolOutlined />} size="large" />
        </Form.Item>

        <Form.Item
          label="Unit Price"
          name="unit_price"
          validateFirst
          rules={[
            { required: true, message: 'Please enter unit price' },
            { type: 'number', min: 0.01, message: 'Unit price must be at least 0.01' },
            { type: 'number', max: 999999999.99, message: 'Unit price cannot exceed 999999999.99' },
          ]}
        >
          <InputNumber
            placeholder="Enter unit price"
            prefix={<DollarOutlined />}
            size="large"
            style={{ width: '100%' }}
            min={0.01}
            max={999999999.99}
            step={0.01}
            precision={2}
          />
        </Form.Item>

        <Form.Item
          label="Category"
          name="category_id"
          rules={[{ required: true, message: 'Please select a category' }]}
        >
          <Select
            placeholder="Select category"
            size="large"
            showSearch
            loading={partCategoriesLoading}
            disabled={isUpdate} // Category cannot be changed after creation
            filterOption={(input, option) => {
              const label = option?.label as string
              return label?.toLowerCase().includes(input.toLowerCase())
            }}
            options={partCategories.map((cat) => ({
              label: cat.category_name,
              value: cat.id,
            }))}
          />
        </Form.Item>

        <Form.Item label="Office Location" name="office_location_id">
          <Select
            placeholder="Select office location (optional)"
            size="large"
            allowClear
            showSearch
            loading={officesLoading}
            filterOption={(input, option) => {
              const label = option?.label as string
              return label?.toLowerCase().includes(input.toLowerCase())
            }}
            options={offices.map((office) => ({
              label: office.office_name,
              value: office.id,
            }))}
          />
        </Form.Item>

        <Form.Item style={{ marginBottom: 0, textAlign: 'right', marginTop: '24px' }}>
          <Space>
            <Button size="large" onClick={onClose}>
              Cancel
            </Button>
            <Button type="primary" htmlType="submit" loading={loading} size="large">
              {isUpdate ? 'Update Part' : 'Create Part'}
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  )
}

export default PartModal
