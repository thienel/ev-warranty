import React, { useEffect, useState } from 'react'
import { Form, Card, message, Typography } from 'antd'
import './Login.less'
import { useDispatch } from 'react-redux'
import { API_BASE_URL, API_ENDPOINTS } from '@constants'
import api from '@services/api.js'
import { loginSuccess } from '@redux/authSlice.js'
import Logo from '@pages/auth/Login/Logo/Logo.jsx'
import { useDelay } from '@/hooks/index.js'
import { useNavigate, useSearchParams } from 'react-router-dom'
import LoginForm from '@pages/auth/Login/LoginForm/LoginForm.jsx'

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
        dispatch(loginSuccess({ user, token: token, remember: values.remember }))
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
          <LoginForm
            form={form}
            onLogin={handleLogin}
            onGoogleLogin={handleGoogleLogin}
            loginLoading={loginLoading}
            googleLoading={googleLoading}
          />
        </div>
      </Card>
    </div>
  )
}

export default Login
