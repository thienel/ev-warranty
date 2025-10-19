import React, { useEffect, useState } from 'react'
import { Form, Card, message, Typography } from 'antd'
import './Login.less'
import { useDispatch } from 'react-redux'
import { API_BASE_URL, API_ENDPOINTS } from '@constants/common-constants.js'
import api from '@services/api.js'
import { loginSuccess } from '@redux/authSlice.js'
import  useDelay  from '@/hooks/useDelay.js'
import { useNavigate, useSearchParams } from 'react-router-dom'
import LoginForm from '@pages/Login/LoginForm/LoginForm.jsx'
import { ThunderboltOutlined } from '@ant-design/icons'
import useHandleApiError from '@/hooks/useHandleApiError.js'

const { Title } = Typography

const Login = () => {
  const [form] = Form.useForm()
  const [loginLoading, setLoginLoading] = useState(false)
  const [googleLoading, setGoogleLoading] = useState(false)
  const navigate = useNavigate()
  const dispatch = useDispatch()
  const [searchParams] = useSearchParams()
  const delay = useDelay(500)
  const handleError = useHandleApiError()

  useEffect(() => {
    const error = searchParams.get('error')
    if (error) {
      navigate('/login')
      message.error(error)
    }
  }, [navigate, searchParams])

  const handleLogin = async (values) => {
    setLoginLoading(true)
    delay(async () => {
      try {
        const res = await api.post(API_ENDPOINTS.AUTH.LOGIN, values)
        const { token, user } = res.data.data
        message.success('Login successful!')
        dispatch(loginSuccess({ user, token: token, remember: values.remember }))
      } catch (error) {
        handleError(error)
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
          <div className="main-logo">
            <ThunderboltOutlined />
          </div>
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
