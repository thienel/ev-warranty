import React, { useEffect } from 'react'
import { Modal, Button, Form, Input, message, Space, InputNumber, Select } from 'antd'
import { SafetyOutlined, ClockCircleOutlined, CarOutlined } from '@ant-design/icons'
import { API_ENDPOINTS } from '@constants/common-constants.js'
import { type WarrantyPolicyModalProps, type WarrantyPolicyFormData } from '@/types/index.js'
import api from '@services/api.js'
import useHandleApiError from '@/hooks/useHandleApiError.js'

const { TextArea } = Input
const { Option } = Select

const WarrantyPolicyModal: React.FC<WarrantyPolicyModalProps> = ({
  loading,
  setLoading,
  onClose,
  policy = null,
  opened = false,
  isUpdate,
  vehicleModels,
}) => {
  const [form] = Form.useForm<WarrantyPolicyFormData>()
  const handleError = useHandleApiError()

  useEffect(() => {
    if (opened) {
      if (policy && isUpdate) {
        const formData: WarrantyPolicyFormData = {
          policy_name: policy.policy_name,
          warranty_duration_months: policy.warranty_duration_months,
          kilometer_limit: policy.kilometer_limit,
          terms_and_conditions: policy.terms_and_conditions,
          vehicle_model_id: policy.vehicle_model_id,
        }
        form.setFieldsValue(formData)
      } else {
        form.resetFields()
      }
    }
  }, [form, policy, isUpdate, opened])

  useEffect(() => {
    if (!opened) {
      form.resetFields()
    }
  }, [form, opened])

  const handleSubmit = async (values: WarrantyPolicyFormData): Promise<void> => {
    setLoading(true)
    try {
      const payload = {
        ...values,
        kilometer_limit: values.kilometer_limit || undefined,
        vehicle_model_id: values.vehicle_model_id || undefined,
      }

      if (isUpdate) {
        await api.put(`${API_ENDPOINTS.WARRANTY_POLICIES}/${policy?.id}`, payload)
        message.success('Warranty policy updated successfully')
      } else {
        await api.post(API_ENDPOINTS.WARRANTY_POLICIES, payload)
        message.success('Warranty policy created successfully')
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
      title={
        <Space style={{ margin: '14px 0' }}>
          {isUpdate ? 'Edit Warranty Policy' : 'Create New Warranty Policy'}
        </Space>
      }
      open={opened}
      onCancel={onClose}
      style={{ margin: 'auto' }}
      footer={null}
      width={600}
      destroyOnHidden
    >
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
        autoComplete="off"
        key={policy?.id || 'new'}
      >
        <Form.Item
          label="Policy Name"
          name="policy_name"
          validateFirst
          rules={[
            { required: true, message: 'Please enter policy name' },
            { min: 1, message: 'Policy name must be at least 1 character' },
            { max: 255, message: 'Policy name cannot exceed 255 characters' },
          ]}
        >
          <Input placeholder="Enter policy name" prefix={<SafetyOutlined />} size="large" />
        </Form.Item>

        <Form.Item
          label="Vehicle Model (Optional)"
          name="vehicle_model_id"
        >
          <Select
            placeholder="Select a vehicle model"
            size="large"
            showSearch
            allowClear
            optionFilterProp="children"
            filterOption={(input, option) => {
              const label = option?.children as unknown as string
              return label?.toLowerCase().includes(input.toLowerCase()) ?? false
            }}
          >
            {vehicleModels.map((model) => (
              <Option key={model.id} value={model.id}>
                {model.brand} {model.model_name} ({model.year})
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label="Warranty Duration (Months)"
          name="warranty_duration_months"
          validateFirst
          rules={[
            { required: true, message: 'Please enter warranty duration' },
            { type: 'number', min: 1, message: 'Duration must be at least 1 month' },
            { type: 'number', max: 600, message: 'Duration cannot exceed 600 months' },
          ]}
        >
          <InputNumber
            placeholder="Enter warranty duration in months"
            prefix={<ClockCircleOutlined />}
            size="large"
            style={{ width: '100%' }}
            min={1}
            max={600}
          />
        </Form.Item>

        <Form.Item
          label="Kilometer Limit (Optional)"
          name="kilometer_limit"
          validateFirst
          rules={[
            { type: 'number', min: 1, message: 'Kilometer limit must be at least 1' },
            { type: 'number', max: 9999999, message: 'Kilometer limit cannot exceed 9,999,999' },
          ]}
        >
          <InputNumber
            placeholder="Enter kilometer limit (optional)"
            prefix={<CarOutlined />}
            size="large"
            style={{ width: '100%' }}
            min={1}
            max={9999999}
          />
        </Form.Item>

        <Form.Item
          label="Terms and Conditions"
          name="terms_and_conditions"
          validateFirst
          rules={[
            { required: true, message: 'Please enter terms and conditions' },
            { min: 1, message: 'Terms and conditions must be at least 1 character' },
            { max: 5000, message: 'Terms and conditions cannot exceed 5000 characters' },
          ]}
        >
          <TextArea
            placeholder="Enter warranty terms and conditions"
            rows={6}
            size="large"
            showCount
            maxLength={5000}
          />
        </Form.Item>

        <Form.Item style={{ marginBottom: 0, textAlign: 'right', marginTop: '24px' }}>
          <Space>
            <Button size="large" onClick={onClose}>
              Cancel
            </Button>
            <Button type="primary" htmlType="submit" loading={loading} size="large">
              {isUpdate ? 'Update Policy' : 'Create Policy'}
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  )
}

export default WarrantyPolicyModal
