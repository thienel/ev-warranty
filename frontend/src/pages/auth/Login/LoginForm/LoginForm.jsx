import React from 'react'
import { Button, Checkbox, Form, Input } from 'antd'
import {
  EyeInvisibleOutlined,
  EyeTwoTone,
  GoogleOutlined,
  LockOutlined,
  UserOutlined,
} from '@ant-design/icons'
import './LoginForm.less'
import { EMAIL_RULES, PASSWORD_RULES } from '@constants'

const LoginForm = ({ form, onLogin, onGoogleLogin, loginLoading, googleLoading }) => {
  return (
    <Form
      form={form}
      name="login"
      className="login-form"
      initialValues={{
        remember: false,
      }}
      onFinish={onLogin}
      autoComplete="off"
      layout="vertical"
    >
      <Form.Item
        label="Email"
        name="email"
        validateFirst
        validateTrigger="onBlur"
        rules={EMAIL_RULES}
      >
        <Input prefix={<UserOutlined />} placeholder="Enter your email" size="large" />
      </Form.Item>

      <Form.Item
        label="Password"
        name="password"
        validateFirst
        validateTrigger="onBlur"
        rules={PASSWORD_RULES}
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
          onClick={onGoogleLogin}
          loading={googleLoading}
          size="large"
          block
          className="social-button"
        >
          {googleLoading ? 'Connecting...' : 'Continue with Google'}
        </Button>
      </div>
    </Form>
  )
}

export default LoginForm
