import React from 'react'
import { Button, Result } from 'antd'
import { useNavigate } from 'react-router-dom'
import './Error.less'
import { ERROR_MESSAGES } from '@constants'

const Error = ({ code = 404 }) => {
  const navigate = useNavigate()

  const handleGoHome = () => {
    navigate('/')
  }

  const handleGoBack = () => {
    code === 403 ? navigate(-3) : navigate(-1)
  }

  return (
    <div className="error-container">
      <Result
        status={code}
        title={code}
        subTitle={ERROR_MESSAGES[code]}
        extra={
          <div className="error-actions">
            <Button type="primary" onClick={handleGoHome}>
              Back Home
            </Button>
            <Button onClick={handleGoBack}>Go Back</Button>
          </div>
        }
      />
    </div>
  )
}

export default Error
