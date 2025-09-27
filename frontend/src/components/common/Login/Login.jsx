import React, { useState } from 'react'
import { Form, Input, Button, Card, Divider, message, Typography, Checkbox } from 'antd'
import {
  UserOutlined,
  LockOutlined,
  GoogleOutlined,
  EyeInvisibleOutlined,
  EyeTwoTone,
} from '@ant-design/icons'
import './Login.less'

const { Title, Text } = Typography

const Login = ({ onLoginSuccess }) => {
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [googleLoading, setGoogleLoading] = useState(false)

  const onFinish = async (values) => {
    setLoading(true)
    try {
      // Simulate API call - replace with actual login API
      console.log('Login values:', values)

      // Mock successful login
      await new Promise((resolve) => setTimeout(resolve, 1000))

      message.success('Login successful!')

      // Call success callback if provided
      if (onLoginSuccess) {
        onLoginSuccess(values)
      }
    } catch {
      message.error('Login failed. Please check your credentials.')
    } finally {
      setLoading(false)
    }
  }

  const handleGoogleLogin = async () => {
    setGoogleLoading(true)
    try {
      // Implement Google OAuth login here
      console.log('Google login initiated')

      // Mock Google login
      await new Promise((resolve) => setTimeout(resolve, 1000))

      message.success('Google login successful!')

      if (onLoginSuccess) {
        onLoginSuccess({ loginType: 'google' })
      }
    } catch {
      message.error('Google login failed. Please try again.')
    } finally {
      setGoogleLoading(false)
    }
  }

  const onFinishFailed = (errorInfo) => {
    console.log('Failed:', errorInfo)
    message.error('Please fill in all required fields correctly.')
  }

  return (
    <div className="login-container">
      <Card className="login-card">
        <div className="login-header">
          <Title level={2} className="login-title">
            EV Warranty System
          </Title>
          <Text type="secondary" className="login-subtitle">
            Sign in to your account
          </Text>
        </div>

        <Form
          form={form}
          name="login"
          className="login-form"
          initialValues={{
            remember: true,
          }}
          onFinish={onFinish}
          onFinishFailed={onFinishFailed}
          autoComplete="off"
          layout="vertical"
        >
          <Form.Item
            label="Email"
            name="email"
            rules={[
              {
                required: true,
                message: 'Please input your email!',
              },
              {
                type: 'email',
                message: 'Please enter a valid email address!',
              },
            ]}
          >
            <Input
              prefix={<UserOutlined className="input-icon" />}
              placeholder="Enter your email"
              size="large"
              className="login-input"
            />
          </Form.Item>

          <Form.Item
            label="Password"
            name="password"
            rules={[
              {
                required: true,
                message: 'Please input your password!',
              },
              {
                min: 6,
                message: 'Password must be at least 6 characters long!',
              },
            ]}
          >
            <Input.Password
              prefix={<LockOutlined className="input-icon" />}
              placeholder="Enter your password"
              size="large"
              className="login-input"
              iconRender={(visible) => (visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />)}
            />
          </Form.Item>

          <Form.Item>
            <div className="login-options">
              <Form.Item name="remember" valuePropName="checked" noStyle>
                <Checkbox>Remember me</Checkbox>
              </Form.Item>
              <Button type="link" className="forgot-password">
                Forgot password?
              </Button>
            </div>
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              loading={loading}
              size="large"
              className="login-button"
              block
            >
              {loading ? 'Signing in...' : 'Sign In'}
            </Button>
          </Form.Item>
        </Form>

        <Divider className="login-divider">
          <Text type="secondary">or</Text>
        </Divider>

        <Button
          icon={<GoogleOutlined />}
          onClick={handleGoogleLogin}
          loading={googleLoading}
          size="large"
          className="google-login-button"
          block
        >
          {googleLoading ? 'Connecting...' : 'Continue with Google'}
        </Button>
      </Card>
    </div>
  )
}

export default Login
