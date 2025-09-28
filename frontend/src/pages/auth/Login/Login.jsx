import React, { useState } from 'react'
import { Form, Input, Button, Card, message, Typography, Checkbox } from 'antd'
import {
  UserOutlined,
  LockOutlined,
  GoogleOutlined,
  EyeInvisibleOutlined,
  EyeTwoTone,
} from '@ant-design/icons'
import './Login.less'
import { useDispatch, useSelector } from 'react-redux'
import { API_ENDPOINTS } from '@constants'
import api from '@services/api.js'
import { loginStart, loginSuccess, loginFailure } from '@redux/authSlice.js'
import Logo from '@pages/Login/Logo/Logo.jsx'

const { Title } = Typography

const Login = ({ onLoginSuccess }) => {
  const [form] = Form.useForm()
  const [googleLoading, setGoogleLoading] = useState(false)

  const dispatch = useDispatch()
  const { isLoading } = useSelector((state) => state.auth)

  const onFinish = async (values) => {
    dispatch(loginStart())
    try {
      const res = await api.post(API_ENDPOINTS.AUTH.LOGIN, values)
      const { refresh_token, access_token, user } = res.data.data

      dispatch(loginSuccess({ user, token: token }))

      message.success('Login successful!')
      if (onLoginSuccess) {
        onLoginSuccess(values)
      }
    } catch (error) {
      console.error(error)
      dispatch(loginFailure())
      message.error('Login failed. Please check your credentials.')
    }
  }

  const handleGoogleLogin = async () => {
    setGoogleLoading(true)
    try {
      console.log('Google login initiated')
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
    <div className={`login-container ${isLoading ? 'login-loading' : ''}`}>
      <Card className="login-card">
        <div className="login-header">
          <Logo />
          <Title level={2} className="login-title">
            EV Warranty System
          </Title>
        </div>

        <div className="login-body">
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
              validateFirst
              validateTrigger="onBlur"
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
              <Input prefix={<UserOutlined />} placeholder="Enter your email" size="large" />
            </Form.Item>

            <Form.Item
              label="Password"
              name="password"
              validateFirst
              validateTrigger="onBlur"
              rules={[
                {
                  required: true,
                  message: 'Please input your password!',
                },
                {
                  min: 8,
                  message: 'Password must be at least 8 characters long!',
                },
                {
                  pattern: /[a-z]/,
                  message: 'Password must contain at least one lowercase letter!',
                },
                {
                  pattern: /[A-Z]/,
                  message: 'Password must contain at least one uppercase letter!',
                },
                {
                  pattern: /\d/,
                  message: 'Password must contain at least one digit!',
                },
                {
                  pattern: /[^A-Za-z0-9]/,
                  message: 'Password must contain at least one special character!',
                },
              ]}
            >
              <Input.Password
                prefix={<LockOutlined />}
                placeholder="Enter your password"
                size="large"
                iconRender={(visible) => (visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />)}
              />
            </Form.Item>

            <div className="form-options">
              <Form.Item name="remember" valuePropName="checked" noStyle>
                <Checkbox>Remember me</Checkbox>
              </Form.Item>
              <a href="#" className="forgot-password">
                Forgot password?
              </a>
            </div>

            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                loading={isLoading}
                size="large"
                block
                className="login-button"
              >
                {isLoading ? 'Signing in...' : 'Sign In'}
              </Button>
            </Form.Item>

            <div className="divider">
              <span className="divider-text">or</span>
            </div>

            <div className="social-login">
              <Button
                icon={<GoogleOutlined />}
                onClick={handleGoogleLogin}
                loading={googleLoading}
                size="large"
                block
                className="social-button"
              >
                {googleLoading ? 'Connecting...' : 'Continue with Google'}
              </Button>
            </div>
          </Form>
        </div>
      </Card>
    </div>
  )
}

export default Login
