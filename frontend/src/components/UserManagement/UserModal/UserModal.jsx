import React from 'react'
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

  const handleSubmit = async (values) => {
    setLoading(true)
    try {
      const payload = { ...values }

      if (!isUpdate && !payload.password) {
        message.warning('Password is required for new user')
        setLoading(false)
        return
      }

      if (isUpdate) {
        delete payload.password
        delete payload.email
        await api.put(`${API_ENDPOINTS.USER}/${user.id}`, payload)
        message.success('User updated successfully')
      } else {
        await api.post(API_ENDPOINTS.USER, payload)
        message.success('User created successfully')
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
      title={<Space style={{ margin: '14px 0' }}>{isUpdate ? 'Edit User' : 'Add New User'}</Space>}
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
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
        autoComplete="off"
        key={user?.id || 'new'}
        initialValues={user || { is_active: true }}
      >
        <Form.Item
          label="Full Name"
          name="name"
          validateFirst
          rules={[
            { required: true, message: 'Please enter full name' },
            { min: 2, message: 'Name must be at least 2 characters' },
            { max: 50, message: 'Name cannot exceed 50 characters' },
            {
              pattern: /^[\p{L}\s'-]+$/u,
              message: 'Name can only contain letters, spaces, apostrophes or hyphens',
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
          <Form.Item label="Password" name="password" validateFirst rules={PASSWORD_RULES}>
            <Input.Password placeholder="Enter password" prefix={<LockOutlined />} size="large" />
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
