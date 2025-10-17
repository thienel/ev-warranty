import React, { useEffect } from 'react'
import { Modal, Button, Form, Input, message, Select, Space, Switch } from 'antd'
import { BankOutlined, EnvironmentOutlined } from '@ant-design/icons'
import { API_ENDPOINTS } from '@constants'
import api from '@services/api.js'

const OfficeModal = ({ loading, setLoading, onClose, office = null, opened = false, isUpdate }) => {
  const [form] = Form.useForm()
  const { Option } = Select

  useEffect(() => {
    if (office) {
      form.setFieldsValue({
        ...office,
      })
    } else {
      form.resetFields()
    }
  })

  const handleSubmit = async (values) => {
    setLoading(true)
    try {
      let response

      const payload = { ...values }

      if (isUpdate) {
        response = await api.put(`${API_ENDPOINTS.OFFICE}/${office.id}`, payload)

        if (response.data.success) {
          message.success('Office updated successfully')
        }
      } else {
        response = await api.post(API_ENDPOINTS.OFFICE, payload)

        if (response.data.success) {
          message.success('Office created successfully')
        }
      }

      form.resetFields()
      onClose()
    } catch (error) {
      const errorMsg =
        error.response?.data?.message ||
        (isUpdate ? 'Failed to update office' : 'Failed to create office')
      message.error(errorMsg)
      console.error('Error saving office:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Modal
      title={
        <Space style={{ margin: '14px 0' }}>{isUpdate ? 'Edit Office' : 'Add New Office'}</Space>
      }
      open={opened}
      onCancel={() => {
        form.resetFields()
        onClose()
      }}
      style={{ margin: 'auto' }}
      footer={null}
      width={500}
      destroyOnHidden
    >
      <Form form={form} layout="vertical" onFinish={handleSubmit} autoComplete="off">
        <Form.Item
          label="Office Name"
          name="office_name"
          validateFirst
          rules={[
            { required: true, message: 'Please enter office name' },
            { min: 2, message: 'Office name must be at least 2 characters' },
            { max: 100, message: 'Office name cannot exceed 100 characters' },
          ]}
        >
          <Input placeholder="Enter office name" prefix={<BankOutlined />} size="large" />
        </Form.Item>

        <Form.Item
          label="Office Type"
          name="office_type"
          rules={[{ required: true, message: 'Please select office type' }]}
        >
          <Select placeholder="Select office type" size="large">
            <Option value="evm">
              <Space>
                <BankOutlined />
                EVM
              </Space>
            </Option>
            <Option value="sc">
              <Space>
                <BankOutlined />
                Service Center
              </Space>
            </Option>
          </Select>
        </Form.Item>

        <Form.Item
          label="Address"
          name="address"
          validateFirst
          rules={[
            { required: true, message: 'Please enter address' },
            { min: 5, message: 'Address must be at least 5 characters' },
            { max: 200, message: 'Address cannot exceed 200 characters' },
          ]}
        >
          <Input.TextArea placeholder="Enter office address" rows={3} size="large" />
        </Form.Item>

        <Form.Item label="Status" name="is_active" valuePropName="checked" initialValue={true}>
          <Switch checkedChildren="Active" unCheckedChildren="Inactive" size="default" />
        </Form.Item>

        <Form.Item style={{ marginBottom: 0, textAlign: 'right', marginTop: '24px' }}>
          <Space>
            <Button
              size="large"
              onClick={() => {
                form.resetFields()
                onClose()
              }}
            >
              Cancel
            </Button>
            <Button type="primary" htmlType="submit" loading={loading} size="large">
              {isUpdate ? 'Update Office' : 'Create Office'}
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  )
}

export default OfficeModal
