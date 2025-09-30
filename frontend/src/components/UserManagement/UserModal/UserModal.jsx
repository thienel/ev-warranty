import React, { useEffect } from 'react'
import { Modal, Button, Form, Input, message, Select, Space, Switch } from 'antd'
import { HomeOutlined, LockOutlined, MailOutlined, UserOutlined } from '@ant-design/icons'
import { API_ENDPOINTS, PASSWORD_RULES, ROLE_LABELS, USER_ROLES } from '@constants'
import api from '@services/api.js'

const UserModal = ({
  loading,
  setLoading,
  onClose,
  user = null,
  opened = false,
  offices,
  isUpdate,
}) => {
  const [form] = Form.useForm()

  const { Option } = Select

  useEffect(() => {
    if (user) {
      form.setFieldsValue({
        ...user,
        password: '',
      })
    } else {
      form.resetFields()
      form.setFieldsValue({ is_active: true })
    }
  })

  const handleSubmit = async (values) => {
    setLoading(true)
    try {
      let response

      const payload = { ...values }

      if (!isUpdate && !payload.password) {
        message.warning('Password is required for new user')
        setLoading(false)
        return
      }

      delete payload.password
      delete payload.email

      if (isUpdate) {
        response = await api.put(`${API_ENDPOINTS.USER}${user.id}`, payload)

        if (response.data.success) {
          message.success('User updated successfully')
        }
      } else {
        response = await api.post(API_ENDPOINTS.USER, payload)

        if (response.data.success) {
          message.success('User created successfully')
        }
      }

      form.resetFields()
      onClose()
    } catch (error) {
      const errorMsg =
        error.response?.data?.message ||
        (isUpdate ? 'Failed to update user' : 'Failed to create user')
      message.error(errorMsg)
      console.error('Error saving user:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Modal
      title={
        <Space>
          <UserOutlined />
          {isUpdate ? 'Edit User' : 'Add New User'}
        </Space>
      }
      open={opened}
      onCancel={() => {
        form.resetFields()
        onClose()
      }}
      footer={null}
      width={600}
      destroyOnHidden
    >
      <Form form={form} layout="vertical" onFinish={handleSubmit} autoComplete="off">
        <Form.Item
          label="Full Name"
          name="name"
          rules={[
            { required: true, message: 'Please enter full name' },
            { min: 2, message: 'Name must be at least 2 characters' },
            { max: 100, message: 'Name cannot exceed 100 characters' },
            {
              pattern: /^[a-zA-ZÀ-ỹ\s]+$/,
              message: 'Name can only contain letters and spaces',
            },
          ]}
        >
          <Input placeholder="Enter full name" prefix={<UserOutlined />} size="large" />
        </Form.Item>

        <Form.Item
          label="Email"
          name="email"
          rules={[
            { required: true, message: 'Please enter email' },
            { type: 'email', message: 'Invalid email format' },
            { max: 100, message: 'Email cannot exceed 100 characters' },
          ]}
        >
          <Input
            placeholder="Enter email"
            prefix={<MailOutlined />}
            size="large"
            disabled={isUpdate}
          />
        </Form.Item>

        {!isUpdate && (
          <Form.Item
            label="Password"
            name="password"
            validateFirst
            rules={[{ required: true, message: 'Please enter password' }, ...PASSWORD_RULES]}
          >
            <Input.Password
              placeholder={isUpdate ? 'Enter new password (optional)' : 'Enter password'}
              prefix={<LockOutlined />}
              size="large"
            />
          </Form.Item>
        )}

        <Form.Item
          label="Role"
          name="role"
          rules={[{ required: true, message: 'Please select a role' }]}
        >
          <Select placeholder="Select role" size="large">
            {Object.entries(USER_ROLES).map(([key, value]) => (
              <Option key={key} value={value}>
                <Space>
                  <UserOutlined />
                  {ROLE_LABELS[value]}
                </Space>
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label="Office"
          name="office_id"
          rules={[{ required: true, message: 'Please select an office' }]}
        >
          <Select
            placeholder="Select office"
            size="large"
            showSearch
            optionFilterProp="children"
            filterOption={(input, option) =>
              option.children.toLowerCase().includes(input.toLowerCase())
            }
          >
            {offices.map((office) => (
              <Option key={office.id} value={office.id}>
                <Space>
                  <HomeOutlined />
                  {office.office_name}
                </Space>
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item label="Status" name="is_active" valuePropName="checked">
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
              {isUpdate ? 'Update User' : 'Create User'}
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  )
}

export default UserModal
