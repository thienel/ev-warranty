import React, { useEffect } from 'react'
import { Modal, Button, Form, Input, message, Select, Space } from 'antd'
import { AppstoreOutlined } from '@ant-design/icons'
import { API_ENDPOINTS } from '@constants/common-constants.js'
import { type PartCategoryModalProps, type PartCategoryFormData } from '@/types/index.js'
import api from '@services/api.js'
import useHandleApiError from '@/hooks/useHandleApiError.js'

const PartCategoryModal: React.FC<PartCategoryModalProps> = ({
  loading,
  setLoading,
  onClose,
  partCategory = null,
  opened = false,
  isUpdate,
  partCategories,
  partCategoriesLoading,
}) => {
  const [form] = Form.useForm<PartCategoryFormData>()
  const handleError = useHandleApiError()

  // Get available parent categories (excluding the current category if editing)
  const availableParentCategories = React.useMemo(() => {
    if (!partCategory || !isUpdate) {
      return partCategories
    }
    // Exclude current category and its descendants from parent options
    return partCategories.filter((cat) => cat.id !== partCategory.id)
  }, [partCategories, partCategory, isUpdate])

  // Populate form when partCategory prop changes or modal opens
  useEffect(() => {
    if (opened) {
      if (partCategory && isUpdate) {
        // When editing, populate form with part category data
        const formData: PartCategoryFormData = {
          category_name: partCategory.category_name,
          description: partCategory.description,
          parent_category_id: partCategory.parent_category_id,
        }
        form.setFieldsValue(formData)
      } else {
        // When creating new, reset to default values
        form.resetFields()
      }
    }
  }, [form, partCategory, isUpdate, opened])

  // Clear form when modal closes
  useEffect(() => {
    if (!opened) {
      form.resetFields()
    }
  }, [form, opened])

  const handleSubmit = async (values: PartCategoryFormData): Promise<void> => {
    setLoading(true)
    try {
      // Ensure empty parent_category_id is sent as undefined/null instead of empty string
      const payload = {
        ...values,
        parent_category_id: values.parent_category_id || undefined,
      }

      if (isUpdate) {
        await api.put(`${API_ENDPOINTS.PART_CATEGORIES}/${partCategory?.id}`, payload)
        message.success('Part category updated successfully')
      } else {
        await api.post(API_ENDPOINTS.PART_CATEGORIES, payload)
        message.success('Part category created successfully')
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
          {isUpdate ? 'Edit Part Category' : 'Add New Part Category'}
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
        key={partCategory?.id || 'new'}
      >
        <Form.Item
          label="Category Name"
          name="category_name"
          validateFirst
          rules={[
            { required: true, message: 'Please enter category name' },
            { min: 2, message: 'Category name must be at least 2 characters' },
            { max: 255, message: 'Category name cannot exceed 255 characters' },
          ]}
        >
          <Input placeholder="Enter category name" prefix={<AppstoreOutlined />} size="large" />
        </Form.Item>

        <Form.Item
          label="Description"
          name="description"
          validateFirst
          rules={[{ max: 1000, message: 'Description cannot exceed 1000 characters' }]}
        >
          <Input.TextArea
            placeholder="Enter category description (optional)"
            rows={3}
            size="large"
          />
        </Form.Item>

        <Form.Item
          label="Parent Category"
          name="parent_category_id"
          tooltip="Select a parent category to create a subcategory. Leave empty for root category."
        >
          <Select
            placeholder="Select parent category (optional)"
            size="large"
            allowClear
            showSearch
            loading={partCategoriesLoading}
            filterOption={(input, option) => {
              const label = option?.label as string
              return label?.toLowerCase().includes(input.toLowerCase())
            }}
            options={availableParentCategories.map((cat) => ({
              label: `${cat.category_name}${cat.parent_category_name ? ` (under ${cat.parent_category_name})` : ''}`,
              value: cat.id,
            }))}
          />
        </Form.Item>

        <Form.Item style={{ marginBottom: 0, textAlign: 'right', marginTop: '24px' }}>
          <Space>
            <Button size="large" onClick={onClose}>
              Cancel
            </Button>
            <Button type="primary" htmlType="submit" loading={loading} size="large">
              {isUpdate ? 'Update Category' : 'Create Category'}
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  )
}

export default PartCategoryModal
