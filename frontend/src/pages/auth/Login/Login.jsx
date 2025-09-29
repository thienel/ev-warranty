import React, { useEffect, useState } from 'react'
import { Form, Input, Button, Card, message, Typography, Checkbox } from 'antd'
import {
  UserOutlined,
  LockOutlined,
  GoogleOutlined,
  EyeInvisibleOutlined,
  EyeTwoTone,
} from '@ant-design/icons'
import './Login.less'
import { useDispatch } from 'react-redux'
import { API_BASE_URL, API_ENDPOINTS } from '@constants'
import api from '@services/api.js'
import { loginSuccess} from '@redux/authSlice.js'
import Logo from '@pages/auth/Login/Logo/Logo.jsx'
import { useDelay } from '@/hooks/index.js'
import { useNavigate, useSearchParams } from 'react-router-dom'

const { Title } = Typography

const Login = () => {
  const [form] = Form.useForm()
  const [loginLoading, setLoginLoading] = useState(false)
  const [googleLoading, setGoogleLoading] = useState(false)
  const navigate = useNavigate()
  const dispatch = useDispatch()
  const [searchParams] = useSearchParams()
  const delay = useDelay(500)

  useEffect(() => {
    const error = searchParams.get('error')
    if (error) {
      navigate('/login')
      message.error(error)
    }
  }, [])

  const handleLogin = async (values) => {
    setLoginLoading(true)
    delay(async () => {
      try {
        const res = await api.post(API_ENDPOINTS.AUTH.LOGIN, values)
        const { token, user } = res.data.data
        message.success('Login successful!')
        dispatch(loginSuccess({ user, token: token }))
      } catch (error) {
        console.error(error.response.data)
        message.error('Login failed. Please check your credentials.')
      } finally {
        setLoginLoading(false)
      }
    })
  }

  const handleGoogleLogin = async () => {
    setGoogleLoading(true)
    delay(() => {
      window.location.href = `${API_BASE_URL}${API_ENDPOINTS.AUTH.GOOGLE}`
    })
  }

  return (
    <div className={`login-container ${loginLoading ? 'login-loading' : ''}`}>
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
            onFinish={handleLogin}
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
                loading={loginLoading}
                size="large"
                block
                className="login-button"
              >
                {loginLoading ? 'Signing in...' : 'Sign In'}
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
