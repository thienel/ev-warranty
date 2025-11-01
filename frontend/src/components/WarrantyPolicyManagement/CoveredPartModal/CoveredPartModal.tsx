import React, { useEffect } from 'react'
import { Modal, Button, Form, Input, message, Select, Space } from 'antd'
import { API_ENDPOINTS } from '@constants/common-constants.js'
import {
  type PolicyCoveragePartModalProps,
  type PolicyCoveragePartFormData,
} from '@/types/index.js'
import api from '@services/api.js'
import useHandleApiError from '@/hooks/useHandleApiError.js'

const { TextArea } = Input

const CoveredPartModal: React.FC<PolicyCoveragePartModalProps> = ({
  loading,
  setLoading,
  onClose,
  coveragePart = null,
  opened = false,
  isUpdate,
  policyId,
  partCategories,
  partCategoriesLoading,
}) => {
  const [form] = Form.useForm<PolicyCoveragePartFormData>()
  const handleError = useHandleApiError()

  useEffect(() => {
    if (opened) {
      if (coveragePart && isUpdate) {
        const formData: PolicyCoveragePartFormData = {
          policy_id: coveragePart.policy_id,
          part_category_id: coveragePart.part_category_id,
          coverage_conditions: coveragePart.coverage_conditions,
        }
        form.setFieldsValue(formData)
      } else {
        // Pre-fill policy_id when creating new
        form.setFieldsValue({ policy_id: policyId })
      }
    }
  }, [form, coveragePart, isUpdate, opened, policyId])

  useEffect(() => {
    if (!opened) {
      form.resetFields()
    }
  }, [form, opened])

  const handleSubmit = async (values: PolicyCoveragePartFormData): Promise<void> => {
    setLoading(true)
    console.log('Submitting policy coverage part data:', values)
    try {
      const payload = {
        ...values,
        coverage_conditions: values.coverage_conditions || undefined,
      }

      if (isUpdate) {
        // For update, only send coverage_conditions
        const updatePayload = {
          coverage_conditions: payload.coverage_conditions,
        }
        await api.put(`${API_ENDPOINTS.POLICY_COVERAGE_PARTS}/${coveragePart?.id}`, updatePayload)
        message.success('Covered part updated successfully')
      } else {
        await api.post(API_ENDPOINTS.POLICY_COVERAGE_PARTS, payload)
        message.success('Covered part added successfully')
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
          {isUpdate ? 'Edit Covered Part' : 'Add Covered Part'}
        </Space>
      }
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
        key={coveragePart?.id || 'new'}
      >
        <Form.Item name="policy_id" hidden>
          <Input />
        </Form.Item>

        <Form.Item
          label="Part Category"
          name="part_category_id"
          rules={[{ required: true, message: 'Please select a part category' }]}
        >
          <Select
            placeholder="Select part category"
            size="large"
            showSearch
            loading={partCategoriesLoading}
            disabled={isUpdate} // Cannot change category after creation
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

        <Form.Item
          label="Coverage Conditions"
          name="coverage_conditions"
          validateFirst
          rules={[{ max: 1000, message: 'Coverage conditions cannot exceed 1000 characters' }]}
        >
          <TextArea
            placeholder="Enter coverage conditions (optional)"
            rows={4}
            size="large"
            showCount
            maxLength={1000}
          />
        </Form.Item>

        <Form.Item style={{ marginBottom: 0, textAlign: 'right', marginTop: '24px' }}>
          <Space>
            <Button size="large" onClick={onClose}>
              Cancel
            </Button>
            <Button type="primary" htmlType="submit" loading={loading} size="large">
              {isUpdate ? 'Update' : 'Add'}
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  )
}

export default CoveredPartModal
